package data

import (
	"math/rand"
	"os"
	"sync"
	"time"

	"../log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

type Instance struct {
	InstanceID string `json:"instance_id"`
	PrivateIP  string `json:"private_ip"`
	Tags       []Tag
}

type Tag struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func getSession(r string) *session.Session {
	return session.New(&aws.Config{
		Region: aws.String(r),
	})
}

func random(min int, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min) + min
}

func getHosts(wg *sync.WaitGroup, r string) {
	randomNum := random(3, 5)
	time.Sleep(time.Duration(randomNum) * time.Second)
	wg.Done()
}

// func getHosts(wg *sync.WaitGroup, r string) (*ec2.DescribeInstancesOutput, error) {
// 	time.Sleep(1 * time.Second)
// 	ec2Svc := ec2.New(getSession(r))

// 	filter := &ec2.DescribeInstancesInput{
// 		Filters: []*ec2.Filter{
// 			{
// 				Name: aws.String("instance-state-name"),
// 				Values: []*string{
// 					aws.String("running"),
// 					aws.String("pending"),
// 				},
// 			},
// 		},
// 	}

// 	reservations, err := ec2Svc.DescribeInstances(filter)
// 	if err != nil {
// 		return nil, errors.New("Could not retrieve instance data")
// 	}
// 	wg.Done()

// 	return reservations, nil
// }

// GetHosts - External function used to fetch host information,
// parameters: r regions (regions) []strings, c bool (clear cache)
func GetHosts(r []string, c bool) {
	stat, err := os.Stat("/tmp/ec2.json")
	if err == nil {
		timestamp := stat.ModTime().Unix()
		timeNow := time.Now().Unix()
		if (timeNow - timestamp) >= 300 {
			log.Warn("Cached file is over 5 minutes old, regenerating")
			// data = getEC2Hosts()
		} else {
			log.Warn("Getting hosts from cached file")
			// data = getCachedHosts()
		}
	} else {
		log.Warn("Cached file does not exist, fetching hosts from AWS")
		// data = getEC2Hosts(r)
		getEC2Hosts(r)
	}
}

// ValidIP - Check if a string is a valid IP address or not
// parameters: ip string
// returns: bool
func ValidIP(ip string) bool {
	return false
}

func getCachedHosts() {}

func getEC2Hosts(r []string) {
	var wg sync.WaitGroup
	wg.Add(len(r))

	log.Info("Fetching results from [%v] regions", len(r))
	startTime := time.Now()
	for _, r := range r {
		go getHosts(&wg, r)
	}
	wg.Wait()
	elapsed := time.Since(startTime)
	log.Info("Data retrived in [%.4v] seconds", elapsed)
}
