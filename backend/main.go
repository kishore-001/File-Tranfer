package main

import (
	"backend/logic"
	"log"
	"net/http"
)

func main() {
	// Initialize system info on startup
	logic.InitSystemInfo()

	// Start peer discovery service
	go logic.StartPeerDiscovery()

	// Start gRPC server for incoming file transfers
	go logic.StartGRPCServer(9002)

	// Create a new ServeMux
	mux := http.NewServeMux()

	// Register API endpoints
	mux.HandleFunc("/api/systeminfo", logic.GetSystemInfo)
	mux.HandleFunc("/api/peers", logic.GetPeers)
	mux.HandleFunc("/api/filetransfer", logic.HandleFileTransfer)

	// Add CORS middleware for frontend communication
	handler := enableCORS(mux)

	// Start REST server
	log.Println("Backend running on port 9001")
	log.Fatal(http.ListenAndServe("80", handler))
}

// CORS middleware for frontend communication
func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:9000")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
