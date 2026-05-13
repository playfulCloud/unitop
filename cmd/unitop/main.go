package main

import (
	"log"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/playfulCloud/unitop/internal/cmdclient"
	"github.com/playfulCloud/unitop/internal/config"
	"github.com/playfulCloud/unitop/internal/store"
	"github.com/playfulCloud/unitop/internal/systemd"
	"github.com/playfulCloud/unitop/internal/tui"
)

func main() {
	cfg, err := config.ReadConfig("configs/unitop.yaml")
	if err != nil {
		log.Fatal(err)
	}

	serviceNames := cfg.ServiceNames
	if strings.EqualFold(cfg.Mode, "all") {
		serviceNames, err = systemd.DiscoverServiceNames(cfg.Discovery, cmdclient.Execute)
		if err != nil {
			log.Fatal(err)
		}
	}

	refreshInterval := 5 * time.Second
	if cfg.RefreshInterval != "" {
		refreshInterval, err = time.ParseDuration(cfg.RefreshInterval)
		if err != nil {
			log.Fatal(err)
		}
	}

	store := store.NewServiceStore(serviceNames, cfg.Properties)
	c := systemd.NewSystemdManager(store, cfg.Properties)
	c.MonitorState()

	p := tea.NewProgram(
		tui.NewModel(c, refreshInterval),
		tea.WithAltScreen(),
	)

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}

}
