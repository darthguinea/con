package data

import (
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/darthguinea/golib/log"
)

type Instance struct {
	Search           string
	Name             string
	InstanceId       string
	Environment      string
	EnvironmentTag   string
	SubnetId         string
	PublicIp         string
	PrivateIp        string
	State            string
	Type             string
	AvailabilityZone string
	LaunchTime       time.Time
}

func ConvertToArray(r *[]*[]*ec2.DescribeInstancesOutput) []*Instance {
	results := []*Instance{}

	log.Debug("converting data to array structure")

	for _, block := range *r {
		for _, reservation := range *block {
			for _, instances := range reservation.Reservations {
				for _, instance := range instances.Instances {

					instanceName := getTagValue(instance.Tags, "Name")

					i := &Instance{
						Search: instanceName +
							" " + *instance.InstanceId +
							" " + *instance.PrivateIpAddress,
						Name:             instanceName,
						InstanceId:       *instance.InstanceId,
						Environment:      getTagValue(instance.Tags, "Environment"),
						EnvironmentTag:   getTagValue(instance.Tags, "EnvironmentTag"),
						SubnetId:         *instance.SubnetId,
						PrivateIp:        *instance.PrivateIpAddress,
						State:            *instance.State.Name,
						Type:             *instance.InstanceType,
						AvailabilityZone: *instance.Placement.AvailabilityZone,
						LaunchTime:       *instance.LaunchTime,
					}

					if instance.PublicIpAddress != nil {
						i.PublicIp = *instance.PublicIpAddress
					}

					results = append(results, i)
				}
			}
		}
	}

	return results
}

func getTagValue(tags []*ec2.Tag, key string) string {
	for _, v := range tags {
		if strings.ToLower(key) == strings.ToLower(*v.Key) {
			return *v.Value
		}
	}
	return ""
}
