package cron

import (
	"fmt"
	"strings"
)

type JobRegistry struct {
	jobs map[string]Job
}

func NewJobRegistry() *JobRegistry {
	registry := &JobRegistry{
		jobs: make(map[string]Job),
	}

	registry.Register(&DueDateRemindersJob{})
	registry.Register(&OverdueNotificationsJob{})
	registry.Register(&WeeklyReportsJob{})
	registry.Register(&AutoArchiveJob{})

	return registry
}

func (r *JobRegistry) Register(job Job) {
	r.jobs[job.Name()] = job
}

func (r *JobRegistry) Get(name string) (Job, error) {
	job, exists := r.jobs[name]
	if !exists {
		return nil, fmt.Errorf("job '%s' not found", name)
	}
	return job, nil
}

func (r *JobRegistry) List() []string {
	names := make([]string, 0, len(r.jobs))
	for name := range r.jobs {
		names = append(names, name)
	}
	return names
}

func (r *JobRegistry) Help() string {
	var help strings.Builder
	help.WriteString("Available cron jobs:\n")

	// Calculate max job name length for alignment
	maxLen := 0
	for name := range r.jobs {
		if len(name) > maxLen {
			maxLen = len(name)
		}
	}

	for _, name := range r.List() {
		job := r.jobs[name]
		help.WriteString(fmt.Sprintf("  %-*s - %s\n", maxLen, name, job.Description()))
	}

	return help.String()
}
