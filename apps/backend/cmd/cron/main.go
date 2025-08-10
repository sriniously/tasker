package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/sriniously/tasker/internal/cron"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "cron",
		Short: "Tasker Cron Job Runner",
		Long:  "Tasker Cron Job Runner - Execute scheduled jobs for the Tasker task management system",
	}

	// List command
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List available cron jobs",
		Run: func(cmd *cobra.Command, args []string) {
			registry := cron.NewJobRegistry()
			fmt.Print(registry.Help())
		},
	}
	rootCmd.AddCommand(listCmd)

	// Create subcommands for each job
	registry := cron.NewJobRegistry()
	for _, jobName := range registry.List() {
		job, _ := registry.Get(jobName)
		// Capture jobName in closure
		name := jobName
		jobCmd := &cobra.Command{
			Use:   job.Name(),
			Short: job.Description(),
			RunE: func(cmd *cobra.Command, args []string) error {
				return runJob(name)
			},
		}
		rootCmd.AddCommand(jobCmd)
	}

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func runJob(jobName string) error {
	registry := cron.NewJobRegistry()

	job, err := registry.Get(jobName)
	if err != nil {
		return fmt.Errorf("job '%s' not found", jobName)
	}

	runner, err := cron.NewJobRunner(job)
	if err != nil {
		return fmt.Errorf("failed to create job runner: %w", err)
	}

	if err := runner.Run(); err != nil {
		return fmt.Errorf("job failed: %w", err)
	}

	return nil
}
