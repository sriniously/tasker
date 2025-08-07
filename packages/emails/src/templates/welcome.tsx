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

interface WelcomeEmailProps {
  userFirstName: string;
}

export const WelcomeEmail = ({
  userFirstName = "{{.UserFirstName}}",
}: WelcomeEmailProps) => {
  return (
    <Html>
      <Head />
      <Preview>Welcome to Boilerplate</Preview>
      <Tailwind>
        <Body className="bg-gray-100 font-sans">
          <Container className="bg-white p-8 rounded-lg shadow-sm my-10 mx-auto max-w-[600px]">
            <Heading className="text-2xl font-bold text-gray-800 mt-4">
              Welcome to Boilerplate!
            </Heading>

            <Section>
              <Text className="text-gray-700 text-base">
                Hi {userFirstName},
              </Text>
              <Text className="text-gray-700 text-base">
                Thank you for joining!
              </Text>
            </Section>

            <Section className="my-8 text-center">
              <Button
                className="bg-orange-600 hover:bg-orange-700 text-white font-medium rounded-md px-6 py-3"
                href={`/dashboard`}
              >
                Get Started
              </Button>
            </Section>

            <Hr className="border-gray-200 my-6" />

            <Section>
              <Text className="text-gray-600 text-sm">
                If you have any questions, feel free to{" "}
                <Link href={`/support`} className="text-orange-600 underline">
                  contact our support team
                </Link>
                .
              </Text>
            </Section>

            <Section className="mt-8 text-center">
              <Text className="text-gray-500 text-xs">
                © {new Date().getFullYear()} Alfred. All rights reserved.
              </Text>
              <Text className="text-gray-500 text-xs">
                123 Project Street, Suite 100, San Francisco, CA 94103
              </Text>
            </Section>
          </Container>
        </Body>
      </Tailwind>
    </Html>
  );
};

WelcomeEmail.PreviewProps = {
  userFirstName: "John",
};

export default WelcomeEmail;
