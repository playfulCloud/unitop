package systemd

import (
	"github.com/playfulCloud/unitop/internal/model"
)

func BuildSystemctlShowWithArgs(serviceID string, properties []string) *model.Command {

	args := []string{"show", serviceID}

	for _, property := range properties {
		args = append(args, "--property="+property)
	}

	return model.NewCommand("systemctl", args)
}

func BuildSystemctlActionCommand(serviceID string, action string) *model.Command {
	args := []string{action, serviceID}
	return model.NewCommand("systemctl", args)
}
