package systemd

import (
	"reflect"
	"testing"
)

func TestParseCommandOutputSuccess(t *testing.T) {
	output := `
		ID=docker.service
		ActiveState=active
		LoadState=running
	`
	expectedMapOutput := map[string]string{
		"ID":          "docker.service",
		"ActiveState": "active",
		"LoadState":   "running",
	}

	result := parseCommandOutput(output)

	if reflect.DeepEqual(expectedMapOutput, result) {
		t.Fatalf("expected output to %v but got %v", expectedMapOutput, result)
	}
}
