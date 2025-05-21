import type { SigninBody, SignupBody, AuthResponse } from "@/lib/api/types";
import { post } from "./common";
import { AuthError } from "./errors";

export async function signup(body: SignupBody): Promise<AuthResponse> {
	const res = await post("auth/signup", body);
	const data = await res.json();
	if (!res.ok) {
		throw new AuthError(
			`Signup failed with status: ${res.status}`,
			"SignUp",
			data,
		);
	}
	return data;
}

export async function signin(body: SigninBody): Promise<AuthResponse> {
	const res = await post("auth/signin", body);
	const data = await res.json();
	if (!res.ok) {
		throw new AuthError(
			`Signin failed with status: ${res.status}`,
			"SignIn",
			data,
		);
	}
	return data;
}
