package main

import (
	"flag"

	"./src/config"
	"./src/data"
	"./src/log"
)

var Config Configuration

type Configuration struct {
	Regions []string `json:regions`
}

func main() {
	var (
		flagLogLevel   int
		flagRegions    string
		flagClearCache bool
	)

	flag.IntVar(&flagLogLevel, "l", 0, "-l <level> Set the log level 1..5")
	flag.StringVar(&flagRegions, "r", "all", "-r <regions> Set the regions separated by comma, default will search all")
	flag.BoolVar(&flagClearCache, "c", false, "-c Clear cache")
	flag.Parse()
	log.SetLevel(flagLogLevel)

	log.Info("Starting con")

	if flagClearCache {
		log.Info("Clearing cache")
	}

	config.Load(&Config, "config.json")
	log.Debug("Using configuration: %v", Config)

	data.GetHosts(Config.Regions, flagClearCache)
}
