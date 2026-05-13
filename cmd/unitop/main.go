package main

import (
	"log"
	"time"

	tea "github.com/charmbracelet/bubbletea"
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
	store := store.NewServiceStore(cfg.ServiceNames, cfg.Properties)
	c := systemd.NewSystemdManager(store, cfg.Properties)
	c.MonitorState()

	p := tea.NewProgram(
		tui.NewModel(c, 5*time.Second),
		tea.WithAltScreen(),
	)

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}

}
