package systemd

import "strings"

func parseCommandOutput(commandOutput string) map[string]string {
	lines := strings.SplitSeq(commandOutput, "\n")
	properties := make(map[string]string)

	for line := range lines {
		keyValueParts := strings.SplitN(line, "=", 2)
		if len(keyValueParts) == 2 {
			properties[keyValueParts[0]] = keyValueParts[1]
		}
	}
	return properties
}
