package job

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"github.com/sriniously/tasker/internal/model/todo"
)

const (
	TaskWelcome           = "email:welcome"
	TaskReminderEmail     = "email:reminder"
	TaskWeeklyReportEmail = "email:weekly_report"
)

type WelcomeEmailPayload struct {
	To        string `json:"to"`
	FirstName string `json:"first_name"`
}

func NewWelcomeEmailTask(to, firstName string) (*asynq.Task, error) {
	payload, err := json.Marshal(WelcomeEmailPayload{
		To:        to,
		FirstName: firstName,
	})
	if err != nil {
		return nil, err
	}

	return asynq.NewTask(TaskWelcome, payload,
		asynq.MaxRetry(3),
		asynq.Queue("default"),
		asynq.Timeout(30*time.Second)), nil
}

type ReminderEmailTask struct {
	UserID    string    `json:"user_id"`
	TodoID    uuid.UUID `json:"todo_id"`
	TodoTitle string    `json:"todo_title"`
	DueDate   time.Time `json:"due_date"`
	TaskType  string    `json:"task_type"` // "due_date_reminder" or "overdue_notification"
}

func EnqueueReminderEmail(client *asynq.Client, task *ReminderEmailTask) error {
	payload, err := json.Marshal(task)
	if err != nil {
		return err
	}

	asynqTask := asynq.NewTask(TaskReminderEmail, payload,
		asynq.MaxRetry(3),
		asynq.Queue("default"),
		asynq.Timeout(30*time.Second))

	_, err = client.Enqueue(asynqTask)
	return err
}

type WeeklyReportEmailTask struct {
	UserID         string               `json:"user_id"`
	WeekStart      time.Time            `json:"week_start"`
	WeekEnd        time.Time            `json:"week_end"`
	CompletedCount int                  `json:"completed_count"`
	ActiveCount    int                  `json:"active_count"`
	OverdueCount   int                  `json:"overdue_count"`
	CompletedTodos []todo.PopulatedTodo `json:"completed_todos"`
	OverdueTodos   []todo.PopulatedTodo `json:"overdue_todos"`
}

func EnqueueWeeklyReportEmail(client *asynq.Client, task *WeeklyReportEmailTask) error {
	payload, err := json.Marshal(task)
	if err != nil {
		return err
	}

	asynqTask := asynq.NewTask(TaskWeeklyReportEmail, payload,
		asynq.MaxRetry(3),
		asynq.Queue("default"),
		asynq.Timeout(60*time.Second)) // Longer timeout for report generation

	_, err = client.Enqueue(asynqTask)
	return err
}
