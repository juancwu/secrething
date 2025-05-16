import type { ReactNode } from "@tanstack/react-router";
import { createContext, useMemo } from "react";

export type AuthContextValue = {
	isLoading: boolean;
};

export const AuthContext = createContext<AuthContextValue>({ isLoading: true });

export function AuthProvider({ children }: { children: ReactNode }) {
	const value = useMemo<AuthContextValue>(() => {
		return { isLoading: true };
	}, []);

	return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}
