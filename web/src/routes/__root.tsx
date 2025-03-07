import { Navbar } from "@/components/ui/navbar";
import { Toaster } from "@/components/ui/sonner";
import { AuthProviderState } from "@/providers/auth/types";
import { createRootRouteWithContext, Outlet } from "@tanstack/react-router";

interface RouterContext {
	auth?: AuthProviderState;
}

export const Route = createRootRouteWithContext<RouterContext>()({
	component: () => (
		<>
			<Navbar />
			<Outlet />
			<Toaster />
		</>
	),
});
