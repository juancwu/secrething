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
import { createFileRoute, redirect, useRouter } from "@tanstack/react-router";
import { zodResolver } from "mantine-form-zod-resolver";
import { useState } from "react";
import { z } from "zod";

const schema = z.object({
	email: z.string().email("Invalid email."),
	password: z.string().nonempty("Password can't be empty."),
});

const fallback = "/dashboard" as const;

export const Route = createFileRoute("/_header/signin")({
	validateSearch: z.object({
		redirect: z.string().optional().catch(""),
	}),
	beforeLoad: ({ context, search }) => {
		if (context.auth?.user) {
			throw redirect({ to: search.redirect || fallback });
		}
	},
	component: SignInPage,
});

function SignInPage() {
	const [isLoading, setIsLoading] = useState(false);
	const search = Route.useSearch();
	const router = useRouter();
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
		<Flex w="100%" align="center" justify="center">
			<Stack w={{ base: "90%", xs: "450px" }}>
				<Title order={5}>Sign in to your account</Title>
				<form
					onSubmit={formProps.onSubmit(async (values) => {
						try {
							setIsLoading(true);
							if (auth.signin) {
								await auth.signin(values);
								await router.invalidate();
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
						Don't have an account?{" "}
						<Anchor to="/signup" search={search}>
							Create account
						</Anchor>
					</Text>
				</form>
			</Stack>
		</Flex>
	);
}
