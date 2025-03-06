import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import { ThemeProvider } from "@/providers/theme-provider.tsx";
import { RouterProvider, createRouter } from "@tanstack/react-router";
import { routeTree } from "./routeTree.gen";

import "./index.css";

// Create a new router instance
const router = createRouter({ routeTree });

// Register the router instance for type safety
declare module "@tanstack/react-router" {
	interface Register {
		router: typeof router;
	}
}

const root = document.getElementById("root");
if (root === null) {
	throw new Error("No root element found");
}
createRoot(root).render(
	<StrictMode>
		<ThemeProvider defaultTheme="system" storageKey="konbini-theme">
			<RouterProvider router={router} />
		</ThemeProvider>
	</StrictMode>,
);
