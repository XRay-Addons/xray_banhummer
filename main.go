package main

import (
	"flag"
	"goipban/banning_loop"
	"goipban/config"
	"log"
)

func main() {
	log.Printf("Tool to ban XRay users running.")

	// parse config
	configPath := getConfigPath()
	config, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("Could not read config file: %s", err)
		return
	}

	// log params
	log.Printf("Update Ban list interval: %s\n", config.UpdateInterval)
	log.Printf("Ban time: %s\n", config.BanDuration)
	log.Printf("Banning xray outbound: %s\n", config.XRayBlacklistOutbound)

	// start main loop
	mainLoop := banning_loop.CreateBanningLoop(*config)
	if err := mainLoop.StartBanningLoop(); err != nil {
		log.Fatalf("banning loop error %s", err)
		return
	}
}

func getConfigPath() string {
	configPath := flag.String("config", "", "Path to the JSON configuration file")
	flag.Parse()
	return *configPath
}
