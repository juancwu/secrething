import { SignInPage } from "@/pages/signin.page";
import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/signin")({
	component: SignInPage,
});
