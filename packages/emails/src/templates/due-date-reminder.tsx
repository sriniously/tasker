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

interface DueDateReminderEmailProps {
  todoTitle: string;
  todoID: string;
  dueDate: string;
  daysUntilDue: string;
}

export const DueDateReminderEmail = ({
  todoTitle = "{{.TodoTitle}}",
  todoID = "{{.TodoID}}",
  dueDate = "{{.DueDate}}",
  daysUntilDue = "{{.DaysUntilDue}}",
}: DueDateReminderEmailProps) => {
  const urgencyColor =
    parseInt(daysUntilDue) <= 1 ? "text-red-600" : "text-orange-600";
  const urgencyMessage =
    parseInt(daysUntilDue) <= 1 ? "due tomorrow" : `due in ${daysUntilDue} days`;

  return (
    <Html>
      <Head />
      <Preview>
        Reminder: "{todoTitle}" is {urgencyMessage}
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
                📅 Todo Reminder
              </Heading>
            </Section>

            <Section className="bg-yellow-50 border-l-4 border-yellow-400 p-4 mb-6">
              <Text className={`font-semibold ${urgencyColor} text-lg mb-2`}>
                "{todoTitle}" is {urgencyMessage}
              </Text>
              <Text className="text-gray-700 text-base">
                Due Date: {dueDate}
              </Text>
            </Section>

            <Section>
              <Text className="text-gray-700 text-base">
                This is a friendly reminder that your todo item is due soon.
                Don't let it slip through the cracks!
              </Text>
            </Section>

            <Section className="my-8 text-center">
              <Button
                className="bg-blue-600 hover:bg-blue-700 text-white font-medium rounded-md px-6 py-3 mr-4"
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

            <Section>
              <Text className="text-gray-700 text-base">
                💡 <strong>Pro tip:</strong> Stay on top of your tasks by
                checking your Tasker dashboard regularly and setting realistic
                due dates.
              </Text>
            </Section>

            <Hr className="border-gray-200 my-6" />

            <Section>
              <Text className="text-gray-600 text-sm">
                You're receiving this reminder because you have an active todo
                item with an upcoming due date.{" "}
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

DueDateReminderEmail.PreviewProps = {
  todoTitle: "Complete quarterly report",
  todoID: "123e4567-e89b-12d3-a456-426614174000",
  dueDate: "Monday, January 15, 2025 at 5:00 PM",
  daysUntilDue: "1",
};

export default DueDateReminderEmail;
