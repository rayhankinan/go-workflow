package main

import (
	"context"
	"fmt"
	"log"
	"time"

	executions "cloud.google.com/go/workflows/executions/apiv1beta"
	executionspb "cloud.google.com/go/workflows/executions/apiv1beta/executionspb"
	"google.golang.org/api/option"
)

func main() {
	// Set your project, location, and workflow details
	projectID := "sequencing-lab"
	location := "asia-southeast2"
	workflow := "workflow-qa-integration"

	ctx := context.Background()

	// Create a new Workflows client
	client, err := executions.NewClient(ctx, option.WithCredentialsFile("sequencing-lab-19c8c46fbaf8.json"))
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	defer client.Close()

	// Define the request to execute the workflow
	req := &executionspb.CreateExecutionRequest{
		Parent: fmt.Sprintf("projects/%s/locations/%s/workflows/%s", projectID, location, workflow),
		Execution: &executionspb.Execution{
			Argument: `{"key":"value"}`, // Optional JSON argument
		},
	}

	// Execute the workflow
	resp, err := client.CreateExecution(ctx, req)
	if err != nil {
		log.Fatalf("failed to execute workflow: %v", err)
	}

	fmt.Printf("Execution started: %s\n", resp.Name)

	// Wait for the workflow execution to complete
	for {
		exec, err := client.GetExecution(ctx, &executionspb.GetExecutionRequest{Name: resp.Name})
		if err != nil {
			log.Fatalf("failed to get execution: %v", err)
		}

		if exec.State == executionspb.Execution_SUCCEEDED || exec.State == executionspb.Execution_FAILED {
			fmt.Printf("Execution finished with state: %s\n", exec.State.String())
			fmt.Printf("Execution result: %s\n", exec.Result)
			break
		}

		fmt.Printf("Execution state: %s\n", exec.State.String())

		// Sleep for a bit before checking the status again
		time.Sleep(5 * time.Second)
	}
}
