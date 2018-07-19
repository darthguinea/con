package data

import (
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

func getSession(region string) *session.Session {
	return session.New(&aws.Config{
		Region: aws.String(region),
	})
}

func getHosts(wg *sync.WaitGroup) {
	time.Sleep(1 * time.Second)
	wg.Done()
}

func GetHosts(regions []string) {
	//ec3Svc := ec2.New(session)
	log.Info("Fetching results from [%v] regions", len(regions))
	// if region == "all" {
	// 	var wg sync.WaitGroup
	// 	wg.Add(3)
	// 	for i := 0; i < 3; i++ {
	// 		go getHosts(&wg)
	// 	}
	// 	wg.Wait()
	// }
}
