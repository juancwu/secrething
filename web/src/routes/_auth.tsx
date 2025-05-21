import { Box } from "@mantine/core";
import { Outlet, createFileRoute, redirect } from "@tanstack/react-router";

export const Route = createFileRoute("/_auth")({
	beforeLoad: ({ context, location }) => {
		if (!context.auth?.isAuthenticated) {
			throw redirect({
				to: "/signin",
				search: {
					redirect: location.href,
				},
			});
		}
	},
	component: AuthLayout,
});

function AuthLayout() {
	return (
		<Box p={2} h="full">
			<h1>Authenticated Route</h1>
			<p>This route's content is only visible to authenticated users.</p>
			<Outlet />
		</Box>
	);
}
