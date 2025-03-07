import { createFileRoute, redirect } from "@tanstack/react-router";

export const Route = createFileRoute("/__noauth")({
	beforeLoad: ({ context }) => {
		if (context.auth?.user) {
			throw redirect({
				to: "/dashboard",
				replace: true,
			});
		}
	},
});
