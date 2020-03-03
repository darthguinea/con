package cache

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/darthguinea/golib/log"
)

func RemoveCacheFile(path *string) {
	f, err := os.Stat(*path)
	if err != nil {
		log.Warn("cache file [%v] does not exist", *path)
		return
	}

	// if is not directory
	if !f.IsDir() {
		os.Remove(*path)
	}
}

func IsCacheExpired(path *string, ttl *int64) bool {
	stat, err := os.Stat(*path)
	if err == nil {
		timestamp := stat.ModTime().Unix()
		timeNow := time.Now().Unix()

		if (timeNow - timestamp) >= *ttl {
			log.Warn("cache file [%v] is over [%v] (expired), re-creating", *path, *ttl)
			return true
		} else {
			log.Info("reading from cache file")
			return false
		}
	}

	log.Warn("cache file [%v] does not exist, creating")
	return true
}

func Save(path *string, v *[]*[]*ec2.DescribeInstancesOutput) {
	j, _ := json.Marshal(v)
	ioutil.WriteFile(*path, j, 0644)
}

func Load(path *string) *[]*[]*ec2.DescribeInstancesOutput {
	var reservations []*[]*ec2.DescribeInstancesOutput
	file, err := ioutil.ReadFile(*path)

	if err != nil {
		log.Fatal(err)
		return nil
	}

	json.Unmarshal(file, &reservations)

	return &reservations
}
