package tui

import (
	"strings"

	"github.com/playfulCloud/unitop/internal/systemd"
)

type ActionBinding struct {
	Key         string
	Description string
	Action      systemd.ServiceAction
}

var actionBindings = []ActionBinding{
	{Key: "r", Description: "restart", Action: systemd.RestartAction},
	{Key: "s", Description: "start", Action: systemd.StartAction},
	{Key: "x", Description: "stop", Action: systemd.StopAction},
	{Key: "e", Description: "enable", Action: systemd.EnableAction},
	{Key: "d", Description: "disable", Action: systemd.DisableAction},
}

func actionForKey(key string) (systemd.ServiceAction, bool) {
	for _, binding := range actionBindings {
		if binding.Key == key {
			return binding.Action, true
		}
	}

	return "", false
}

func footerText() string {
	parts := []string{
		"↑↓ move",
		"/ filter",
		"l logs",
		"r restart",
		"s start",
		"x stop",
		"e/d enable/disable",
		"q quit",
	}

	return strings.Join(parts, " | ")
}
