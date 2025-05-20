import { Anchor } from "@/components/ui/anchor";
import { useAuth } from "@/contexts/auth";
import { AuthError } from "@/lib/api/errors";
import {
	Button,
	Flex,
	PasswordInput,
	Stack,
	Text,
	TextInput,
	Title,
} from "@mantine/core";
import { useForm } from "@mantine/form";
import { notifications } from "@mantine/notifications";
import { zodResolver } from "mantine-form-zod-resolver";
import { z } from "zod";

const schema = z.object({
	email: z.string().email("Invalid email."),
	password: z.string().nonempty("Password can't be empty."),
});

export function SignInPage() {
	const auth = useAuth();
	const formProps = useForm({
		mode: "uncontrolled",
		initialValues: {
			email: "",
			password: "",
		},
		validate: zodResolver(schema),
	});

	return (
		<Flex h="100vh" w="100%" align="center" justify="center">
			<Stack w={{ base: "90%", xs: "450px" }}>
				<Title order={5}>Sign in to your account</Title>
				<form
					onSubmit={formProps.onSubmit(async (values) => {
						try {
							if (auth.signin) {
								await auth.signin(values);
							}
						} catch (error) {
							let message =
								"Oops, something went wrong. Please try again later.";
							if (error instanceof AuthError) {
								message = error.message;
							}
							notifications.show({
								title: "Sign In Failure",
								color: "red",
								message,
							});
						}
					})}
				>
					<Stack>
						<TextInput
							withAsterisk
							label="Email"
							placeholder="your@mail.com"
							key={formProps.key("email")}
							required
							{...formProps.getInputProps("email")}
						/>
						<PasswordInput
							withAsterisk
							label="Password"
							key={formProps.key("password")}
							required
							{...formProps.getInputProps("password")}
						/>
						<Button type="submit">Sign In</Button>
					</Stack>
					<Text size="sm" c="dimmed" mt="sm">
						Don't have an account? <Anchor to="/signup">Create account</Anchor>
					</Text>
				</form>
			</Stack>
		</Flex>
	);
}
