package main

import (
	"flag"

	"./src/config"
	"./src/data"
	"./src/log"
)

var cfg Config

type Config struct {
	Regions []string `json:regions`
}

func main() {
	var (
		flagLogLevel int
		flagRegions  string
	)

	flag.IntVar(&flagLogLevel, "l", 0, "-l <level> Set the log level 1..5")
	flag.StringVar(&flagRegions, "r", "all", "-r <regions> Set the regions separated by comma, default will search all")
	flag.Parse()
	log.SetLevel(flagLogLevel)

	log.Info("Starting con")

	config.Load(&cfg, "config.json")
	log.Debug("Using configuration: %v", cfg)

	data.GetHosts(cfg.Regions)
}
