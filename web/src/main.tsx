import { AuthProvider, useAuth } from "@/contexts/auth";
import { MantineProvider, createTheme } from "@mantine/core";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { ReactQueryDevtools } from "@tanstack/react-query-devtools";
import { RouterProvider, createRouter } from "@tanstack/react-router";
import { StrictMode } from "react";
import { createRoot } from "react-dom/client";

import "@mantine/core/styles.css";
import "@mantine/notifications/styles.css";

// Import the generated route tree
import { routeTree } from "@/routeTree.gen";
import { Notifications } from "@mantine/notifications";

// Create a new router instance
const router = createRouter({
	routeTree,
	defaultPreload: "intent",
	scrollRestoration: true,
	context: {
		auth: undefined,
	},
});

// Register the router instance for type safety
declare module "@tanstack/react-router" {
	interface Register {
		router: typeof router;
	}
}

const theme = createTheme({
	primaryColor: "violet",
});

const queryClient = new QueryClient();

function InnerApp() {
	const auth = useAuth();
	return (
		<>
			<Notifications />
			<RouterProvider router={router} context={{ auth }} />
			<ReactQueryDevtools />
		</>
	);
}

function App() {
	return (
		<QueryClientProvider client={queryClient}>
			<AuthProvider>
				<MantineProvider theme={theme}>
					<InnerApp />
				</MantineProvider>
			</AuthProvider>
		</QueryClientProvider>
	);
}

const rootElement = document.getElementById("root");

if (!rootElement) {
	console.error("Root element not found. Cannot mount React application.");
} else {
	const root = createRoot(rootElement);
	root.render(
		<StrictMode>
			<App />
		</StrictMode>,
	);
}
