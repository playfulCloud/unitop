package systemd

import (
	"fmt"
	"time"

	"github.com/playfulCloud/unitop/internal/cmdclient"
	"github.com/playfulCloud/unitop/internal/store"
)

type SystemdManager struct {
	Store      store.ServiceStore
	Properties []string
}

type ServiceAction string

const (
	RestartAction ServiceAction = "restart"
	StartAction   ServiceAction = "start"
	StopAction    ServiceAction = "stop"
	EnableAction  ServiceAction = "enable"
	DisableAction ServiceAction = "disable"
)

func NewSystemdManager(store store.ServiceStore, properties []string) *SystemdManager {
	return &SystemdManager{
		Store:      store,
		Properties: properties,
	}
}

func (c *SystemdManager) MonitorState() error {
	entries := c.Store.GetServiceEntries()

	for key := range entries {
		command := BuildSystemctlShowWithArgs(key, c.Properties)
		commandOutput, err := cmdclient.Execute(*command)
		if err != nil {
			return fmt.Errorf("Error while executing command: %w", err)
		}
		updatedProperties := parseCommandOutput(commandOutput)
		c.Store.UpdateServiceEntry(key, updatedProperties)
	}
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

	output, err := cmdclient.Execute(*command)
	if err != nil {
		return fmt.Errorf(
			"command failed: %s %v: %w: %s",
			command.Name,
			command.Args,
			err,
			string(output),
		)
	}
	time.Sleep(10 * time.Second)

	return nil
}
