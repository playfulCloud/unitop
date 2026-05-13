package store

import (
	"sync"

	"github.com/playfulCloud/unitop/internal/model"
)

type ServiceStore struct {
	mu      sync.RWMutex
	entries map[string]*model.ServiceEntry
}

func NewServiceStore(serviceIDs []string, parameters []string) *ServiceStore {
	return &ServiceStore{
		entries: createServiceEntries(serviceIDs, parameters),
	}
}

func createServiceEntries(serviceIDs []string, parameters []string) map[string]*model.ServiceEntry {
	serviceEntries := make(map[string]*model.ServiceEntry, len(serviceIDs))

	for _, serviceID := range serviceIDs {
		serviceEntry := model.NewServiceEntry(serviceID)

		for _, parameter := range parameters {
			serviceEntry.Params[parameter] = ""
		}

		serviceEntries[serviceID] = serviceEntry
	}

	return serviceEntries
}

func (s *ServiceStore) UpdateServiceEntry(serviceID string, params map[string]string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry, exists := s.entries[serviceID]
	if !exists {
		return false
	}

	entry.Params = params
	return true
}

func (s *ServiceStore) GetServiceEntry(serviceID string) (*model.ServiceEntry, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	entry, exists := s.entries[serviceID]
	if !exists {
		return nil, false
	}

	return entry.Clone(), true
}

func (s *ServiceStore) GetServiceEntries() map[string]*model.ServiceEntry {
	s.mu.RLock()
	defer s.mu.RUnlock()

	entriesCopy := make(map[string]*model.ServiceEntry, len(s.entries))

	for serviceID, entry := range s.entries {
		entriesCopy[serviceID] = entry.Clone()
	}

	return entriesCopy
}
