package main

import (
	"flag"
	"os"
	"strings"
	"sync"

	"./src/cache"
	// "./src/cmd"
	"./src/config"
	"./src/data"
	"./src/results"

	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/darthguinea/golib/iv"
	"github.com/darthguinea/golib/log"
	"github.com/ktr0731/go-fuzzyfinder"
)

func main() {
	var (
		flagVerbose         bool
		flagDebug           bool
		flagRemoveCacheFile bool
		flagConfigPath      string
		flagDontSSH         bool
	)

	flagConfigDefault := strings.Join([]string{iv.CWD(), "con.json"}, "/")
	flag.StringVar(&flagConfigPath, "c", flagConfigDefault, "-c <config_file> config file location")
	flag.BoolVar(&flagVerbose, "v", false, "-v verbose.")
	flag.BoolVar(&flagDebug, "D", false, "-D debug mode.")
	flag.BoolVar(&flagDontSSH, "s", false, "-s dont ssh to server")
	flag.BoolVar(&flagRemoveCacheFile, "C", false, "-C remove cache file (force scan).")
	flag.Parse()

	if flagVerbose {
		log.SetLevel(log.INFO)
	}
	if flagDebug {
		log.SetLevel(log.DEBUG)
	}

	cfg := config.Config{}
	iv.Load(&cfg, &flagConfigPath)

	if flagRemoveCacheFile || cache.IsCacheExpired(cfg.CacheFile, cfg.CacheTTLSeconds) {
		log.Debug("reading cache file [%v]", *cfg.CacheFile)
		cache.RemoveCacheFile(cfg.CacheFile)

		wg := sync.WaitGroup{}

		var reservations []*[]*ec2.DescribeInstancesOutput
		for _, r := range *cfg.Regions {
			wg.Add(1)
			go func(wg *sync.WaitGroup, region *string) {
				log.Info("scanning region %v", *region)

				result := results.GetAllReservations(region)
				reservations = append(reservations, result)

				wg.Done()
			}(&wg, r)
		}
		wg.Wait()

		cache.Save(cfg.CacheFile, &reservations)
	}

	instancesOutput := cache.Load(cfg.CacheFile)
	instances := data.ConvertToArray(instancesOutput)
	idx, err := fuzzyfinder.FindMulti(
		instances,
		func(i int) string {
			return instances[i].Search
		},
		fuzzyfinder.WithPreviewWindow(func(i, w, h int) string {
			if i == -1 {
				return ""
			}
			return log.Sprintf("Name: \t%s\n"+
				"InstanceId: %s\n"+
				"Environment: %s\n"+
				"Stack: \t%s\n"+
				"Subnet: \t%s\n"+
				"Availability Zone: \t%s\n"+
				"Public Ip: \t%s\n"+
				"Private Ip: \t%s\n"+
				"Type: \t%s\n"+
				"State: \t%s\n"+
				"LaunchTime: \t%s\n",
				instances[i].Name,
				instances[i].InstanceId,
				instances[i].Environment,
				instances[i].EnvironmentTag,
				instances[i].SubnetId,
				instances[i].AvailabilityZone,
				instances[i].PublicIp,
				instances[i].PrivateIp,
				instances[i].Type,
				instances[i].State,
				instances[i].LaunchTime,
			)
		}))
	if err != nil {
		log.Print("exiting")
		os.Exit(1)
	}
	if flagDontSSH {
		log.Info("not connecting to server")
		os.Exit(0)
	}
	log.Info("connecting to server")
	if len(instances[idx[0]].PublicIp) > 1 {
		log.Info("using public ip address [%v]", instances[idx[0]].PublicIp)
		// cmd.Exec("ssh " + instances[idx[0]].PublicIp)
		log.Print("%v", instances[idx[0]].PublicIp)
	} else {
		log.Info("using private ip address [%v]", instances[idx[0]].PrivateIp)
		// cmd.Exec("ssh " + instances[idx[0]].PrivateIp)
		log.Print("%v", instances[idx[0]].PrivateIp)
	}
}
