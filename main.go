package main

import (
	"flag"
	"strings"

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
		flagOr         bool
	)

	flag.IntVar(&flagLogLevel, "l", 0, "-l <level> Set the log level 1..5")
	flag.StringVar(&flagRegions, "r", "all", "-r <regions> Set the regions separated by comma, default will use the config")
	flag.BoolVar(&flagClearCache, "c", false, "-c Clear cache")
	flag.BoolVar(&flagOr, "o", false, "-o Search using OR instead of AND")
	flag.Parse()

	config.Load(&Config, "config.json")
	log.SetLevel(flagLogLevel)

	search := flag.Args()

	log.Info("Starting con")

	log.Debug("Searching values %v", search)
	if flagOr {
		log.Debug("Using OR flag to search")
	}

	if flagClearCache {
		log.Info("Clearing cache")
	}

	if strings.Compare(flagRegions, "all") != 0 {
		Config.Regions = strings.Split(flagRegions, ",")
		log.Debug("Overriding region settings %v", Config.Regions)
	}

	log.Debug("Using configuration: %v", Config)

	data.GetHosts(Config.Regions, flagClearCache, search)
}
