package results

import (
	"os"
	"strings"

	"../cmd"
	"../log"

	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/olekukonko/tablewriter"
)

// RemoveCacheFile Removes the cache file
func RemoveCacheFile(path string) {
	if f, _ := os.Stat(path); f.IsDir() == false {
		log.Warn("Removing cache file")
		os.Remove(path)
	}
}

// DrawTable takes the resultset and converts it into a table
// parameters: []ResultSet
func DrawTable(display map[string]string, rs []ResultSet) {
	customHeader := []string{}
	rowOrdering := []string{}

	for k, v := range display {
		customHeader = append(customHeader, k)
		rowOrdering = append(rowOrdering, v)
	}

	defaultHeader := []string{
		"Instance Id",
		"Private Ip",
		"Launch Time",
		"State",
		"Region",
	}

	header := append(customHeader, defaultHeader...)

	log.Info("Drawing table of results")
	log.Debug("Custom header: %v", customHeader)
	log.Debug("Default header: %v", defaultHeader)
	log.Debug("Header: %v", header)
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(header)
	table.SetRowLine(true)

	pip := ""
	count := 0
	for _, ins := range rs {
		customRow := []string{}
		defaultRow := []string{}

		log.Debug("New row:")

		for _, v := range rowOrdering {
			value := getTagValue(ins.Instance.Tags, v)
			log.Debug("Searching for element item [%v] found [%v]", v, value)
			customRow = append(customRow, value)
		}
		log.Debug("Appending custom row item [%v]", customRow)

		// dt := ins.LaunchTime.Local().Format("Mon Jan 2 15:04:05 MST 2006"),
		defaultRow = []string{
			*ins.Instance.InstanceId,
			*ins.Instance.PrivateIpAddress,
			ins.Instance.LaunchTime.Local().Format("Mon Jan 2 15:04:05 MST 2006"),
			*ins.Instance.State.Name,
			*ins.Instance.Placement.AvailabilityZone,
		}

		row := append(customRow, defaultRow...)

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
func Filter(r []*ec2.DescribeInstancesOutput, s []string, tags []string) []ResultSet {
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
					// Search the placement as well:
					if strings.Contains(*ins.Placement.AvailabilityZone, search) {
						score = score + 200
					}
					for _, t := range ins.Tags {
						if searchTag(tags, *t.Key) {
							if strings.Contains(strings.ToLower(*t.Value), strings.ToLower(search)) {
								foundKeyword = true
								score = score + 100
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

func getTagValue(tags []*ec2.Tag, key string) string {
	for _, v := range tags {
		if strings.ToLower(key) == strings.ToLower(*v.Key) {
			return *v.Value
		}
	}
	return ""
}

func searchTag(tags []string, val string) bool {
	for _, tag := range tags {
		if strings.Contains(strings.ToLower(tag), strings.ToLower(val)) {
			return true
		}
	}
	return false
}
