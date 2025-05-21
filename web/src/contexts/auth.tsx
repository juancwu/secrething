import { signin, signout, signup } from "@/lib/api/auth";
import type { SigninBody, SignupBody } from "@/lib/api/types";
import type { User } from "@/lib/types/auth";
import type { ReactNode } from "@tanstack/react-router";
import { useCallback, useContext, useState } from "react";
import { createContext, useMemo } from "react";
import { flushSync } from "react-dom";

export type AuthContextValue = {
	isLoading: boolean;
	isAuthenticated: boolean;
	signup?: typeof signup;
	signin?: typeof signin;
	signout?: typeof signout;
	user?: User;
	token?: string;
	expiresAt?: number;
};

export const AuthContext = createContext<AuthContextValue>({
	isLoading: true,
	isAuthenticated: false,
});

export function AuthProvider({ children }: { children: ReactNode }) {
	const [isLoading, _setIsLoading] = useState(true);
	const [user, setUser] = useState<User | undefined>();
	const [token, setToken] = useState<string | undefined>();
	const [expiresAt, setExpiresAt] = useState<number | undefined>();
	const isAuthenticated = !!user;

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

	const handleSignout = useCallback(async () => {
		const res = await signout();
		flushSync(() => {
			setUser(undefined);
			setToken(undefined);
			setExpiresAt(undefined);
		});
		return res;
	}, []);

	const value = useMemo<AuthContextValue>(() => {
		return {
			isLoading,
			isAuthenticated,
			user,
			token,
			expiresAt,
			signup: handleSignup,
			signin: handleSignin,
			signout: handleSignout,
		} satisfies AuthContextValue;
	}, [
		isLoading,
		isAuthenticated,
		user,
		token,
		expiresAt,
		handleSignin,
		handleSignup,
		handleSignout,
	]);

	return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

export function useAuth() {
	const value = useContext(AuthContext);
	if (!value) {
		throw new Error("useAuth must be used within an AuthProvider");
	}
	return value;
}
