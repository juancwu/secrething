import { useContext } from "react";
import { AuthProviderContext } from "./context";

export function useAuth() {
	const context = useContext(AuthProviderContext);

	if (context === undefined) {
		throw new Error("useAuth must be used within an AuthProvider");
	}

	return context;
}
