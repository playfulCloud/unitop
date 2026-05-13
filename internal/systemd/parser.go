package systemd

import "strings"

func parseCommandOutput(output string) map[string]string {
	result := make(map[string]string)

	lines := strings.SplitSeq(output, "\n")

	for line := range lines {
		line = strings.TrimSpace(line)

		if line == "" {
			continue
		}

		parts := strings.SplitN(line, "=", 2)

		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		result[key] = value
	}

	return result
}
