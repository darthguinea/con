package results

import (
	"os"
	"strings"

	"../log"

	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/olekukonko/tablewriter"
)

func DrawTable(r []*ec2.DescribeInstancesOutput) {
	header := []string{
		"Name",
		"Environment",
		"Instance Id",
		"Private Ip",
		"Launch Time",
		"State",
	}

	log.Info("Drawing table of results")
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(header)
	table.SetRowLine(true)

	count := 0
	for _, rvs := range r {
		for _, rv := range rvs.Reservations {
			for _, ins := range rv.Instances {
				count++
				var name, environment string
				for _, t := range ins.Tags {
					if *t.Key == "Name" {
						name = *t.Value
					}
					if *t.Key == "Environment" {
						environment = *t.Value
					}
				}

				// dt := ins.LaunchTime.Local().Format("Mon Jan 2 15:04:05 MST 2006"),
				row := []string{
					name,
					environment,
					*ins.InstanceId,
					*ins.PrivateIpAddress,
					ins.LaunchTime.Local().Format("Mon Jan 2 15:04:05 MST 2006"),
					*ins.State.Name,
				}
				table.Append(row)
				if count > 9 {
					table.Render()
					return
				}
			}
		}
	}
}

func Filter(r []*ec2.DescribeInstancesOutput, s []string, o bool) {
	for _, rv := range r {
		for _, res := range rv.Reservations {
			for _, ins := range res.Instances {
				shouldAppend := true

				for _, search := range s {
					if !strings.Contains(*ins.InstanceId, search) {
						shouldAppend = false
					}
				}

				if shouldAppend {
					log.Info("ins: %v", *ins.InstanceId)
				}
			}
		}
	}
}
