package queue

import (
	"ghostnet/pkg/common"
	"math/rand"
	"sync"
	"time"
)

type Config struct {
	MinBatchSize    int
	MaxWaitTime     time.Duration
	RandomInterval  time.Duration
}

type MixQueue struct {
	mu         sync.Mutex
	config     Config
	packets    []*common.Packet
	out        chan []*common.Packet
	lastUpdate time.Time
}

type Manager struct {
	mu     sync.RWMutex
	queues map[common.TrafficClass]*MixQueue
	out    chan []*common.Packet
}

func NewManager(out chan []*common.Packet) *Manager {
	return &Manager{
		queues: make(map[common.TrafficClass]*MixQueue),
		out:    out,
	}
}

func (m *Manager) Enqueue(p *common.Packet) {
	m.mu.Lock()
	q, ok := m.queues[p.Type]
	if !ok {
		q = &MixQueue{
			config: Config{
				MinBatchSize:   5,
				MaxWaitTime:    500 * time.Millisecond,
				RandomInterval: 100 * time.Millisecond,
			},
			packets:    make([]*common.Packet, 0),
			out:        m.out,
			lastUpdate: time.Now(),
		}
		m.queues[p.Type] = q
		go q.run()
	}
	m.mu.Unlock()

	q.mu.Lock()
	q.packets = append(q.packets, p)
	q.mu.Unlock()
}

func (q *MixQueue) run() {
	for {
		// Randomize release interval
		jitter := time.Duration(rand.Int63n(int64(q.config.RandomInterval)))
		time.Sleep(jitter)

		q.mu.Lock()

		// Randomize batch size slightly
		dynamicMinBatch := q.config.MinBatchSize + rand.Intn(3) - 1

		shouldRelease := len(q.packets) >= dynamicMinBatch ||
			(len(q.packets) > 0 && time.Since(q.lastUpdate) >= q.config.MaxWaitTime)

		if shouldRelease {
			// Randomize order
			rand.Shuffle(len(q.packets), func(i, j int) {
				q.packets[i], q.packets[j] = q.packets[j], q.packets[i]
			})

			batch := q.packets
			q.packets = make([]*common.Packet, 0)
			q.lastUpdate = time.Now()
			q.mu.Unlock()

			q.out <- batch
		} else {
			q.mu.Unlock()
		}
	}
}
