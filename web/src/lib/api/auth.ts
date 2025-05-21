import type {
	SigninBody,
	SignupBody,
	AuthResponse,
	ApiResponse,
} from "@/lib/api/types";
import { post } from "./common";
import { AuthError } from "./errors";

export async function signup(body: SignupBody): Promise<AuthResponse> {
	const res = await post("auth/signup", body);
	const data = await res.json();
	if (res.status !== 201) {
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

export async function signout(): Promise<ApiResponse> {
	const res = await post("auth/signout");
	const data = await res.json();
	if (!res.ok) {
		throw new AuthError(
			`Signout failed with status: ${res.status}`,
			"SignOut",
			data,
		);
	}
	return data;
}
