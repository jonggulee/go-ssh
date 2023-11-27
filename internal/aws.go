package internal

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2_types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

type (
	Target struct {
		Name          string
		PublicDomain  string
		PrivateDomain string
	}
)

func NewConfig(ctx context.Context) (cfg aws.Config, err error) {
	cfg, err = config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}
	return
}

func FindeInstances(ctx context.Context, cfg aws.Config) (map[string]*Target, error) {
	var (
		client     = ec2.NewFromConfig(cfg)
		table      = make(map[string]*Target)
		outputFunc = func(table map[string]*Target, output *ec2.DescribeInstancesOutput) {
			for _, rv := range output.Reservations {
				for _, inst := range rv.Instances {
					name := ""
					for _, tag := range inst.Tags {
						if aws.ToString(tag.Key) == "Name" {
							name = aws.ToString(tag.Value)
							break
						}
					}
					table[fmt.Sprintf("%s\t(%s)", name, *inst.InstanceId)] = &Target{
						Name:          aws.ToString(inst.InstanceId),
						PublicDomain:  aws.ToString(inst.PublicDnsName),
						PrivateDomain: aws.ToString(inst.PrivateDnsName),
					}
				}
			}
		}
	)

	output, err := client.DescribeInstances(context.Background(),
		&ec2.DescribeInstancesInput{
			Filters: []ec2_types.Filter{
				{Name: aws.String("instance-state-name"), Values: []string{"running"}},
			},
		})
	if err != nil {
		fmt.Println(err)
	}

	outputFunc(table, output)
	return table, err
}
