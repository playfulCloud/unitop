package store

import "testing"

func TestCreateServiceEntries(t *testing.T) {
	serviceIDs := []string{"docker.service", "bluetooth.service"}
	parameters := []string{"ID", "LoadState", "ActiveState"}

	serviceEntries := createServiceEntries(serviceIDs, parameters)

	if serviceEntries == nil {
		t.Fatal("expected map, got nil")
	}

	for _, serviceID := range serviceIDs {
		entry, exists := serviceEntries[serviceID]
		if !exists {
			t.Fatalf("expected entries to contain %s, but did not", serviceID)
		}

		for _, parameter := range parameters {
			parameterValue, exists := entry.Params[parameter]
			if !exists {
				t.Fatalf("expected entry %s to contain key %s", serviceID, parameter)
			}

			if parameterValue != "" {
				t.Fatalf("expected parameter value to be blank, got %s", parameterValue)
			}
		}
	}
}

func TestUpdateServiceEntrySuccess(t *testing.T) {
	serviceStore := newTestServiceStore()

	states := map[string]string{
		"ID":          "docker.service",
		"LoadState":   "active",
		"ActiveState": "running",
	}

	success := serviceStore.UpdateServiceEntry("docker.service", states)

	if !success {
		t.Fatalf("expected operation to succeed, got %t", success)
	}

	entry, exists := serviceStore.GetServiceEntry("docker.service")
	if !exists {
		t.Fatal("expected entry to exist")
	}

	for key, expectedValue := range states {
		actualValue := entry.Params[key]

		if actualValue != expectedValue {
			t.Fatalf("expected value of %s to be %s, got %s", key, expectedValue, actualValue)
		}
	}
}

func TestUpdateServiceEntryFail(t *testing.T) {
	serviceStore := newTestServiceStore()

	states := map[string]string{
		"ID":          "docker.service",
		"LoadState":   "active",
		"ActiveState": "running",
	}

	success := serviceStore.UpdateServiceEntry("non-existing.service", states)

	if success {
		t.Fatalf("expected operation to fail, got %t", success)
	}
}

func newTestServiceStore() *ServiceStore {
	serviceIDs := []string{"docker.service", "bluetooth.service"}
	parameters := []string{"ID", "LoadState", "ActiveState"}

	return NewServiceStore(serviceIDs, parameters)
}

func TestGetServiceEntrySuccess(t *testing.T) {
	serviceStore := newTestServiceStore()

	entry, exists := serviceStore.GetServiceEntry("docker.service")

	if !exists {
		t.Fatal("expected entry to exist")
	}

	if entry == nil {
		t.Fatal("expected entry, got nil")
	}

	if entry.ID != "docker.service" {
		t.Fatalf("expected ID docker.service, got %s", entry.ID)
	}
}

func TestGetServiceEntryNotFound(t *testing.T) {
	serviceStore := newTestServiceStore()

	entry, exists := serviceStore.GetServiceEntry("non-existing.service")

	if exists {
		t.Fatal("expected entry not to exist")
	}

	if entry != nil {
		t.Fatalf("expected nil entry, got %+v", entry)
	}
}

func TestGetServiceEntryReturnsClone(t *testing.T) {
	serviceStore := newTestServiceStore()

	entry, exists := serviceStore.GetServiceEntry("docker.service")
	if !exists {
		t.Fatal("expected entry to exist")
	}

	entry.Params["ActiveState"] = "modified"

	entryAgain, exists := serviceStore.GetServiceEntry("docker.service")
	if !exists {
		t.Fatal("expected entry to exist")
	}

	if entryAgain.Params["ActiveState"] == "modified" {
		t.Fatal("expected store to be protected from external mutation")
	}
}

func TestGetServiceEntriesReturnsClones(t *testing.T) {
	serviceStore := newTestServiceStore()

	entries := serviceStore.GetServiceEntries()

	entries["docker.service"].Params["ActiveState"] = "modified"

	entryAgain, exists := serviceStore.GetServiceEntry("docker.service")
	if !exists {
		t.Fatal("expected entry to exist")
	}

	if entryAgain.Params["ActiveState"] == "modified" {
		t.Fatal("expected store to be protected from external mutation")
	}
}

func TestServiceStoreConcurrentAccess(t *testing.T) {
	serviceStore := newTestServiceStore()

	done := make(chan struct{})

	for range 100 {
		go func() {
			for range 1000 {
				serviceStore.UpdateServiceEntry("docker.service", map[string]string{
					"ID":          "docker.service",
					"LoadState":   "loaded",
					"ActiveState": "active",
				})
			}

			done <- struct{}{}
		}()

		go func() {
			for range 1000 {
				serviceStore.GetServiceEntry("docker.service")
				serviceStore.GetServiceEntries()
			}

			done <- struct{}{}
		}()
	}

	for range 200 {
		<-done
	}
}
