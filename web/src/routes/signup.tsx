import { SignUpPage } from "@/components/pages/signup";
import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/signup")({
	component: SignUpPage,
});
