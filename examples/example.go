package main

import (
	"fmt"
	"ghostnet/control-plane/pkg/common"
	"ghostnet/control-plane/pkg/pipeline"
	"time"
)

func main() {
	// 1. Initialize the GhostNet pipeline
	p := pipeline.NewPipeline()

	// 2. Configure privacy mode (LOW, MEDIUM, HIGH, MAXIMUM PRIVACY)
	p.CoverEngine.SetMode(common.ModeHigh)
	p.CoverEngine.Start()

	fmt.Println("GhostNet Pipeline initialized...")

	// 3. Process some simulated application traffic
	site := "example.com"
	data := []byte("GET /index.html HTTP/1.1\r\nHost: example.com\r\n\r\n")

	fmt.Printf("Sending request to %s via GhostNet...\n", site)
	p.ProcessRequest(site, data)

	// 4. Listen for processed packets arriving at the end of the pipeline
	// In a real app, these would be sent over the network to the first relay.
	select {
	case pkt := <-p.FinalOut:
		fmt.Printf("Packet processed and ready for transit!\n")
		fmt.Printf(" - ID: %s\n", pkt.ID)
		fmt.Printf(" - Size: %d bytes (Normalized)\n", pkt.Size)
		fmt.Printf(" - Route: %v (Multi-hop)\n", pkt.Metadata.Route)
		fmt.Printf(" - Traffic Class: %s\n", pkt.Type)
	case <-time.After(2 * time.Second):
		fmt.Println("Timed out waiting for packet to clear the mix queue.")
	}
}
