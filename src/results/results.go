package results

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/darthguinea/golib/log"
)

func getSession(r string) *session.Session {
	return session.New(&aws.Config{
		Region: aws.String(r),
	})
}

func getAllReservations(s *ec2.EC2, nextToken *string) (*ec2.DescribeInstancesOutput, error) {
	filters := []*ec2.Filter{}
	values := []*string{}
	values = append(values, aws.String("running"))

	filters = append(filters, &ec2.Filter{
		Name:   aws.String("instance-state-name"),
		Values: values,
	})

	input := &ec2.DescribeInstancesInput{
		Filters: filters,
	}

	if nextToken == nil {
		input.NextToken = nextToken
	}

	result, err := s.DescribeInstances(input)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func GetAllReservations(region *string) *[]*ec2.DescribeInstancesOutput {
	results := []*ec2.DescribeInstancesOutput{}
	svc := ec2.New(getSession(*region))

	log.Debug("getting all reservations")

	var nextToken *string
	for {
		result, err := getAllReservations(svc, nextToken)
		if err != nil {
			log.Errors(err)
			return nil
		}

		results = append(results, result)

		if result.NextToken == nil {
			break
		}

		nextToken = result.NextToken
	}

	i := 0
	for _, res := range results {
		for _, instances := range res.Reservations {
			i = i + len(instances.Instances)
		}
	}
	log.Debug("%v instances found.", i)

	return &results
}
