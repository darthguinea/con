package main

import (
	"flag"

	"./src/con"
	"./src/config"
	"./src/log"
)

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

	cfg := config.Load("config.json")

	con.GetHosts("regions: %v")
}
