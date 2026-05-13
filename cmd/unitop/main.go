package main

import (
	"flag"
	"fmt"
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

var version = "dev"

func main() {
	configPath := flag.String("config", "", "path to config file")
	showVersion := flag.Bool("version", false, "print version and exit")
	flag.Parse()

	if *showVersion {
		fmt.Printf("unitop %s\n", version)
		return
	}

	if *configPath == "" {
		defaultConfigPath, err := config.DefaultConfigPath()
		if err != nil {
			log.Fatal(err)
		}

		*configPath = defaultConfigPath
	}

	cfg, err := config.ReadOrCreateConfig(*configPath)
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

		if refreshInterval <= 0 {
			log.Fatal("refresh_interval must be greater than zero")
		}
	}

	store := store.NewServiceStore(serviceNames, systemd.DefaultProperties)
	c := systemd.NewSystemdManager(store, systemd.DefaultProperties)
	c.MonitorState()

	p := tea.NewProgram(
		tui.NewModel(c, refreshInterval),
		tea.WithAltScreen(),
	)

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}

}
