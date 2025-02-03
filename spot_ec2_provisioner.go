package main

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func provisionSpotEC2(repoFullName string, numInstances int) (string, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-west-2"),
	})

	if err != nil {
		return "", fmt.Errorf("failed to create session: %w", err)
	}

	svc := ec2.New(sess)

	runResult, err := svc.RunInstances(&ec2.RunInstancesInput{
		LaunchTemplate: &ec2.LaunchTemplateSpecification{
			LaunchTemplateName: aws.String("gitsetrun-launch-template"),
			Version:            aws.String("$Latest"),
		},
		InstanceMarketOptions: &ec2.InstanceMarketOptionsRequest{
			MarketType: aws.String("spot"),
		},

		MinCount: aws.Int64(int64(numInstances)),
		MaxCount: aws.Int64(int64(numInstances)),
		TagSpecifications: []*ec2.TagSpecification{
			{
				ResourceType: aws.String("instance"),
				Tags: []*ec2.Tag{
					{
						Key:   aws.String("Repository"),
						Value: aws.String(repoFullName),
					},
					{
						Key:   aws.String("Purpose"),
						Value: aws.String("githubrunner"),
					},
				},
			},
		},
	})
	if err != nil {
		return "", fmt.Errorf("failed to run instance: %w", err)
	}

	instanceId := *runResult.Instances[0].InstanceId
	log.Printf("Instance %s created", instanceId)
	return instanceId, nil
}
