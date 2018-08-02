package main

import (
	"flag"
	"os"
	"path"
	"strings"

	"./src/config"
	"./src/data"
	"./src/log"
	"./src/results"
)

// Config global struct for configuration
var Config Configuration

// Configuration data structure of config file
type Configuration struct {
	CacheFile string              `json:"cache_file"`
	Regions   []string            `json:"regions"`
	Search    SearchConfiguration `json:"search"`
	Display   map[string]string   `json:"display"`
}

// SearchConfiguration struct of tags to search (loaded from config file)
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
		flagTable      bool
	)

	flag.IntVar(&flagLogLevel, "l", 0, "-l <level> Set the log level 1..5")
	flag.StringVar(&flagRegions, "r", "config", "-r <regions> Set the regions separated by comma, default will use the config")
	flag.StringVar(&flagTags, "t", "config", "-t <tags> Which tags to search seperated by a comma, default uses config")
	flag.BoolVar(&flagClearCache, "C", false, "-C Clear cache")
	flag.IntVar(&flagCount, "c", 5, "-c <count> Number of results to show")
	flag.BoolVar(&flagTable, "Q", false, "-Q Do not render table (for debugging)")
	flag.Parse()

	ex, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	dir := path.Dir(ex)

	config.Load(&Config, dir+"/config.json")
	log.SetLevel(flagLogLevel)

	search := flag.Args()

	log.Info("Starting con")

	log.Debug("Searching values %v", search)

	if flagClearCache {
		log.Info("Clearing cache")
		results.RemoveCacheFile(Config.CacheFile)
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
	rs := results.Filter(d, search, Config.Search.Tags)

	score := func(rs1, rs2 *results.ResultSet) bool {
		return rs1.Score > rs2.Score
	}
	results.By(score).Sort(rs)

	if !flagTable {
		if flagCount > len(rs) {
			flagCount = len(rs)
		}
		log.Info("Custom display items %v", Config.Display)
		results.DrawTable(Config.Display, rs[0:flagCount])
	}
}
