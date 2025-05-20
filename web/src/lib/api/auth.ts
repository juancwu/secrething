import type { SigninBody, SignupBody, AuthResponse } from "@/lib/api/types";
import { post } from "./common";
import { AuthError } from "./errors";

export async function signup(body: SignupBody): Promise<AuthResponse> {
	try {
		const res = await post("auth/signup", body);
		return res.json();
	} catch (error) {
		throw new AuthError("failed to signup", "SignUp");
	}
}

export async function signin(body: SigninBody): Promise<AuthResponse> {
	try {
		const res = await post("auth/signin", body);
		return res.json();
	} catch (error) {
		throw new AuthError("failed to signin", "SignIn");
	}
}
