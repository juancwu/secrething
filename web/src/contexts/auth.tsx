import { signin, signup } from "@/lib/api/auth";
import type { SigninBody, SignupBody } from "@/lib/api/types";
import type { User } from "@/lib/types/auth";
import type { ReactNode } from "@tanstack/react-router";
import { useCallback, useContext, useState } from "react";
import { createContext, useMemo } from "react";
import { flushSync } from "react-dom";

export type AuthContextValue = {
	isLoading: boolean;
	signup?: typeof signup;
	signin?: typeof signin;
	user?: User;
	token?: string;
	expiresAt?: number;
};

export const AuthContext = createContext<AuthContextValue>({
	isLoading: true,
});

export function AuthProvider({ children }: { children: ReactNode }) {
	const [isLoading, setIsLoading] = useState(true);
	const [user, setUser] = useState<User | undefined>();
	const [token, setToken] = useState<string | undefined>();
	const [expiresAt, setExpiresAt] = useState<number | undefined>();

	const handleSignup = useCallback(async (body: SignupBody) => {
		const res = await signup(body);
		flushSync(() => {
			setUser(res.user);
			setToken(res.token);
			setExpiresAt(res.expires_at);
		});
		return res;
	}, []);

	const handleSignin = useCallback(async (body: SigninBody) => {
		const res = await signin(body);
		flushSync(() => {
			setUser(res.user);
			setToken(res.token);
			setExpiresAt(res.expires_at);
		});
		return res;
	}, []);

	const value = useMemo<AuthContextValue>(() => {
		return {
			isLoading,
			user,
			token,
			expiresAt,
			signup: handleSignup,
			signin: handleSignin,
		} satisfies AuthContextValue;
	}, [isLoading, handleSignup, user, token, expiresAt]);

	return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

export function useAuth() {
	const value = useContext(AuthContext);
	return value;
}
