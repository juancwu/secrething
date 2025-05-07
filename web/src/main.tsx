import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import "./index.css";
import { ClerkProvider } from "@clerk/clerk-react";
import { MantineProvider } from "@mantine/core";
import { BrowserRouter, Route, Routes } from "react-router";
import App from "./App.tsx";

// Import your Publishable Key
const PUBLISHABLE_KEY = import.meta.env.VITE_CLERK_PUBLISHABLE_KEY;

if (!PUBLISHABLE_KEY) {
	throw new Error("Add your Clerk Publishable Key to the .env file");
}

const rootElement = document.getElementById("root");

if (!rootElement) {
	console.error("Root element not found. Cannot mount React application.");
} else {
	const root = createRoot(rootElement);
	root.render(
		<StrictMode>
			<BrowserRouter>
				<ClerkProvider publishableKey={PUBLISHABLE_KEY} afterSignOutUrl="/">
					<MantineProvider>
						<Routes>
							<Route index path="/" element={<App />} />
						</Routes>
					</MantineProvider>
				</ClerkProvider>
			</BrowserRouter>
		</StrictMode>,
	);
}
