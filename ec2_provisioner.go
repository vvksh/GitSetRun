package main

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func provisionEC2(jobId int64, repoFullName string) (string, error) {
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
		MinCount: aws.Int64(1),
		MaxCount: aws.Int64(1),
		TagSpecifications: []*ec2.TagSpecification{
			{
				ResourceType: aws.String("instance"),
				Tags: []*ec2.Tag{
					{
						Key:   aws.String("JobId"),
						Value: aws.String(fmt.Sprintf("gitsetrun-%d", jobId)),
					},
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
