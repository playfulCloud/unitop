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

func BuildSystemctlListUnitFilesCommand() *model.Command {
	args := []string{
		"list-unit-files",
		"--type=service",
		"--no-legend",
		"--no-pager",
	}

	return model.NewCommand("systemctl", args)
}

func BuildSystemctlActionCommand(serviceID string, action string) *model.Command {
	args := []string{action, serviceID}
	return model.NewCommand("systemctl", args)
}

func BuildJournalctlCommand(serviceID string) *model.Command {
	args := []string{"-u", serviceID, "-n", "200"}
	return model.NewCommand("journalctl", args)
}
