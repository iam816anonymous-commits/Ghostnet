package routing

import (
	"fmt"
	"math/rand"
	"sync"
)

type RelayRole string

const (
	RoleEntry     RelayRole = "Entry"
	RoleMiddle    RelayRole = "Middle"
	RoleExit      RelayRole = "Exit"
	RoleBridge    RelayRole = "Bridge"
	RoleDirectory RelayRole = "Directory"
)

type Relay struct {
	ID         string
	IP         string
	Country    string
	Role       RelayRole
	Reputation float64
	Load       float64
}

type Engine struct {
	mu     sync.RWMutex
	relays map[string]*Relay
}

func NewEngine() *Engine {
	e := &Engine{
		relays: make(map[string]*Relay),
	}
	e.seedRelays()
	return e
}

func (e *Engine) seedRelays() {
	countries := []string{"Germany", "Singapore", "Canada", "USA", "India", "UK", "Japan"}
	for i := 0; i < 50; i++ {
		role := RoleMiddle
		if i < 10 {
			role = RoleEntry
		} else if i > 40 {
			role = RoleExit
		}

		id := fmt.Sprintf("relay-%d", i)
		e.relays[id] = &Relay{
			ID:         id,
			IP:         fmt.Sprintf("1.2.3.%d", i),
			Country:    countries[rand.Intn(len(countries))],
			Role:       role,
			Reputation: 1.0,
			Load:       rand.Float64(),
		}
	}
}

func (e *Engine) SelectRoute() []string {
	e.mu.Lock()
	defer e.mu.Unlock()

	var entries, middles, exits []*Relay
	for _, r := range e.relays {
		// Basic load balancing and reputation filtering
		if r.Load > 0.9 || r.Reputation < 0.5 {
			continue
		}

		switch r.Role {
		case RoleEntry:
			entries = append(entries, r)
		case RoleMiddle:
			middles = append(middles, r)
		case RoleExit:
			exits = append(exits, r)
		}
	}

	if len(entries) == 0 || len(middles) < 2 || len(exits) == 0 {
		return nil
	}

	// Pick relays and increment load
	e1 := entries[rand.Intn(len(entries))]
	m1 := middles[rand.Intn(len(middles))]
	m2 := middles[rand.Intn(len(middles))]
	ex := exits[rand.Intn(len(exits))]

	e1.Load += 0.01
	m1.Load += 0.01
	m2.Load += 0.01
	ex.Load += 0.01

	return []string{e1.ID, m1.ID, m2.ID, ex.ID}
}

func (e *Engine) GetRelay(id string) *Relay {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.relays[id]
}
