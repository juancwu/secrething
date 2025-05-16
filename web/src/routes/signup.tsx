import { SignUpPage } from "@/pages/signup.page";
import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/signup")({
	component: SignUpPage,
});
