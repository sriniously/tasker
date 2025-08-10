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

interface OverdueNotificationEmailProps {
  todoTitle: string;
  todoID: string;
  dueDate: string;
  daysOverdue: string;
}

export const OverdueNotificationEmail = ({
  todoTitle = "{{.TodoTitle}}",
  todoID = "{{.TodoID}}",
  dueDate = "{{.DueDate}}",
  daysOverdue = "{{.DaysOverdue}}",
}: OverdueNotificationEmailProps) => {
  const overdueMessage =
    parseInt(daysOverdue) === 1 ? "1 day overdue" : `${daysOverdue} days overdue`;

  return (
    <Html>
      <Head />
      <Preview>Overdue: "{todoTitle}" needs your attention</Preview>
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
                ⚠️ Overdue Todo
              </Heading>
            </Section>

            <Section className="bg-red-50 border-l-4 border-red-500 p-4 mb-6">
              <Text className="font-semibold text-red-600 text-lg mb-2">
                "{todoTitle}" is {overdueMessage}
              </Text>
              <Text className="text-gray-700 text-base">
                Was due: {dueDate}
              </Text>
            </Section>

            <Section>
              <Text className="text-gray-700 text-base">
                Your todo item is now overdue and needs immediate attention.
                Don't let important tasks fall behind schedule!
              </Text>
            </Section>

            <Section className="my-8 text-center">
              <Button
                className="bg-red-600 hover:bg-red-700 text-white font-medium rounded-md px-6 py-3 mr-4"
                href={`/todos?id=${todoID}`}
              >
                View Todo
              </Button>
              <Button
                className="bg-green-600 hover:bg-green-700 text-white font-medium rounded-md px-6 py-3"
                href={`/todos?id=${todoID}&action=complete`}
              >
                Mark Complete
              </Button>
            </Section>

            <Section className="bg-blue-50 border-l-4 border-blue-400 p-4 mb-6">
              <Text className="text-blue-800 text-base font-medium mb-2">
                💡 Need to reschedule?
              </Text>
              <Text className="text-blue-700 text-sm">
                If this todo is no longer relevant or needs a new timeline, you
                can:
              </Text>
              <ul className="list-disc pl-6 text-blue-700 text-sm mt-2">
                <li>Update the due date to a more realistic timeline</li>
                <li>Break it down into smaller, manageable tasks</li>
                <li>Archive it if it's no longer needed</li>
              </ul>
            </Section>

            <Section>
              <Text className="text-gray-700 text-base">
                🎯 <strong>Stay organized:</strong> Regular review of your todos
                helps prevent items from becoming overdue. Consider setting
                aside time each week to review and prioritize your tasks.
              </Text>
            </Section>

            <Hr className="border-gray-200 my-6" />

            <Section>
              <Text className="text-gray-600 text-sm">
                You're receiving this notification because you have an overdue
                todo item.{" "}
                <Link
                  href={`/settings/notifications`}
                  className="text-blue-600 underline"
                >
                  Manage notification preferences
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

OverdueNotificationEmail.PreviewProps = {
  todoTitle: "Review budget proposal",
  todoID: "123e4567-e89b-12d3-a456-426614174000",
  dueDate: "Friday, January 12, 2024 at 3:00 PM",
  daysOverdue: "3",
};

export default OverdueNotificationEmail;
