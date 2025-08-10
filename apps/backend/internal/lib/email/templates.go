package email

type Template string

const (
	TemplateWelcome             Template = "welcome"
	TemplateDueDateReminder     Template = "due-date-reminder"
	TemplateOverdueNotification Template = "overdue-notification"
	TemplateWeeklyReport        Template = "weekly-report"
)
