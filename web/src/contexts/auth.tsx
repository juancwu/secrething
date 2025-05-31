import { signin, signout, signup } from "@/lib/api/auth";
import type { SigninBody, SignupBody } from "@/lib/api/types";
import type { User } from "@/lib/types/auth";
import type { ReactNode } from "@tanstack/react-router";
import { useCallback, useContext, useRef, useState } from "react";
import { createContext, useMemo } from "react";
import { flushSync } from "react-dom";

export type AuthContextValue = {
	isLoading: boolean;
	isAuthenticated: boolean;
	signup?: typeof signup;
	signin?: typeof signin;
	signout?: typeof signout;
	user?: User;
	getAuthToken: () => string | undefined;
	expiresAt?: number;
};

export const AuthContext = createContext<AuthContextValue>({
	isLoading: true,
	isAuthenticated: false,
	getAuthToken: () => undefined,
});

export function AuthProvider({ children }: { children: ReactNode }) {
	const [isLoading, _setIsLoading] = useState(true);
	const [user, setUser] = useState<User | undefined>();
	const [expiresAt, setExpiresAt] = useState<number | undefined>();
	// NOTE: Using a ref to not cause re-render when the auth token refreshes
	//       The auth token acts as an access token.
	const authTokenRef = useRef<string | undefined>(undefined);

	const handleSignup = useCallback(async (body: SignupBody) => {
		const res = await signup(body);
		authTokenRef.current = res.token;
		flushSync(() => {
			setUser(res.user);
			setExpiresAt(res.expires_at);
		});
		return res;
	}, []);

	const handleSignin = useCallback(async (body: SigninBody) => {
		const res = await signin(body);
		authTokenRef.current = res.token;
		flushSync(() => {
			setUser(res.user);
			setExpiresAt(res.expires_at);
		});
		return res;
	}, []);

	const handleSignout = useCallback(async () => {
		const res = await signout();
		authTokenRef.current = undefined;
		flushSync(() => {
			setUser(undefined);
			setExpiresAt(undefined);
		});
		return res;
	}, []);

	const value = useMemo<AuthContextValue>(() => {
		return {
			isLoading,
			isAuthenticated: !!user,
			user,
			expiresAt,
			signup: handleSignup,
			signin: handleSignin,
			signout: handleSignout,
			getAuthToken: () => authTokenRef.current,
		} satisfies AuthContextValue;
	}, [isLoading, user, expiresAt, handleSignin, handleSignup, handleSignout]);

	return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

export function useAuth() {
	const value = useContext(AuthContext);
	if (!value) {
		throw new Error("useAuth must be used within an AuthProvider");
	}
	return value;
}
