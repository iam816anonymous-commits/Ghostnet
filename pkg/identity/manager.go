package identity

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"sync"
)

type Identity struct {
	ID           string
	Site         string
	Cookies      map[string]string
	LocalStorage map[string]string
	Cache        []byte
}

type Manager struct {
	mu         sync.RWMutex
	identities map[string]*Identity // site -> Identity
}

func NewManager() *Manager {
	return &Manager{
		identities: make(map[string]*Identity),
	}
}

func (m *Manager) GetIdentity(site string) *Identity {
	m.mu.Lock()
	defer m.mu.Unlock()

	if id, ok := m.identities[site]; ok {
		return id
	}

	id := &Identity{
		ID:           generateID(),
		Site:         site,
		Cookies:      make(map[string]string),
		LocalStorage: make(map[string]string),
	}
	m.identities[site] = id
	return id
}

func (m *Manager) RotateIdentity(site string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.identities, site)
}

func generateID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

func (m *Manager) GetStatus() string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return fmt.Sprintf("Active Identities: %d", len(m.identities))
}
