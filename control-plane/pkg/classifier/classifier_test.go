package classifier

import (
	"ghostnet/control-plane/pkg/common"
	"testing"
)

func TestClassify(t *testing.T) {
	c := NewClassifier()

	tests := []struct {
		name     string
		packet   *common.Packet
		expected common.TrafficClass
	}{
		{
			name: "Small Web",
			packet: &common.Packet{
				Payload: []byte("GET /index.html HTTP/1.1"),
				Size:    400,
			},
			expected: common.ClassSmallWeb,
		},
		{
			name: "API Request",
			packet: &common.Packet{
				Payload: []byte("POST /v1/API/data"),
				Size:    1200,
			},
			expected: common.ClassAPI,
		},
		{
			name: "Video",
			packet: &common.Packet{
				Payload: []byte("VIDEOSTREAMDATA"),
				Size:    2000000,
			},
			expected: common.ClassVideo,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := c.Classify(tt.packet)
			if got != tt.expected {
				t.Errorf("Classify() = %v, want %v", got, tt.expected)
			}
		})
	}
}
