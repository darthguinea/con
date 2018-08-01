package data

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"sync"
	"time"

	"../log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func getSession(r string) *session.Session {
	return session.New(&aws.Config{
		Region: aws.String(r),
	})
}

func filterHosts(search []string) {

}

// GetHosts - External function used to fetch host information,
// parameters: regions []string, clear_cache bool, search []string
func GetHosts(r []string, c bool) []*ec2.DescribeInstancesOutput {
	stat, err := os.Stat("/tmp/ec2.json")

	if err == nil {
		timestamp := stat.ModTime().Unix()
		timeNow := time.Now().Unix()

		if (timeNow - timestamp) >= 300 {
			log.Warn("Cached file is over 5 minutes old, regenerating")
			data := getEC2Hosts(r)
			writeCachedFile(data)
			return data
		} else {
			log.Warn("Getting hosts from cached file")
			data := getCachedHosts()
			return data
		}
	}

	log.Warn("Cached file does not exist, fetching hosts from AWS")
	data := getEC2Hosts(r)
	writeCachedFile(data)

	return data
}

func getCachedHosts() []*ec2.DescribeInstancesOutput {
	var r []*ec2.DescribeInstancesOutput
	file, err := ioutil.ReadFile("/tmp/ec2.json")

	if err != nil {
		log.Fatal(err)
		return nil
	}

	json.Unmarshal(file, &r)

	return r
}

func writeCachedFile(d []*ec2.DescribeInstancesOutput) {
	j, _ := json.Marshal(d)
	ioutil.WriteFile("/tmp/ec2.json", j, 0644)
}

// ValidIP - Check if a string is a valid IP address or not
// parameters: ip string
// returns: bool
func ValidIP(ip string) bool {
	return false
}

func getEC2Hosts(r []string) []*ec2.DescribeInstancesOutput {
	var wg sync.WaitGroup
	wg.Add(len(r))

	log.Info("Fetching results from [%v] regions", len(r))
	startTime := time.Now()
	var instancesOutput []*ec2.DescribeInstancesOutput
	for _, region := range r {
		go func(region string) {
			i, _ := getEC2HostsThreads(&wg, region)
			instancesOutput = append(instancesOutput, i)
		}(region)
	}
	wg.Wait()
	elapsed := time.Since(startTime)
	log.Info("Data retrived in [%v]", elapsed)
	return instancesOutput
}

func getEC2HostsThreads(wg *sync.WaitGroup, r string) (*ec2.DescribeInstancesOutput, error) {
	ec2Svc := ec2.New(getSession(r))

	log.Debug("Searching region [%v]", r)

	filter := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			{
				Name: aws.String("instance-state-name"),
				Values: []*string{
					aws.String("running"),
					aws.String("pending"),
				},
			},
		},
	}

	reservations, err := ec2Svc.DescribeInstances(filter)
	if err != nil {
		return nil, errors.New("Could not retrieve instance data")
	}
	wg.Done()

	return reservations, nil
}
