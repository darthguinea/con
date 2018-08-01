package results

import (
	"os"
	"strings"

	"../cmd"
	"../log"

	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/olekukonko/tablewriter"
)

func DrawTable(rs []ResultSet) {
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

	pip := ""
	count := 0
	for _, ins := range rs {
		var name, environment string
		for _, t := range ins.Instance.Tags {
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
			*ins.Instance.InstanceId,
			*ins.Instance.PrivateIpAddress,
			ins.Instance.LaunchTime.Local().Format("Mon Jan 2 15:04:05 MST 2006"),
			*ins.Instance.State.Name,
		}
		pip = *ins.Instance.PrivateIpAddress
		count++
		table.Append(row)
	}
	table.Render()

	if count == 1 {
		cmd.Exec("ssh " + pip)
	}
}

// Filter - External function used to filter ec2 hosts
// parameters: Reservations []*ec2.DescribeInstancesOutput, search []string, tags []string, or_flag bool
// returns: []*ec2.Instance
func Filter(r []*ec2.DescribeInstancesOutput, s []string, tags []string, o bool) []ResultSet {
	var rs []ResultSet

	for _, rv := range r {
		for _, res := range rv.Reservations {
			for _, ins := range res.Instances {
				score := 0
				keywordScore := 0

				for _, search := range s {
					foundKeyword := false
					if *ins.InstanceId == search {
						score = 1000
					}
					if *ins.PrivateIpAddress == search {
						score = 1000
					}
					for _, t := range ins.Tags {
						if searchTag(tags, *t.Key) {
							if strings.Contains(strings.ToLower(*t.Value), strings.ToLower(search)) {
								foundKeyword = true
								score = score + 100
								log.Info("%v = %v %v", *t.Key, *t.Value, score)
							}
						}
					}
					if foundKeyword {
						keywordScore = keywordScore + 100
						foundKeyword = false
					}
				}

				if score > 0 {
					tmpInstance := ResultSet{
						score + keywordScore,
						ins,
					}
					rs = append(rs, tmpInstance)
				}
			}
		}
	}
	log.Debug("Found [%v]", len(rs))
	return rs
}

func searchTag(tags []string, val string) bool {
	for _, tag := range tags {
		if strings.Contains(strings.ToLower(tag), strings.ToLower(val)) {
			return true
		}
	}
	return false
}
