import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import { ThemeProvider } from "@/providers/theme-provider.tsx";
import { RouterProvider, createRouter } from "@tanstack/react-router";
import { routeTree } from "./routeTree.gen";

import { AuthProvider } from "@/providers/auth";
import { useAuth } from "@/providers/auth/useAuth";

import "./index.css";

// Create a new router instance
const router = createRouter({
	routeTree,
	context: {
		auth: undefined,
	},
});

function App() {
	const auth = useAuth();
	return <RouterProvider router={router} context={{ auth }} />;
}

const root = document.getElementById("root");
if (root === null) {
	throw new Error("No root element found");
}
createRoot(root).render(
	<StrictMode>
		<ThemeProvider defaultTheme="system" storageKey="konbini-theme">
			<AuthProvider>
				<App />
			</AuthProvider>
		</ThemeProvider>
	</StrictMode>,
);
