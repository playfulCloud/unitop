package systemd

import (
	"path/filepath"
	"sort"
	"strings"

	"github.com/playfulCloud/unitop/internal/config"
)

var defaultDiscoveryStates = []string{
	"enabled",
	"enabled-runtime",
	"linked",
	"linked-runtime",
}

func DiscoverServiceNames(
	discovery config.DiscoveryConfig,
	execute ExecuteFunc,
) ([]string, error) {
	command := BuildSystemctlListUnitFilesCommand()

	output, err := execute(*command)
	if err != nil {
		return nil, err
	}

	return parseDiscoveredServices(output, discovery), nil
}

func parseDiscoveredServices(output string, discovery config.DiscoveryConfig) []string {
	states := discovery.States
	if len(states) == 0 {
		states = defaultDiscoveryStates
	}

	include := discovery.Include
	if len(include) == 0 {
		include = []string{"*.service"}
	}

	stateSet := makeSet(states)
	serviceSet := make(map[string]struct{})

	for _, line := range strings.Split(output, "\n") {
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}

		serviceID := fields[0]
		state := fields[1]

		if strings.Contains(serviceID, "@.") {
			continue
		}

		if _, ok := stateSet[state]; !ok {
			continue
		}

		if !matchesAny(serviceID, include) {
			continue
		}

		if matchesAny(serviceID, discovery.Exclude) {
			continue
		}

		serviceSet[serviceID] = struct{}{}
	}

	services := make([]string, 0, len(serviceSet))
	for serviceID := range serviceSet {
		services = append(services, serviceID)
	}

	sort.Strings(services)

	return services
}

func makeSet(values []string) map[string]struct{} {
	set := make(map[string]struct{}, len(values))
	for _, value := range values {
		set[value] = struct{}{}
	}

	return set
}

func matchesAny(value string, patterns []string) bool {
	for _, pattern := range patterns {
		matched, _ := filepath.Match(pattern, value)
		if matched {
			return true
		}
	}

	return false
}
