package cmdclient

import (
	"context"
	"fmt"
	"os/exec"
	"time"

	"github.com/playfulCloud/unitop/internal/model"
)

func Execute(command model.Command) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, command.Name, command.Args...)

	output, err := cmd.CombinedOutput()
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return string(output), fmt.Errorf(
				"command timed out: %s %v",
				command.Name,
				command.Args,
			)
		}

		return string(output), fmt.Errorf(
			"command failed: %s %v: %w: %s",
			command.Name,
			command.Args,
			err,
			string(output),
		)
	}

	return string(output), nil
}
