package systemd

import (
	"fmt"

	"github.com/playfulCloud/unitop/internal/cmdclient"
	"github.com/playfulCloud/unitop/internal/store"
)

type Collector struct {
	Store      store.ServiceStore
	Properties []string
}

func NewCollector(store store.ServiceStore, properties []string) *Collector {
	return &Collector{
		Store:      store,
		Properties: properties,
	}
}

func (c *Collector) MonitorState() error {
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
