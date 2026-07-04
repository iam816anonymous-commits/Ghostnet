package cover

import (
	"crypto/rand"
	"ghostnet/pkg/common"
	"sync"
	"time"
)

type Engine struct {
	mu          sync.Mutex
	mode        common.PrivacyMode
	out         chan *common.Packet
	stop        chan struct{}
	activeUsers int
}

func NewEngine(out chan *common.Packet) *Engine {
	return &Engine{
		mode: common.ModeMedium,
		out:  out,
		stop: make(chan struct{}),
	}
}

func (e *Engine) SetMode(mode common.PrivacyMode) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.mode = mode
}

func (e *Engine) SetActiveUsers(count int) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.activeUsers = count
}

func (e *Engine) GetOutChan() chan *common.Packet {
	return e.out
}

func (e *Engine) Start() {
	go func() {
		for {
			select {
			case <-e.stop:
				return
			case <-time.After(e.getInterval()):
				if e.mode == common.ModeOff {
					continue
				}
				e.out <- e.generateDummy()
			}
		}
	}()
}

func (e *Engine) getInterval() time.Duration {
	e.mu.Lock()
	defer e.mu.Unlock()

	base := 1 * time.Second
	switch e.mode {
	case common.ModeLow:
		base = 2 * time.Second
	case common.ModeMedium:
		base = 1 * time.Second
	case common.ModeHigh:
		base = 500 * time.Millisecond
	case common.ModeMaximumPrivacy:
		base = 100 * time.Millisecond
	}

	// Adapt based on active users - more users means we can send less dummy traffic
	// as the anonymity set is already large.
	if e.activeUsers > 1000 {
		base *= 2
	}

	return base
}

func (e *Engine) generateDummy() *common.Packet {
	payload := make([]byte, 512)
	_, _ = rand.Read(payload)

	return &common.Packet{
		ID:        "dummy",
		Payload:   payload,
		Size:      len(payload),
		Type:      common.ClassBackground,
		IsDummy:   true,
		Timestamp: time.Now(),
	}
}
