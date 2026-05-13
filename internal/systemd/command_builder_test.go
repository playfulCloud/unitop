package systemd

import (
	"reflect"
	"testing"

	"github.com/playfulCloud/unitop/internal/model"
)

func TestBuildSystemctlShowWithArgsSuccess(t *testing.T) {
	serviceID := "docker.service"
	properties := []string{"ID", "LoadState", "ActiveState"}

	command := BuildSystemctlShowWithArgs(serviceID, properties)
	expectedCommand := createExpectedCommand()

	if !reflect.DeepEqual(expectedCommand, command) {
		t.Fatalf("expected %v but got %v", expectedCommand, command)
	}
}

func TestBuildSystemctlListUnitFilesCommand(t *testing.T) {
	command := BuildSystemctlListUnitFilesCommand()
	expectedArgs := []string{
		"list-unit-files",
		"--type=service",
		"--no-legend",
		"--no-pager",
	}

	if command.Name != "systemctl" {
		t.Fatalf("expected command name systemctl, got %s", command.Name)
	}

	if !reflect.DeepEqual(expectedArgs, command.Args) {
		t.Fatalf("expected args %v, got %v", expectedArgs, command.Args)
	}
}

func createExpectedCommand() *model.Command {
	commandName := "systemctl"
	args := []string{
		"show",
		"docker.service",
		"--property=ID",
		"--property=LoadState",
		"--property=ActiveState",
	}
	return model.NewCommand(commandName, args)
}
