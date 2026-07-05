package queue

import (
	"fmt"
	"ghostnet/control-plane/pkg/common"
)

type Stats struct {
	TotalPackets int
	TotalBatches int
	QueueLengths map[common.TrafficClass]int
}

func (m *Manager) GetStats() Stats {
	m.mu.RLock()
	defer m.mu.RUnlock()

	stats := Stats{
		QueueLengths: make(map[common.TrafficClass]int),
	}

	for class, q := range m.queues {
		q.mu.Lock()
		stats.QueueLengths[class] = len(q.packets)
		q.mu.Unlock()
	}

	return stats
}

func (s Stats) String() string {
	return fmt.Sprintf("Queues: %v", s.QueueLengths)
}
