import { createContext } from "react";
import type { AuthProviderState } from "./types";

// This is just a placeholder, the actual implementation will be in the provider component
const noop = async () => {
	throw new Error("Auth Provider not initialized");
};

const initialState: AuthProviderState = {
	user: null,
	isLoading: false,
	isAuthenticated: false,
	isPartialAuth: false,
	login: noop,
	register: noop,
	logout: () => {},
	checkSession: noop,
	setupTOTP: noop,
	lockTOTP: noop,
	removeTOTP: noop,
	resendVerificationEmail: noop,
};

export const AuthProviderContext =
	createContext<AuthProviderState>(initialState);
