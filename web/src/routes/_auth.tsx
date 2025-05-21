import { Header } from "@/components/ui/header";
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
		<Box h="full">
			<Header />
			<Outlet />
		</Box>
	);
}
