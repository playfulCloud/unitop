package store

import "github.com/playfulCloud/unitop/internal/model"

type ServiceStore struct {
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
	entry, exists := s.entries[serviceID]
	if !exists {
		return false
	}
	entry.Params = params

	return true
}

func (s *ServiceStore) GetServiceEntry(serviceID string) (*model.ServiceEntry, bool) {
	entry, exists := s.entries[serviceID]
	return entry, exists
}
