package normalizer

import (
	"ghostnet/control-plane/pkg/common"
	"testing"
)

func TestNormalize(t *testing.T) {
	n := NewNormalizer()

	t.Run("Padding", func(t *testing.T) {
		p := &common.Packet{Payload: []byte("short"), Size: 5}
		packets := n.Normalize(p)
		if len(packets) != 1 {
			t.Errorf("Expected 1 packet, got %d", len(packets))
		}
		if len(packets[0].Payload) != FixedPacketSize {
			t.Errorf("Expected size %d, got %d", FixedPacketSize, len(packets[0].Payload))
		}
	})

	t.Run("Fragmentation", func(t *testing.T) {
		largePayload := make([]byte, 2500)
		p := &common.Packet{Payload: largePayload, Size: 2500}
		packets := n.Normalize(p)
		if len(packets) != 3 {
			t.Errorf("Expected 3 packets, got %d", len(packets))
		}
		for _, frag := range packets {
			if len(frag.Payload) != FixedPacketSize {
				t.Errorf("Expected fragment size %d, got %d", FixedPacketSize, len(frag.Payload))
			}
		}
	})
}
