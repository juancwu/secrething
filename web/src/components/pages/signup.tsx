import { Anchor } from "@/components/ui/anchor";
import { useAuth } from "@/contexts/auth";
import { AuthError } from "@/lib/api/errors";
import {
	Box,
	Button,
	Center,
	Flex,
	Grid,
	Group,
	PasswordInput,
	Progress,
	Stack,
	Text,
	TextInput,
	Title,
} from "@mantine/core";
import { useForm } from "@mantine/form";
import { notifications } from "@mantine/notifications";
import { IconCheck, IconX } from "@tabler/icons-react";
import { zodResolver } from "mantine-form-zod-resolver";
import { useState } from "react";
import { z } from "zod";

const PASSWORD_MIN_LEN = 8;

function PasswordRequirement({
	meets,
	label,
}: { meets: boolean; label: string }) {
	return (
		<Text component="div" c={meets ? "teal" : "red"} mt={5} size="sm">
			<Center inline>
				{meets ? (
					<IconCheck size={14} stroke={1.5} />
				) : (
					<IconX size={14} stroke={1.5} />
				)}
				<Box ml={7}>{label}</Box>
			</Center>
		</Text>
	);
}

const requirements = [
	{ re: /[a-z]/, label: "Includes lowercase letter" },
	{ re: /[A-Z]/, label: "Includes uppercase letter" },
	{ re: /[0-9]/, label: "Includes number" },
	{ re: /[$&+,:;=?@#|'<>.^*()%!-]/, label: "Includes special symbol" },
];

function getStrength(password: string) {
	let multiplier = password.length >= PASSWORD_MIN_LEN ? 0 : 1;

	for (const requirement of requirements) {
		if (!requirement.re.test(password)) {
			multiplier += 1;
		}
	}

	return Math.max(100 - (100 / (requirements.length + 1)) * multiplier, 0);
}

const passwordSchema = z
	.string()
	.nonempty({ message: "Password can't be empty." })
	.min(
		PASSWORD_MIN_LEN,
		`Password must be at least ${PASSWORD_MIN_LEN} characters long.`,
	)
	.regex(
		requirements[0].re,
		"Password must contain at least one lowercase letter.",
	)
	.regex(
		requirements[1].re,
		"Password must contain at least one uppercase letter.",
	)
	.regex(requirements[2].re, "Password must contain at least one number.")
	.regex(
		requirements[3].re,
		"Password must contain at least one special character.",
	);

const schema = z
	.object({
		first_name: z.string().nonempty("First name is empty."),
		last_name: z.string().nonempty("Last name is empty."),
		email: z.string().email("Invalid email."),
		password: passwordSchema,
		confirmPassword: z.string().nonempty("Confirm password can't be empty."),
	})
	.refine((data) => data.password === data.confirmPassword, {
		message: "Passwords don't match.",
		path: ["confirmPassword"],
	});

export function SignUpPage() {
	const [isLoading, setIsLoading] = useState(false);
	const auth = useAuth();
	const formProps = useForm({
		mode: "controlled",
		initialValues: {
			first_name: "",
			last_name: "",
			email: "",
			password: "",
			confirmPassword: "",
		},
		validate: zodResolver(schema),
	});

	const strength = getStrength(formProps.values.password);
	const checks = requirements.map((requirement, index) => (
		<PasswordRequirement
			key={index}
			label={requirement.label}
			meets={requirement.re.test(formProps.values.password)}
		/>
	));
	const bars = Array(4)
		.fill(0)
		.map((_, index) => (
			<Progress
				styles={{ section: { transitionDuration: "0ms" } }}
				value={
					formProps.values.password.length > 0 && index === 0
						? 100
						: strength >= ((index + 1) / 4) * 100
							? 100
							: 0
				}
				color={strength > 80 ? "teal" : strength > 50 ? "yellow" : "red"}
				key={index}
				size={4}
			/>
		));

	return (
		<Flex h="100vh" w="100%" align="center" justify="center">
			<Stack w={{ base: "90%", xs: "450px" }}>
				<Title order={5}>Join Secrething community</Title>
				<form
					onSubmit={formProps.onSubmit(async (values) => {
						try {
							setIsLoading(true);
							if (auth.signup) {
								await auth.signup(values);
							}
						} catch (error) {
							let message =
								"Oops, something went wrong. Please try again later.";
							if (error instanceof AuthError) {
								if (error.data.errors) {
									formProps.setErrors(error.data.errors);
								}
								message = error.data.message || error.message;
							}
							notifications.show({
								title: "Sign Up Failure",
								color: "red",
								message,
							});
						} finally {
							setIsLoading(false);
						}
					})}
				>
					<Stack>
						<Grid>
							<Grid.Col span={6}>
								<TextInput
									withAsterisk
									label="First name"
									placeholder="John"
									autoComplete="given-name"
									key={formProps.key("first_name")}
									required
									{...formProps.getInputProps("first_name")}
								/>
							</Grid.Col>
							<Grid.Col span={6}>
								<TextInput
									withAsterisk
									label="Last name"
									placeholder="Doe"
									autoComplete="family-name"
									key={formProps.key("last_name")}
									required
									{...formProps.getInputProps("last_name")}
								/>
							</Grid.Col>
						</Grid>
						<TextInput
							withAsterisk
							label="Email"
							type="email"
							placeholder="your@mail.com"
							autoComplete="email"
							key={formProps.key("email")}
							required
							{...formProps.getInputProps("email")}
						/>
						<div>
							<PasswordInput
								withAsterisk
								label="Password"
								placeholder="Your password"
								key={formProps.key("password")}
								autoComplete="new-password"
								required
								{...formProps.getInputProps("password")}
							/>
							<Group gap={5} grow mt="xs" mb="md">
								{bars}
							</Group>
							<PasswordRequirement
								label={`Has at least ${PASSWORD_MIN_LEN} characters`}
								meets={formProps.values.password.length >= PASSWORD_MIN_LEN}
							/>
							{checks}
						</div>
						<PasswordInput
							withAsterisk
							label="Confirm password"
							placeholder="You password, just one more time"
							key={formProps.key("confirmPassword")}
							autoComplete="new-password"
							required
							{...formProps.getInputProps("confirmPassword")}
						/>
						<Button type="submit" loading={isLoading}>
							Create Account
						</Button>
					</Stack>
					<Text size="sm" c="dimmed" mt="sm">
						Have an account already? <Anchor to="/signin">Sign in</Anchor>
					</Text>
				</form>
			</Stack>
		</Flex>
	);
}
