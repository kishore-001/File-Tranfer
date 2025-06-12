package logic

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
	"log"
	"net/http"
	"os"
	"runtime"
	"time"
)

type SystemInfo struct {
	Hostname  string `json:"hostname"`
	CPU       string `json:"cpu"`
	RAM       string `json:"ram"`
	OS        string `json:"os"`
	PeerID    string `json:"peer_id"`
	Timestamp string `json:"timestamp"`
}

var systemInfo SystemInfo

const systemInfoFile = "system_info.json"

// InitSystemInfo initializes system info on server startup
func InitSystemInfo() {
	// Check if system_info.json exists
	if _, err := os.Stat(systemInfoFile); os.IsNotExist(err) {
		log.Println("system_info.json not found, creating new file...")
		createSystemInfoFile()
	} else {
		log.Println("Loading existing system_info.json...")
		loadSystemInfoFromFile()
	}

	log.Printf("System initialized - Hostname: %s, PeerID: %s",
		systemInfo.Hostname, systemInfo.PeerID)
}

// createSystemInfoFile gathers system data and creates the JSON file
func createSystemInfoFile() {
	// Gather system information
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}

	// Get CPU information
	cpuInfo, err := cpu.Info()
	cpuModel := "Unknown CPU"
	if err == nil && len(cpuInfo) > 0 {
		cpuModel = cpuInfo[0].ModelName
	}

	// Get memory information
	memInfo, err := mem.VirtualMemory()
	ramSize := "Unknown RAM"
	if err == nil {
		ramSize = fmt.Sprintf("%.2f GB", float64(memInfo.Total)/1024/1024/1024)
	}

	// Get OS information
	hostInfo, err := host.Info()
	osInfo := runtime.GOOS
	if err == nil {
		osInfo = fmt.Sprintf("%s %s", hostInfo.Platform, hostInfo.PlatformVersion)
	}

	// Generate unique peer ID
	peerID := generatePeerID(hostname)

	// Create system info struct
	systemInfo = SystemInfo{
		Hostname:  hostname,
		CPU:       cpuModel,
		RAM:       ramSize,
		OS:        osInfo,
		PeerID:    peerID,
		Timestamp: time.Now().Format(time.RFC3339),
	}

	// Save to file
	saveSystemInfoToFile()
}

// loadSystemInfoFromFile loads existing system info from JSON file
func loadSystemInfoFromFile() {
	file, err := os.Open(systemInfoFile)
	if err != nil {
		log.Printf("Error opening system_info.json: %v", err)
		createSystemInfoFile() // Fallback to creating new file
		return
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&systemInfo); err != nil {
		log.Printf("Error decoding system_info.json: %v", err)
		createSystemInfoFile() // Fallback to creating new file
		return
	}

	log.Println("System info loaded successfully from file")
}

// saveSystemInfoToFile saves current system info to JSON file
func saveSystemInfoToFile() {
	file, err := os.Create(systemInfoFile)
	if err != nil {
		log.Printf("Error creating system_info.json: %v", err)
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "    ") // Pretty print JSON

	if err := encoder.Encode(systemInfo); err != nil {
		log.Printf("Error encoding system_info.json: %v", err)
		return
	}

	log.Println("System info saved to system_info.json")
}

// GetSystemInfo HTTP handler that returns system info as JSON
func GetSystemInfo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Set response headers
	w.Header().Set("Content-Type", "application/json")

	// Encode and send system info
	if err := json.NewEncoder(w).Encode(systemInfo); err != nil {
		log.Printf("Error encoding system info response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// generatePeerID creates a unique peer identifier
func generatePeerID(hostname string) string {
	// Generate random bytes
	bytes := make([]byte, 4)
	if _, err := rand.Read(bytes); err != nil {
		// Fallback to timestamp if random generation fails
		return fmt.Sprintf("peer_%s_%d", hostname, time.Now().Unix())
	}

	randomStr := hex.EncodeToString(bytes)
	return fmt.Sprintf("peer_%s_%s", hostname, randomStr)
}

// GetSystemInfoStruct returns the current system info struct (for internal use)
func GetSystemInfoStruct() SystemInfo {
	return systemInfo
}

// RefreshSystemInfo re-gathers system info and updates the file
func RefreshSystemInfo() {
	log.Println("Refreshing system information...")
	createSystemInfoFile()
}
