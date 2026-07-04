package pipeline

import (
	"testing"
	"time"
)

func TestPipeline(t *testing.T) {
	p := NewPipeline()
	p.CoverEngine.Start()

	site := "example.com"
	data := []byte("GET / HTTP/1.1\r\nHost: example.com\r\n\r\n")

	// Process a few requests
	for i := 0; i < 10; i++ {
		p.ProcessRequest(site, data)
	}

	// Expect some packets in FinalOut
	select {
	case pkt := <-p.FinalOut:
		if pkt == nil {
			t.Error("Received nil packet")
		}
		if len(pkt.Metadata.Route) == 0 {
			t.Error("Packet has no route")
		}
	case <-time.After(2 * time.Second):
		t.Error("No packets reached final output")
	}
}
