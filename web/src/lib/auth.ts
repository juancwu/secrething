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
	user_id: string;
	email: string;
	name?: string;
	access_token: string;
	refresh_token?: string; // Only present for CLI clients
	expires_in: number;
}

// Register
export interface RegisterRequest {
	email: string;
	password: string;
	nickname: string;
}

export interface RegisterResponse {
	user_id: string;
	email: string;
	name?: string;
	access_token: string;
	refresh_token?: string; // Only present for CLI clients
	expires_in: number;
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
	LOGIN: "/api/auth/sign-in",
	REGISTER: "/api/auth/sign-up",
	CHECK_TOKEN: "/api/auth/token/check", // This will need to be updated to match server endpoint
	TOTP_SETUP: "/api/auth/totp/activate",
	TOTP_LOCK: "/api/auth/totp/verify",
	TOTP_DELETE: "/api/auth/totp/remove",
	VERIFY_EMAIL: "/api/auth/email/verify", // This will need to be updated to match server endpoint
	RESEND_VERIFICATION: "/api/auth/email/resend-verification", // This will need to be updated to match server endpoint
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
