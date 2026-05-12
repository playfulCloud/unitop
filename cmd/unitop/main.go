package main

import (
	"fmt"
	"log"

	"github.com/playfulCloud/unitop/internal/config"
)

func main() {
	cfg, err := config.ReadConfig("configs/unitop.yaml")
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	fmt.Println(cfg)
}
