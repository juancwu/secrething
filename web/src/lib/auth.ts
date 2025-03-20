export type User = {
	id: string;
	email: string;
	emailVerified: boolean;
	nickname: string;
	totpEnabled: boolean;
};

export type TokenType = "partial" | "full";

export type AuthToken = {
	token: string;
	type: TokenType;
};

// Login
export interface LoginRequest {
	email: string;
	password: string;
	totp_code?: string;
}

export interface LoginResponse {
	token: string;
	type: TokenType;
}

// Register
export interface RegisterRequest {
	email: string;
	password: string;
	nickname: string;
}

export interface RegisterResponse {
	token: string;
	type: TokenType;
}

// Check Token
export interface CheckTokenRequest {
	auth_token: string;
}

export interface CheckTokenResponse {
	token: string;
	type: TokenType;
	email: string;
	email_verified: boolean;
	totp: boolean;
	nickname: string;
}

// TOTP Setup
export interface SetupTOTPResponse {
	url: string;
}

// TOTP Lock
export interface LockTOTPRequest {
	code: string;
}

export interface LockTOTPResponse {
	recovery_codes: string[];
	token: string;
	type: TokenType;
}

// TOTP Delete
export interface RemoveTOTPRequest {
	code: string;
}

// API routes
export const API_ROUTES = {
	LOGIN: "/auth/login",
	REGISTER: "/auth/register",
	CHECK_TOKEN: "/auth/token/check",
	TOTP_SETUP: "/auth/totp/setup",
	TOTP_LOCK: "/auth/totp/lock",
	TOTP_DELETE: "/auth/totp",
	VERIFY_EMAIL: "/auth/email/verify",
	RESEND_VERIFICATION: "/auth/email/resend-verification",
};

// Local storage keys
export const STORAGE_KEYS = {
	AUTH_TOKEN: "konbini_auth_token",
};

// Helper functions
export function saveToken(token: string): void {
	localStorage.setItem(STORAGE_KEYS.AUTH_TOKEN, token);
}

export function getToken(): string | null {
	return localStorage.getItem(STORAGE_KEYS.AUTH_TOKEN);
}

export function removeToken(): void {
	localStorage.removeItem(STORAGE_KEYS.AUTH_TOKEN);
}

export function transformResponseToUser(response: CheckTokenResponse): User {
	return {
		id: "", // Not provided in the response, might need to be handled differently
		email: response.email,
		emailVerified: response.email_verified,
		nickname: response.nickname,
		totpEnabled: response.totp,
	};
}
