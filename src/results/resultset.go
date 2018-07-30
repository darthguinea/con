package results

import (
	"sort"

	"github.com/aws/aws-sdk-go/service/ec2"
)

// ResultSet contains AWS instance and score for search results
type ResultSet struct {
	Score    int           `json:"score"`
	Instance *ec2.Instance `json:"Instance"`
}

type By func(rs1, rs2 *ResultSet) bool

// Len is part of sort.Interface.
func (s *resultsSorter) Len() int {
	return len(s.results)
}

// Swap is part of sort.Interface.
func (s *resultsSorter) Swap(i, j int) {
	s.results[i], s.results[j] = s.results[j], s.results[i]
}

// Less is part of sort.Interface. It is implemented by calling the "by" closure in the sorter.
func (s *resultsSorter) Less(i, j int) bool {
	return s.by(&s.results[i], &s.results[j])
}

func (by By) Sort(results []ResultSet) {
	rs := &resultsSorter{
		results: results,
		by:      by,
	}
	sort.Sort(rs)
}

type resultsSorter struct {
	results []ResultSet
	by      func(rs1, rs2 *ResultSet) bool // Closure used in the Less method.
}
