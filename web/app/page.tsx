import { Button, Text, Title } from "@mantine/core";
import { ColorSchemeToggle } from "../components/ColorSchemeToggle/ColorSchemeToggle";
import { Welcome } from "../components/Welcome/Welcome";

export default function HomePage() {
  
  return (
    <>
      <Title>LibreLog</Title>
      <Text>
        Logging service for daily operations, allowing you to access and analyze
        your personal data whenever you'd like.
      </Text>
      <Button>Sign In</Button>
      <Button>Create Account</Button>
    </>
  );
}
