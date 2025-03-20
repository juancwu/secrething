import { LoginForm } from "@/components/forms/login-form";

function LoginPage() {
	return (
		<div className="min-h-svh flex flex-1 items-center justify-center">
			<div className="w-full max-w-xs">
				<LoginForm />
			</div>
		</div>
	);
}

export default LoginPage;
