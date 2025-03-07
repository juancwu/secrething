import { ReactNode, useEffect, useMemo, useState } from "react";
import { AuthProviderContext } from "./context";
import type {
	User,
	LoginRequest,
	RegisterRequest,
	LockTOTPRequest,
	RemoveTOTPRequest,
	TokenType,
	CheckTokenResponse,
	LoginResponse,
	RegisterResponse,
	SetupTOTPResponse,
	LockTOTPResponse,
} from "@/lib/auth";
import {
	API_ROUTES,
	getToken,
	removeToken,
	saveToken,
	transformResponseToUser,
} from "@/lib/auth";
import { apiRequest } from "@/lib/utils";
import { toast } from "sonner";

interface AuthProviderProps {
	children: ReactNode;
}

export function AuthProvider({ children }: AuthProviderProps) {
	const [user, setUser] = useState<User | null>(null);
	const [isLoading, setIsLoading] = useState<boolean>(true);
	const [tokenType, setTokenType] = useState<TokenType | null>(null);

	// Derived state
	const isAuthenticated = !!user;
	const isPartialAuth = tokenType === "partial";

	// Auth functions
	const login = async (credentials: LoginRequest): Promise<void> => {
		setIsLoading(true);
		try {
			const response = await apiRequest<LoginResponse>(
				API_ROUTES.LOGIN,
				{
					method: "POST",
					body: JSON.stringify(credentials),
				},
				false,
			);

			saveToken(response.token);
			setTokenType(response.type);

			// After successful login, check the session to get user data
			await checkSession();
		} catch (error) {
			console.error("Login failed:", error);
			toast.error("Login Failed");
			throw error;
		} finally {
			setIsLoading(false);
		}
	};

	const register = async (data: RegisterRequest): Promise<void> => {
		setIsLoading(true);
		try {
			const response = await apiRequest<RegisterResponse>(
				API_ROUTES.REGISTER,
				{
					method: "POST",
					body: JSON.stringify(data),
				},
				false,
			);

			saveToken(response.token);
			setTokenType(response.type);

			// After successful registration, check the session to get user data
			await checkSession();
		} catch (error) {
			console.error("Registration failed:", error);
			throw error;
		} finally {
			setIsLoading(false);
		}
	};

	const logout = (): void => {
		removeToken();
		setUser(null);
		setTokenType(null);
	};

	const checkSession = async (): Promise<void> => {
		const token = getToken();
		if (!token) {
			setIsLoading(false);
			return;
		}

		setIsLoading(true);
		try {
			const response = await apiRequest<CheckTokenResponse>(
				API_ROUTES.CHECK_TOKEN,
				{
					method: "POST",
					body: JSON.stringify({ auth_token: token }),
				},
				false,
			);

			// Update token if a new one was returned
			if (response.token !== token) {
				saveToken(response.token);
			}

			setTokenType(response.type);
			setUser(transformResponseToUser(response));
		} catch (error) {
			console.error("Session check failed:", error);
			// Invalid token or other error, clear the session
			logout();
		} finally {
			setIsLoading(false);
		}
	};

	const setupTOTP = async (): Promise<SetupTOTPResponse> => {
		setIsLoading(true);
		try {
			return await apiRequest<SetupTOTPResponse>(API_ROUTES.TOTP_SETUP, {
				method: "POST",
			});
		} catch (error) {
			console.error("TOTP setup failed:", error);
			throw error;
		} finally {
			setIsLoading(false);
		}
	};

	const lockTOTP = async (data: LockTOTPRequest): Promise<LockTOTPResponse> => {
		setIsLoading(true);
		try {
			const response = await apiRequest<LockTOTPResponse>(
				API_ROUTES.TOTP_LOCK,
				{
					method: "POST",
					body: JSON.stringify(data),
				},
			);

			// Update token and token type
			saveToken(response.token);
			setTokenType(response.type);

			// Re-check session to update user data
			await checkSession();

			return response;
		} catch (error) {
			console.error("TOTP lock failed:", error);
			throw error;
		} finally {
			setIsLoading(false);
		}
	};

	const removeTOTP = async (data: RemoveTOTPRequest): Promise<void> => {
		setIsLoading(true);
		try {
			await apiRequest(API_ROUTES.TOTP_DELETE, {
				method: "DELETE",
				body: JSON.stringify(data),
			});

			// Re-check session to update user data
			await checkSession();
		} catch (error) {
			console.error("TOTP removal failed:", error);
			throw error;
		} finally {
			setIsLoading(false);
		}
	};

	const resendVerificationEmail = async (): Promise<void> => {
		setIsLoading(true);
		try {
			await apiRequest(API_ROUTES.RESEND_VERIFICATION, {
				method: "POST",
			});
		} catch (error) {
			console.error("Resend verification email failed:", error);
			throw error;
		} finally {
			setIsLoading(false);
		}
	};

	// Check for existing session on mount
	useEffect(() => {
		checkSession();
	}, []);

	// Create the context value
	const contextValue = useMemo(
		() => ({
			user,
			isLoading,
			isAuthenticated,
			isPartialAuth,
			login,
			register,
			logout,
			checkSession,
			setupTOTP,
			lockTOTP,
			removeTOTP,
			resendVerificationEmail,
		}),
		[user, isLoading, isAuthenticated, isPartialAuth, tokenType],
	);

	return (
		<AuthProviderContext.Provider value={contextValue}>
			{children}
		</AuthProviderContext.Provider>
	);
}
