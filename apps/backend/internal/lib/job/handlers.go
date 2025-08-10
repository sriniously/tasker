package job

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog"
	"github.com/sriniously/tasker/internal/config"
	"github.com/sriniously/tasker/internal/lib/email"
)

func (j *JobService) InitHandlers(config *config.Config, logger *zerolog.Logger) {
	j.emailClient = email.NewClient(config, logger)
}

func (j *JobService) handleWelcomeEmailTask(ctx context.Context, t *asynq.Task) error {
	var p WelcomeEmailPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("failed to unmarshal welcome email payload: %w", err)
	}

	j.logger.Info().
		Str("type", "welcome").
		Str("to", p.To).
		Msg("Processing welcome email task")

	err := j.emailClient.SendWelcomeEmail(
		p.To,
		p.FirstName,
	)
	if err != nil {
		j.logger.Error().
			Str("type", "welcome").
			Str("to", p.To).
			Err(err).
			Msg("Failed to send welcome email")
		return err
	}

	j.logger.Info().
		Str("type", "welcome").
		Str("to", p.To).
		Msg("Successfully sent welcome email")
	return nil
}

func (j *JobService) handleReminderEmailTask(ctx context.Context, t *asynq.Task) error {
	var p ReminderEmailTask
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("failed to unmarshal reminder email payload: %w", err)
	}

	j.logger.Info().
		Str("type", p.TaskType).
		Str("user_id", p.UserID).
		Str("todo_id", p.TodoID.String()).
		Str("todo_title", p.TodoTitle).
		Msg("Processing reminder email task")

	userEmail, err := j.authService.GetUserEmail(ctx, p.UserID)
	if err != nil {
		j.logger.Error().
			Str("type", p.TaskType).
			Str("user_id", p.UserID).
			Err(err).
			Msg("Failed to resolve user email")
		return fmt.Errorf("failed to resolve user email for user %s: %w", p.UserID, err)
	}

	switch p.TaskType {
	case "due_date_reminder":
		err = j.emailClient.SendDueDateReminderEmail(
			userEmail,
			p.TodoTitle,
			p.TodoID,
			p.DueDate,
		)
	case "overdue_notification":
		err = j.emailClient.SendOverdueNotificationEmail(
			userEmail,
			p.TodoTitle,
			p.TodoID,
			p.DueDate,
		)
	default:
		return fmt.Errorf("unknown reminder task type: %s", p.TaskType)
	}

	if err != nil {
		j.logger.Error().
			Str("type", p.TaskType).
			Str("user_id", p.UserID).
			Str("todo_id", p.TodoID.String()).
			Err(err).
			Msg("Failed to send reminder email")
		return err
	}

	j.logger.Info().
		Str("type", p.TaskType).
		Str("user_id", p.UserID).
		Str("todo_id", p.TodoID.String()).
		Msg("Successfully sent reminder email")
	return nil
}

func (j *JobService) handleWeeklyReportEmailTask(ctx context.Context, t *asynq.Task) error {
	var p WeeklyReportEmailTask
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("failed to unmarshal weekly report email payload: %w", err)
	}

	j.logger.Info().
		Str("type", "weekly_report").
		Str("user_id", p.UserID).
		Int("completed_count", p.CompletedCount).
		Int("active_count", p.ActiveCount).
		Int("overdue_count", p.OverdueCount).
		Msg("Processing weekly report email task")

	userEmail, err := j.authService.GetUserEmail(ctx, p.UserID)
	if err != nil {
		j.logger.Error().
			Str("type", "weekly_report").
			Str("user_id", p.UserID).
			Err(err).
			Msg("Failed to resolve user email")
		return fmt.Errorf("failed to resolve user email for user %s: %w", p.UserID, err)
	}

	err = j.emailClient.SendWeeklyReportEmail(
		userEmail,
		p.WeekStart,
		p.WeekEnd,
		p.CompletedCount,
		p.ActiveCount,
		p.OverdueCount,
		p.CompletedTodos,
		p.OverdueTodos,
	)
	if err != nil {
		j.logger.Error().
			Str("type", "weekly_report").
			Str("user_id", p.UserID).
			Err(err).
			Msg("Failed to send weekly report email")
		return err
	}

	j.logger.Info().
		Str("type", "weekly_report").
		Str("user_id", p.UserID).
		Msg("Successfully sent weekly report email")
	return nil
}
