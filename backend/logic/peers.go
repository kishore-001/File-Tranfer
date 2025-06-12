package logic

import (
	"context"
	"encoding/json"
	"github.com/grandcat/zeroconf"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

type Peer struct {
	ID       string `json:"peer_id"`
	Hostname string `json:"hostname"`
	IP       string `json:"ip"`
	Port     int    `json:"port"`
	CPU      string `json:"cpu"`
	RAM      string `json:"ram"`
	OS       string `json:"os"`
	Status   string `json:"status"`
}

type PeersResponse struct {
	Peers []Peer `json:"peers"`
	Count int    `json:"count"`
}

var (
	discoveredPeers []Peer
	peersMutex      sync.RWMutex
	isDiscovering   bool
	discoveryMutex  sync.Mutex
)

const (
	serviceType = "_p2pfileshare._tcp"
	domain      = "local."
	grpcPort    = 9002
)

// StartPeerDiscovery initializes mDNS service registration and peer browsing
func StartPeerDiscovery() {
	log.Println("Starting peer discovery service...")

	// Register this device on mDNS network
	go registerService()

	// Start continuous peer browsing
	go continuousPeerBrowsing()

	log.Println("Peer discovery service started")
}

// registerService registers this device as a discoverable service
func registerService() {
	systemInfo := GetSystemInfoStruct()

	// Create TXT records with system information
	txtRecords := []string{
		"peer_id=" + systemInfo.PeerID,
		"cpu=" + systemInfo.CPU,
		"ram=" + systemInfo.RAM,
		"os=" + systemInfo.OS,
	}

	// Register the service
	server, err := zeroconf.Register(
		systemInfo.Hostname, // service instance name
		serviceType,         // service type
		domain,              // domain
		grpcPort,            // port for gRPC communication
		txtRecords,          // TXT records with metadata
		nil,                 // network interfaces (nil = all)
	)

	if err != nil {
		log.Printf("Failed to register mDNS service: %v", err)
		return
	}

	log.Printf("mDNS service registered: %s on port %d", systemInfo.Hostname, grpcPort)

	// Keep service running
	defer server.Shutdown()
	select {} // Block forever
}

// continuousPeerBrowsing continuously scans for peers every 10 seconds
func continuousPeerBrowsing() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	// Initial discovery
	discoverPeers()

	// Periodic discovery
	for range ticker.C {
		discoverPeers()
	}
}

// discoverPeers performs a single peer discovery scan
func discoverPeers() {
	discoveryMutex.Lock()
	if isDiscovering {
		discoveryMutex.Unlock()
		return
	}
	isDiscovering = true
	discoveryMutex.Unlock()

	defer func() {
		discoveryMutex.Lock()
		isDiscovering = false
		discoveryMutex.Unlock()
	}()

	log.Println("Discovering peers...")

	resolver, err := zeroconf.NewResolver(nil)
	if err != nil {
		log.Printf("Failed to initialize resolver: %v", err)
		return
	}

	entries := make(chan *zeroconf.ServiceEntry)

	// Process discovered entries
	go func(results <-chan *zeroconf.ServiceEntry) {
		tempPeers := []Peer{}

		for entry := range results {
			peer := processPeerEntry(entry)
			if peer != nil && !isOwnPeer(peer.ID) {
				tempPeers = append(tempPeers, *peer)
			}
		}

		// Update discovered peers list
		updatePeersList(tempPeers)
	}(entries)

	// Browse for services with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = resolver.Browse(ctx, serviceType, domain, entries)
	if err != nil {
		log.Printf("Failed to browse for peers: %v", err)
		return
	}

	<-ctx.Done() // Wait for timeout
	log.Printf("Discovery completed. Found %d peers", len(discoveredPeers))
}

// processPeerEntry converts a zeroconf service entry to a Peer struct
func processPeerEntry(entry *zeroconf.ServiceEntry) *Peer {
	if entry == nil || len(entry.AddrIPv4) == 0 {
		return nil
	}

	// Extract TXT record data
	txtData := parseTXTRecords(entry.Text)

	peer := &Peer{
		ID:       txtData["peer_id"],
		Hostname: entry.Instance,
		IP:       entry.AddrIPv4[0].String(),
		Port:     entry.Port,
		CPU:      txtData["cpu"],
		RAM:      txtData["ram"],
		OS:       txtData["os"],
		Status:   "online",
	}

	// Validate required fields
	if peer.ID == "" || peer.Hostname == "" {
		return nil
	}

	return peer
}

// parseTXTRecords parses TXT records into a key-value map
func parseTXTRecords(txtRecords []string) map[string]string {
	data := make(map[string]string)

	for _, record := range txtRecords {
		parts := strings.SplitN(record, "=", 2)
		if len(parts) == 2 {
			data[parts[0]] = parts[1]
		}
	}

	return data
}

// isOwnPeer checks if the discovered peer is this device itself
func isOwnPeer(peerID string) bool {
	systemInfo := GetSystemInfoStruct()
	return peerID == systemInfo.PeerID
}

// updatePeersList safely updates the discovered peers list
func updatePeersList(newPeers []Peer) {
	peersMutex.Lock()
	defer peersMutex.Unlock()

	discoveredPeers = newPeers

	// Mark offline peers (peers not seen in latest discovery)
	for i := range discoveredPeers {
		discoveredPeers[i].Status = "online"
	}
}

// GetPeers HTTP handler that returns the list of discovered peers
func GetPeers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Trigger immediate discovery if no peers found
	peersMutex.RLock()
	peerCount := len(discoveredPeers)
	peersMutex.RUnlock()

	if peerCount == 0 {
		go discoverPeers() // Non-blocking discovery
	}

	// Get current peers list
	peersMutex.RLock()
	currentPeers := make([]Peer, len(discoveredPeers))
	copy(currentPeers, discoveredPeers)
	peersMutex.RUnlock()

	response := PeersResponse{
		Peers: currentPeers,
		Count: len(currentPeers),
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding peers response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	log.Printf("Returned %d peers to client", len(currentPeers))
}

// GetPeerByID finds a peer by their ID (for internal use)
func GetPeerByID(peerID string) *Peer {
	peersMutex.RLock()
	defer peersMutex.RUnlock()

	for i := range discoveredPeers {
		if discoveredPeers[i].ID == peerID {
			return &discoveredPeers[i]
		}
	}

	return nil
}

// GetPeersCount returns the current number of discovered peers
func GetPeersCount() int {
	peersMutex.RLock()
	defer peersMutex.RUnlock()

	return len(discoveredPeers)
}
