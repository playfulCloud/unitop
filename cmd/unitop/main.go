package main

import (
	"fmt"
	"log"

	"github.com/playfulCloud/unitop/internal/cmdclient"
	"github.com/playfulCloud/unitop/internal/config"
	"github.com/playfulCloud/unitop/internal/store"
	"github.com/playfulCloud/unitop/internal/systemd"
)

func main() {
	cfg, err := config.ReadConfig("configs/unitop.yaml")
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	store := store.NewServiceStore(cfg.ServiceNames, cfg.Properties)
	command := systemd.BuildSystemctlShowWithArgs("docker.service", cfg.Properties)
	result, err := cmdclient.Execute(*command)
	fmt.Println(cfg)
	fmt.Println(store)
	fmt.Println(result)
}
