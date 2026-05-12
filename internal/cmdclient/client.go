package cmdclient

import (
	"github.com/playfulCloud/unitop/internal/model"
	"os/exec"
)

func Execute(command model.Command) (string, error) {
	byteOutput, err := exec.Command(
		command.Name,
		command.Args...,
	).CombinedOutput()

	return string(byteOutput), err
}
