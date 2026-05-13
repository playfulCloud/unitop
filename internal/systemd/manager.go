package systemd

import (
	"fmt"
	"sync"

	"github.com/playfulCloud/unitop/internal/cmdclient"
	"github.com/playfulCloud/unitop/internal/model"
	"github.com/playfulCloud/unitop/internal/store"
)

type SystemdManager struct {
	Store         *store.ServiceStore
	Properties    []string
	Execute       ExecuteFunc
	ActionExecute ExecuteFunc
}

type ExecuteFunc func(command model.Command) (string, error)

type ServiceAction string

const (
	RestartAction ServiceAction = "restart"
	StartAction   ServiceAction = "start"
	StopAction    ServiceAction = "stop"
	EnableAction  ServiceAction = "enable"
	DisableAction ServiceAction = "disable"
)

func NewSystemdManager(store *store.ServiceStore, properties []string) *SystemdManager {
	return &SystemdManager{
		Store:         store,
		Properties:    properties,
		Execute:       cmdclient.Execute,
		ActionExecute: cmdclient.ExecuteAction,
	}
}

func (m *SystemdManager) MonitorState() error {
	entries := m.Store.GetServiceEntries()

	var wg sync.WaitGroup

	sem := make(chan struct{}, 10)

	for serviceID := range entries {
		wg.Add(1)

		go func(serviceID string) {
			defer wg.Done()

			sem <- struct{}{}
			defer func() { <-sem }()

			command := BuildSystemctlShowWithArgs(serviceID, m.Properties)

			commandOutput, err := m.Execute(*command)
			if err != nil {
				m.Store.UpdateServiceEntry(serviceID, failedServiceParams(serviceID, m.Properties))
				return
			}

			updatedProperties := parseCommandOutput(commandOutput)

			m.Store.UpdateServiceEntry(serviceID, updatedProperties)
		}(serviceID)
	}

	wg.Wait()

	return nil
}

func (m *SystemdManager) ExecuteAction(
	serviceID string,
	action ServiceAction,
) error {

	command := BuildSystemctlActionCommand(
		serviceID,
		string(action),
	)

	output, err := m.ActionExecute(*command)
	if err != nil {
		return fmt.Errorf(
			"command failed: %s %v: %w: %s",
			command.Name,
			command.Args,
			err,
			string(output),
		)
	}

	return nil
}

func failedServiceParams(serviceID string, properties []string) map[string]string {
	params := make(map[string]string, len(properties))
	for _, property := range properties {
		params[property] = ""
	}

	params["Id"] = serviceID
	params["ID"] = serviceID
	params["LoadState"] = "error"
	params["ActiveState"] = "unknown"
	params["SubState"] = "show failed"

	return params
}
