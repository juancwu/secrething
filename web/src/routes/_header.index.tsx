import { Space } from "@mantine/core";
import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/_header/")({
	component: Index,
});

function Index() {
	return (
		<>
			<Space />
			<p>Yeah... nothing to see here, just sign in or sign up.</p>
		</>
	);
}
