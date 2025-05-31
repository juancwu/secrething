import type { User } from "@/lib/types/auth";

export type ApiResponse = {
	message: string;
};

export type ApiErrorResponse = {
	errors?: Record<string, string>;
} & ApiResponse;

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
};

export type SigninBody = {
	email: string;
	password: string;
	remember_me: boolean;
};
