import { SignInPage } from "@/components/pages/signin";
import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/signin")({
	component: SignInPage,
});
