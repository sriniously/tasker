import {
  Body,
  Button,
  Container,
  Head,
  Heading,
  Hr,
  Html,
  Img,
  Link,
  Preview,
  Section,
  Text,
  Tailwind,
} from "@react-email/components";

interface Todo {
  id: string;
  title: string;
  priority: string;
  dueDate?: string;
}

interface WeeklyReportEmailProps {
  weekStart: string;
  weekEnd: string;
  completedCount: string;
  activeCount: string;
  overdueCount: string;
  completedTodos: Todo[];
  overdueTodos: Todo[];
  hasCompleted: boolean;
  hasOverdue: boolean;
}

export const WeeklyReportEmail = ({
  weekStart = "{{.WeekStart}}",
  weekEnd = "{{.WeekEnd}}",
  completedCount = "{{.CompletedCount}}",
  activeCount = "{{.ActiveCount}}",
  overdueCount = "{{.OverdueCount}}",
  completedTodos = [],
  overdueTodos = [],
  hasCompleted = false,
  hasOverdue = false,
}: WeeklyReportEmailProps) => {
  const totalTodos =
    parseInt(completedCount) + parseInt(activeCount) + parseInt(overdueCount);
  const completionRate =
    totalTodos > 0
      ? Math.round((parseInt(completedCount) / totalTodos) * 100)
      : 0;

  const getMotivationalMessage = () => {
    if (completionRate >= 80) return "🌟 Outstanding work this week!";
    if (completionRate >= 60) return "👍 Great progress this week!";
    if (completionRate >= 40) return "💪 Keep pushing forward!";
    return "🎯 Let's focus on the priorities ahead!";
  };

  const getPriorityBadgeColor = (priority: string) => {
    switch (priority?.toLowerCase()) {
      case "high":
        return "bg-red-100 text-red-800";
      case "medium":
        return "bg-yellow-100 text-yellow-800";
      case "low":
        return "bg-green-100 text-green-800";
      default:
        return "bg-gray-100 text-gray-800";
    }
  };

  return (
    <Html>
      <Head />
      <Preview>
        Your Weekly Productivity Report ({weekStart} - {weekEnd})
      </Preview>
      <Tailwind>
        <Body className="bg-gray-100 font-sans">
          <Container className="bg-white p-8 rounded-lg shadow-sm my-10 mx-auto max-w-[600px]">
            <Section className="mb-6 text-center">
              <Img
                src="http://localhost:8080/static/full_logo.png?height=48&width=48"
                width="48"
                height="48"
                alt="Tasker Logo"
                className="mx-auto"
              />
              <Heading className="text-2xl font-bold text-gray-800 mt-4">
                📊 Weekly Report
              </Heading>
              <Text className="text-gray-600 text-lg">
                {weekStart} - {weekEnd}
              </Text>
            </Section>

            <Section className="mb-6">
              <Text className="text-xl font-semibold text-gray-800 mb-4">
                {getMotivationalMessage()}
              </Text>
            </Section>

            {/* Stats Overview */}
            <Section className="mb-8">
              <div className="grid grid-cols-3 gap-4 text-center">
                <div className="bg-green-50 p-4 rounded-lg">
                  <Text className="text-2xl font-bold text-green-600 mb-1">
                    {completedCount}
                  </Text>
                  <Text className="text-sm text-green-700">Completed</Text>
                </div>
                <div className="bg-blue-50 p-4 rounded-lg">
                  <Text className="text-2xl font-bold text-blue-600 mb-1">
                    {activeCount}
                  </Text>
                  <Text className="text-sm text-blue-700">Active</Text>
                </div>
                <div className="bg-red-50 p-4 rounded-lg">
                  <Text className="text-2xl font-bold text-red-600 mb-1">
                    {overdueCount}
                  </Text>
                  <Text className="text-sm text-red-700">Overdue</Text>
                </div>
              </div>
            </Section>

            {/* Progress Bar */}
            <Section className="mb-8">
              <Text className="text-lg font-semibold text-gray-800 mb-2">
                Weekly Completion Rate: {completionRate}%
              </Text>
              <div className="w-full bg-gray-200 rounded-full h-2">
                <div
                  className={`h-2 rounded-full ${
                    completionRate >= 80
                      ? "bg-green-500"
                      : completionRate >= 60
                        ? "bg-blue-500"
                        : completionRate >= 40
                          ? "bg-yellow-500"
                          : "bg-red-500"
                  }`}
                  style={{ width: `${completionRate}%` }}
                ></div>
              </div>
            </Section>

            {/* Completed Todos */}
            {hasCompleted && completedTodos.length > 0 && (
              <Section className="mb-6">
                <Heading className="text-lg font-semibold text-green-700 mb-3">
                  ✅ Completed This Week ({completedCount})
                </Heading>
                <div className="space-y-2">
                  {completedTodos.slice(0, 5).map((todo, index) => (
                    <div
                      key={index}
                      className="flex items-center justify-between p-3 bg-green-50 rounded-lg"
                    >
                      <Text className="text-gray-800 flex-1">{todo.title}</Text>
                      <span
                        className={`px-2 py-1 rounded-full text-xs font-medium ${getPriorityBadgeColor(todo.priority)}`}
                      >
                        {todo.priority}
                      </span>
                    </div>
                  ))}
                  {completedTodos.length > 5 && (
                    <Text className="text-sm text-gray-600 text-center">
                      ... and {completedTodos.length - 5} more completed todos
                    </Text>
                  )}
                </div>
              </Section>
            )}

            {/* Overdue Todos */}
            {hasOverdue && overdueTodos.length > 0 && (
              <Section className="mb-6">
                <Heading className="text-lg font-semibold text-red-700 mb-3">
                  ⚠️ Needs Attention ({overdueCount} overdue)
                </Heading>
                <div className="space-y-2">
                  {overdueTodos.slice(0, 5).map((todo, index) => (
                    <div
                      key={index}
                      className="flex items-center justify-between p-3 bg-red-50 rounded-lg"
                    >
                      <div className="flex-1">
                        <Text className="text-gray-800">{todo.title}</Text>
                        {todo.dueDate && (
                          <Text className="text-sm text-red-600">
                            Due: {todo.dueDate}
                          </Text>
                        )}
                      </div>
                      <span
                        className={`px-2 py-1 rounded-full text-xs font-medium ${getPriorityBadgeColor(todo.priority)}`}
                      >
                        {todo.priority}
                      </span>
                    </div>
                  ))}
                  {overdueTodos.length > 5 && (
                    <Text className="text-sm text-gray-600 text-center">
                      ... and {overdueTodos.length - 5} more overdue todos
                    </Text>
                  )}
                </div>
              </Section>
            )}

            {/* Action Buttons */}
            <Section className="my-8 text-center">
              <Button
                className="bg-blue-600 hover:bg-blue-700 text-white font-medium rounded-md px-6 py-3 mr-4"
                href="/dashboard"
              >
                View Dashboard
              </Button>
              {hasOverdue && (
                <Button
                  className="bg-red-600 hover:bg-red-700 text-white font-medium rounded-md px-6 py-3"
                  href="/todos?filter=overdue"
                >
                  Review Overdue
                </Button>
              )}
            </Section>

            {/* Weekly Tips */}
            <Section className="bg-blue-50 border-l-4 border-blue-400 p-4 mb-6">
              <Text className="text-blue-800 text-base font-medium mb-2">
                💡 Productivity Tip
              </Text>
              <Text className="text-blue-700 text-sm">
                {completionRate >= 80
                  ? "You're on fire! Keep up this momentum by planning next week's priorities."
                  : completionRate >= 60
                    ? "Good progress! Try time-blocking your most important tasks for better focus."
                    : parseInt(overdueCount) > 0
                      ? "Focus on clearing overdue items first, then plan realistic due dates for new tasks."
                      : "Start your week by identifying 3 key priorities and tackle them first."}
              </Text>
            </Section>

            <Hr className="border-gray-200 my-6" />

            <Section>
              <Text className="text-gray-600 text-sm">
                This is your weekly productivity summary.{" "}
                <Link
                  href="/settings/notifications"
                  className="text-blue-600 underline"
                >
                  Manage notification preferences
                </Link>{" "}
                or{" "}
                <Link href="/dashboard" className="text-blue-600 underline">
                  view your full dashboard
                </Link>
                .
              </Text>
            </Section>

            <Section className="mt-8 text-center">
              <Text className="text-gray-500 text-xs">
                © {new Date().getFullYear()} Tasker. All rights reserved.
              </Text>
            </Section>
          </Container>
        </Body>
      </Tailwind>
    </Html>
  );
};

WeeklyReportEmail.PreviewProps = {
  weekStart: "January 8, 2024",
  weekEnd: "January 14, 2024",
  completedCount: "8",
  activeCount: "12",
  overdueCount: "3",
  completedTodos: [
    { id: "1", title: "Complete quarterly budget review", priority: "high" },
    { id: "2", title: "Update project documentation", priority: "medium" },
    { id: "3", title: "Schedule team one-on-ones", priority: "low" },
  ],
  overdueTodos: [
    {
      id: "4",
      title: "Submit expense reports",
      priority: "high",
      dueDate: "January 10, 2024",
    },
    {
      id: "5",
      title: "Review marketing proposals",
      priority: "medium",
      dueDate: "January 12, 2024",
    },
  ],
  hasCompleted: true,
  hasOverdue: true,
};

export default WeeklyReportEmail;
