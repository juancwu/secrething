import { createFileRoute } from "@tanstack/react-router";
import Register from "@/pages/register";

export const Route = createFileRoute("/__noauth/register")({
	component: Register,
});
