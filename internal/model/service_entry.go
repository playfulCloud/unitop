package model

import "maps"

type ServiceEntry struct {
	ID     string
	Params map[string]string
}

func NewServiceEntry(
	id string) *ServiceEntry {

	return &ServiceEntry{
		ID:     id,
		Params: make(map[string]string),
	}
}

func (s *ServiceEntry) Clone() *ServiceEntry {
	paramsCopy := make(map[string]string, len(s.Params))

	maps.Copy(paramsCopy, s.Params)

	return &ServiceEntry{
		ID:     s.ID,
		Params: paramsCopy,
	}
}
