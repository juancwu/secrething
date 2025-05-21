import { useAuth } from "@/contexts/auth";
import { Button } from "@mantine/core";
import { createFileRoute, useRouter } from "@tanstack/react-router";

export const Route = createFileRoute("/_auth/dashboard")({
	component: RouteComponent,
});

function RouteComponent() {
	const router = useRouter();
	const auth = useAuth();

	const handleSignout = () => {
		if (auth.signout && window.confirm("Are you sure you want to sign out?")) {
			auth.signout().then(() => {
				router.invalidate();
			});
		}
	};

	return <Button onClick={handleSignout}>Sign Out</Button>;
}
