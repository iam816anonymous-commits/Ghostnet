package normalizer

import (
	"fmt"
	"ghostnet/pkg/common"
)

const (
	FixedPacketSize = 1024
)

type Normalizer struct{}

func NewNormalizer() *Normalizer {
	return &Normalizer{}
}

func (n *Normalizer) Normalize(p *common.Packet) []*common.Packet {
	if p.Size == FixedPacketSize {
		return []*common.Packet{p}
	}

	if p.Size < FixedPacketSize {
		// Pad small packets
		paddedPayload := make([]byte, FixedPacketSize)
		copy(paddedPayload, p.Payload)
		p.Payload = paddedPayload
		p.Size = FixedPacketSize
		return []*common.Packet{p}
	}

	// Split large packets
	var fragments []*common.Packet
	for i := 0; i < p.Size; i += FixedPacketSize {
		end := i + FixedPacketSize
		if end > p.Size {
			end = p.Size
		}

		fragPayload := make([]byte, FixedPacketSize)
		copy(fragPayload, p.Payload[i:end])

		fragments = append(fragments, &common.Packet{
			ID:        fmt.Sprintf("%s-frag-%d", p.ID, i/FixedPacketSize),
			Payload:   fragPayload,
			Size:      FixedPacketSize,
			Type:      p.Type,
			IsDummy:   p.IsDummy,
			Metadata:  p.Metadata,
			Timestamp: p.Timestamp,
		})
	}
	return fragments
}
