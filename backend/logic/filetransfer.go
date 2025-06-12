package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"time"

	pb "backend/proto" // Replace with your actual module path
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type FileTransferRequest struct {
	PeerID string `json:"peerid"`
	File   string `json:"file"`
}

type FileTransferResponse struct {
	Message string `json:"message"`
	Peer    string `json:"peer"`
	File    string `json:"file"`
	Status  string `json:"status"`
}

// gRPC server implementation
type fileTransferServer struct {
	pb.UnimplementedFileTransferServiceServer
}

// StartGRPCServer starts the gRPC server on specified port
func StartGRPCServer(port int) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("Failed to listen on port %d: %v", port, err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterFileTransferServiceServer(grpcServer, &fileTransferServer{})

	log.Printf("gRPC server listening on port %d", port)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC: %v", err)
	}
}

// SendFile handles incoming file transfers via gRPC streaming
func (s *fileTransferServer) SendFile(stream pb.FileTransferService_SendFileServer) error {
	var fileName string
	var totalBytes int64
	var receivedChunks int64

	// Create downloads directory if it doesn't exist
	downloadsDir := "./downloads"
	if err := os.MkdirAll(downloadsDir, 0755); err != nil {
		log.Printf("Error creating downloads directory: %v", err)
		return err
	}

	var file *os.File

	for {
		chunk, err := stream.Recv()
		if err == io.EOF {
			// End of stream - close file and send response
			if file != nil {
				file.Close()
			}

			log.Printf("File transfer completed: %s (%d bytes)", fileName, totalBytes)

			return stream.SendAndClose(&pb.FileTransferResponse{
				Success:       true,
				Message:       fmt.Sprintf("File %s received successfully", fileName),
				BytesReceived: totalBytes,
			})
		}

		if err != nil {
			log.Printf("Error receiving chunk: %v", err)
			if file != nil {
				file.Close()
				// Clean up partial file
				os.Remove(filepath.Join(downloadsDir, fileName))
			}
			return err
		}

		// First chunk - create file
		if file == nil {
			fileName = chunk.FileName
			filePath := filepath.Join(downloadsDir, fileName)

			file, err = os.Create(filePath)
			if err != nil {
				log.Printf("Error creating file %s: %v", filePath, err)
				return err
			}

			log.Printf("Starting to receive file: %s", fileName)
		}

		// Write chunk data to file
		bytesWritten, err := file.Write(chunk.Data)
		if err != nil {
			log.Printf("Error writing to file: %v", err)
			file.Close()
			os.Remove(filepath.Join(downloadsDir, fileName))
			return err
		}

		totalBytes += int64(bytesWritten)
		receivedChunks++

		log.Printf("Received chunk %d for %s (%d bytes)", receivedChunks, fileName, bytesWritten)
	}
}

// HandleFileTransfer HTTP handler for file transfer requests
func HandleFileTransfer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req FileTransferRequest

	// Parse JSON request body
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.PeerID == "" || req.File == "" {
		http.Error(w, "Missing required fields: peerid and file", http.StatusBadRequest)
		return
	}

	// Find peer by ID
	peer := GetPeerByID(req.PeerID)
	if peer == nil {
		http.Error(w, "Peer not found", http.StatusNotFound)
		return
	}

	// Check if file exists
	if _, err := os.Stat(req.File); os.IsNotExist(err) {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	// Start file transfer in goroutine
	go func() {
		if err := sendFileToP2P(peer, req.File); err != nil {
			log.Printf("File transfer failed: %v", err)
		}
	}()

	response := FileTransferResponse{
		Message: "File transfer initiated",
		Peer:    peer.Hostname,
		File:    filepath.Base(req.File),
		Status:  "started",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// sendFileToP2P sends a file to a peer via gRPC streaming
func sendFileToP2P(peer *Peer, filePath string) error {
	// Connect to peer's gRPC server
	conn, err := grpc.Dial(
		fmt.Sprintf("%s:%d", peer.IP, peer.Port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return fmt.Errorf("failed to connect to peer %s: %v", peer.Hostname, err)
	}
	defer conn.Close()

	client := pb.NewFileTransferServiceClient(conn)

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// Start streaming
	stream, err := client.SendFile(ctx)
	if err != nil {
		return fmt.Errorf("failed to create stream: %v", err)
	}

	// Open file for reading
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	// Get file info
	fileInfo, err := file.Stat()
	if err != nil {
		return fmt.Errorf("failed to get file info: %v", err)
	}

	fileName := filepath.Base(filePath)
	fileSize := fileInfo.Size()
	chunkSize := int64(1024 * 64) // 64KB chunks
	totalChunks := (fileSize + chunkSize - 1) / chunkSize

	log.Printf("Sending file %s to %s (%d bytes, %d chunks)", fileName, peer.Hostname, fileSize, totalChunks)

	// Send file in chunks
	buffer := make([]byte, chunkSize)
	chunkNumber := int64(0)

	for {
		bytesRead, err := file.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read file: %v", err)
		}

		chunkNumber++

		chunk := &pb.FileChunk{
			FileName:    fileName,
			Data:        buffer[:bytesRead],
			ChunkNumber: chunkNumber,
			TotalChunks: totalChunks,
		}

		if err := stream.Send(chunk); err != nil {
			return fmt.Errorf("failed to send chunk %d: %v", chunkNumber, err)
		}

		log.Printf("Sent chunk %d/%d (%d bytes)", chunkNumber, totalChunks, bytesRead)
	}

	// Close stream and get response
	response, err := stream.CloseAndRecv()
	if err != nil {
		return fmt.Errorf("failed to close stream: %v", err)
	}

	if response.Success {
		log.Printf("File transfer successful: %s", response.Message)
	} else {
		log.Printf("File transfer failed: %s", response.Message)
	}

	return nil
}
