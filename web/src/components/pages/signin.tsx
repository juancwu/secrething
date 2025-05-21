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
import { useState } from "react";
import { z } from "zod";

const schema = z.object({
	email: z.string().email("Invalid email."),
	password: z.string().nonempty("Password can't be empty."),
});

export function SignInPage() {
	const [isLoading, setIsLoading] = useState(false);
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
							setIsLoading(true);
							if (auth.signin) {
								await auth.signin(values);
							}
						} catch (error) {
							let showNotification = true;
							let message =
								"Oops, something went wrong. Please try again later.";
							if (error instanceof AuthError) {
								if (error.data.errors) {
									formProps.setErrors(error.data.errors);
									showNotification = false;
								} else if (error.data.message) {
									formProps.setErrors({
										email: error.data.message,
										password: error.data.message,
									});
									showNotification = false;
								}
								message = error.message;
							} else if (error instanceof Error) {
								message = error.message;
							}
							if (showNotification) {
								notifications.show({
									title: "Sign In Failure",
									color: "red",
									message,
								});
							}
						} finally {
							setIsLoading(false);
						}
					})}
				>
					<Stack>
						<TextInput
							withAsterisk
							label="Email"
							// type="email"
							autoComplete="email"
							placeholder="your@mail.com"
							key={formProps.key("email")}
							required
							{...formProps.getInputProps("email")}
						/>
						<PasswordInput
							withAsterisk
							label="Password"
							placeholder="Your password"
							autoComplete="current-password"
							key={formProps.key("password")}
							required
							{...formProps.getInputProps("password")}
						/>
						<Button type="submit" loading={isLoading}>
							Sign In
						</Button>
					</Stack>
					<Text size="sm" c="dimmed" mt="sm">
						Don't have an account? <Anchor to="/signup">Create account</Anchor>
					</Text>
				</form>
			</Stack>
		</Flex>
	);
}
