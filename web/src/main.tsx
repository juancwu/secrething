import { MantineProvider, createTheme } from "@mantine/core";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { ReactQueryDevtools } from "@tanstack/react-query-devtools";
import { RouterProvider, createRouter } from "@tanstack/react-router";
import { StrictMode } from "react";
import { createRoot } from "react-dom/client";

import "@mantine/core/styles.css";

import { AuthProvider } from "./providers/auth/provider";
// Import the generated route tree
import { routeTree } from "./routeTree.gen";

// Create a new router instance
const router = createRouter({ routeTree });

// Register the router instance for type safety
declare module "@tanstack/react-router" {
	interface Register {
		router: typeof router;
	}
}

const theme = createTheme({
	primaryColor: "violet",
});

const rootElement = document.getElementById("root");

if (!rootElement) {
	console.error("Root element not found. Cannot mount React application.");
} else {
	const root = createRoot(rootElement);
	const queryClient = new QueryClient();
	root.render(
		<StrictMode>
			<QueryClientProvider client={queryClient}>
				<AuthProvider>
					<MantineProvider theme={theme}>
						<RouterProvider router={router} />
						<ReactQueryDevtools />
					</MantineProvider>
				</AuthProvider>
			</QueryClientProvider>
		</StrictMode>,
	);
}
