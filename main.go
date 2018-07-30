package main

import (
	"flag"
	"strings"

	"./src/config"
	"./src/data"
	"./src/log"
	"./src/results"
)

var Config Configuration

type Configuration struct {
	Regions []string            `json:"regions"`
	Search  SearchConfiguration `json:"search"`
}
type SearchConfiguration struct {
	Tags []string `json:"tags"`
}

func main() {
	var (
		flagLogLevel   int
		flagRegions    string
		flagTags       string
		flagClearCache bool
		flagCount      int
		flagOr         bool
		flagTable      bool
	)

	flag.IntVar(&flagLogLevel, "l", 0, "-l <level> Set the log level 1..5")
	flag.StringVar(&flagRegions, "r", "config", "-r <regions> Set the regions separated by comma, default will use the config")
	flag.StringVar(&flagTags, "t", "config", "-t <tags> Which tags to search seperated by a comma, default uses config")
	flag.BoolVar(&flagClearCache, "C", false, "-C Clear cache")
	flag.IntVar(&flagCount, "c", 5, "-c <count> Number of results to show")
	flag.BoolVar(&flagOr, "o", false, "-o Search using OR instead of AND")
	flag.BoolVar(&flagTable, "D", false, "-D Do not render table (for debugging)")
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

	if strings.Compare(flagRegions, "config") != 0 {
		Config.Regions = strings.Split(flagRegions, ",")
		log.Debug("Overriding region settings %v", Config.Regions)
	}
	if strings.Compare(flagTags, "config") != 0 {
		Config.Search.Tags = strings.Split(flagTags, ",")
		log.Debug("Overriding tags search %v", Config.Search.Tags)
	}
	log.Debug("Searching %v", Config.Search)
	log.Debug("Searching tags %v", Config.Search.Tags)

	log.Debug("Using configuration: %v", Config)

	d := data.GetHosts(Config.Regions, flagClearCache)
	rs := results.Filter(d, search, Config.Search.Tags, flagOr)

	score := func(rs1, rs2 *results.ResultSet) bool {
		return rs1.Score > rs2.Score
	}
	results.By(score).Sort(rs)

	if !flagTable {
		results.DrawTable(rs[0:flagCount])
	}
}
