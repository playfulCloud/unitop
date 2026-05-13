package cmdclient

import (
	"errors"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/playfulCloud/unitop/internal/model"
)

func TestExecuteWithTimeoutSuccess(t *testing.T) {
	output, err := ExecuteWithTimeout(
		*model.NewCommand("sh", []string{"-c", "printf unitop"}),
		time.Second,
	)
	if err != nil {
		t.Fatalf("expected command to succeed, got %v", err)
	}

	if output != "unitop" {
		t.Fatalf("expected output unitop, got %q", output)
	}
}

func TestExecuteWithTimeoutTimesOut(t *testing.T) {
	_, err := ExecuteWithTimeout(
		*model.NewCommand("sh", []string{"-c", "sleep 1"}),
		10*time.Millisecond,
	)
	if err == nil {
		t.Fatal("expected timeout error, got nil")
	}

	if !strings.Contains(err.Error(), "command timed out") {
		t.Fatalf("expected timeout error, got %v", err)
	}
}

func TestExecuteWithTimeoutReturnsCommandFailure(t *testing.T) {
	_, err := ExecuteWithTimeout(
		*model.NewCommand("sh", []string{"-c", "exit 7"}),
		time.Second,
	)
	if err == nil {
		t.Fatal("expected command failure, got nil")
	}

	var exitErr *exec.ExitError
	if !errors.As(err, &exitErr) {
		t.Fatalf("expected wrapped exit error, got %v", err)
	}
}
