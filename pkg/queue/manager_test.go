package queue

import (
	"ghostnet/pkg/common"
	"testing"
	"time"
)

func TestMixQueueBatching(t *testing.T) {
	out := make(chan []*common.Packet, 10)
	mgr := NewManager(out)

	// Send 3 packets (less than MinBatchSize=5)
	for i := 0; i < 3; i++ {
		mgr.Enqueue(&common.Packet{Type: common.ClassSmallWeb, ID: "test"})
	}

	select {
	case <-out:
		t.Error("Batch released too early")
	case <-time.After(100 * time.Millisecond):
		// Success
	}

	// Send 2 more packets to reach MinBatchSize
	for i := 0; i < 2; i++ {
		mgr.Enqueue(&common.Packet{Type: common.ClassSmallWeb, ID: "test"})
	}

	select {
	case batch := <-out:
		if len(batch) != 5 {
			t.Errorf("Expected batch size 5, got %d", len(batch))
		}
	case <-time.After(1 * time.Second):
		t.Error("Batch not released in time")
	}
}
