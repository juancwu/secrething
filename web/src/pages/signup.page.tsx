import { Anchor } from "@/components/ui/anchor";
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
import { zodResolver } from "mantine-form-zod-resolver";
import { z } from "zod";

const passwordSchema = z
	.string()
	.nonempty({ message: "Password can't be empty." })
	.min(8, "Password must be at least 8 characters long.")
	.regex(/[A-Z]/, "Password must contain at least one uppercase letter.")
	.regex(/[a-z]/, "Password must contain at least one lowercase letter.")
	.regex(/[0-9]/, "Password must contain at least one number.")
	.regex(
		/[^A-Za-z0-9]/,
		"Password must contain at least one special character.",
	);

const schema = z
	.object({
		name: z.string().min(2, "Name must be at least 3 characters long."),
		email: z.string().email("Invalid email."),
		password: passwordSchema,
		confirmPassword: z.string().nonempty("Confirm password can't be empty."),
	})
	.refine((data) => data.password === data.confirmPassword, {
		message: "Passwords don't match.",
		path: ["confirmPassword"],
	});

export function SignUpPage() {
	const formProps = useForm({
		mode: "uncontrolled",
		initialValues: {
			name: "",
			email: "",
			password: "",
			confirmPassword: "",
		},
		validate: zodResolver(schema),
	});

	return (
		<Flex h="100vh" w="100%" align="center" justify="center">
			<Stack w={{ base: "90%", xs: "450px" }}>
				<Title order={5}>Join Secrething community</Title>
				<form onSubmit={formProps.onSubmit((values) => console.log(values))}>
					<Stack>
						<TextInput
							withAsterisk
							label="Name"
							placeholder="John Doe"
							autoComplete="name"
							key={formProps.key("name")}
							required
							{...formProps.getInputProps("name")}
						/>
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
						<PasswordInput
							withAsterisk
							label="Password"
							key={formProps.key("password")}
							autoComplete="new-password"
							required
							{...formProps.getInputProps("password")}
						/>
						<PasswordInput
							withAsterisk
							label="Confirm password"
							key={formProps.key("confirmPassword")}
							autoComplete="new-password"
							required
							{...formProps.getInputProps("confirmPassword")}
						/>
						<Button type="submit">Create Account</Button>
					</Stack>
					<Text size="sm" c="dimmed" mt="sm">
						Have an account already? <Anchor to="/signin">Sign in</Anchor>
					</Text>
				</form>
			</Stack>
		</Flex>
	);
}
