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
