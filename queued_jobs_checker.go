package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/google/go-github/v68/github"
	"golang.org/x/oauth2"
)

func CheckQueuedJobs(ctx context.Context, checkInterval time.Duration) {
	token := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)

	githubRepoOwner := os.Getenv("GITHUB_REPO_OWNER")
	githubRepoName := os.Getenv("GITHUB_REPO_NAME")

	tc := oauth2.NewClient(ctx, token)

	githubClient := github.NewClient(tc)

	for {
		select {
		case <-ctx.Done():
			return
		default:
			// Get the list of jobs in the queue
			runs, _, err := githubClient.Actions.ListRepositoryWorkflowRuns(ctx, githubRepoOwner, githubRepoName, &github.ListWorkflowRunsOptions{
				Status: "queued",
			})
			if err != nil {
				log.Printf("Failed to list workflow runs: %v", err)
				time.Sleep(checkInterval)
				continue
			}

			for _, run := range runs.WorkflowRuns {
				log.Printf("title: %s, jobsURL: %s ", run.GetDisplayTitle(), run.GetJobsURL()) // Dereference each struct inside the slice
			}

			if len(runs.WorkflowRuns) > 0 {
				log.Printf("Found %d queued job(s). Launching new runners...", len(runs.WorkflowRuns))
				_, err := provisionSpotEC2(githubRepoName, len(runs.WorkflowRuns))
				if err != nil {
					log.Printf("Failed to launch Spot instance: %v", err)
				}
				log.Println("%d runner launched.", len(runs.WorkflowRuns))
			} else {
				log.Println("No queued jobs found.")
			}

			time.Sleep(checkInterval)
		}
	}
}
