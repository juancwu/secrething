import { useContext } from "react";
import { AuthContext } from "./provider";

export function useAuth() {
	const value = useContext(AuthContext);
	return value;
}
