import { createFileRoute } from "@tanstack/react-router";
import LoginPage from "@/pages/login";

export const Route = createFileRoute("/__noauth/login")({
	component: LoginPage,
});
