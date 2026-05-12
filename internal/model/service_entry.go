package model

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
