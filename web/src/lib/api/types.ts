import type { User } from "@/lib/types/auth";

export type ApiResponse = {
	message: string;
	success: boolean;
};

export type SignupBody = {
	first_name: string;
	last_name: string;
	email: string;
	password: string;
};

export type AuthResponse = {
	token: string;
	expires_at: number;
	user: User;
} & ApiResponse;

export type SigninBody = {
	email: string;
	password: string;
};
