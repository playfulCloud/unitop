package systemd

import (
	"reflect"
	"testing"

	"github.com/playfulCloud/unitop/internal/config"
	"github.com/playfulCloud/unitop/internal/model"
)

func TestDiscoverServiceNamesFiltersToMeaningfulServicesByDefault(t *testing.T) {
	discovery := config.DiscoveryConfig{
		Include: []string{"*.service"},
		Exclude: []string{"systemd-*", "user@*.service"},
	}

	output := `
docker.service enabled enabled
postgresql.service disabled enabled
systemd-journald.service static -
user@1000.service enabled enabled
dbus.service static -
tmp.mount generated -
masked-app.service masked enabled
serial-getty@.service enabled enabled
`

	services, err := DiscoverServiceNames(discovery, func(command model.Command) (string, error) {
		return output, nil
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expected := []string{
		"docker.service",
	}

	if !reflect.DeepEqual(expected, services) {
		t.Fatalf("expected services %v, got %v", expected, services)
	}
}

func TestDiscoverServiceNamesHonorsConfiguredStates(t *testing.T) {
	discovery := config.DiscoveryConfig{
		States: []string{"static"},
	}

	output := `
dbus.service static -
docker.service enabled enabled
`

	services, err := DiscoverServiceNames(discovery, func(command model.Command) (string, error) {
		return output, nil
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expected := []string{"dbus.service"}

	if !reflect.DeepEqual(expected, services) {
		t.Fatalf("expected services %v, got %v", expected, services)
	}
}
