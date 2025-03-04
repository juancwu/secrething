import { StrictMode } from "react";
import { createRoot } from "react-dom/client";

import "./index.css";

import App from "./App.tsx";
import { ThemeProvider } from "./components/theme-provider.tsx";

const root = document.getElementById("root");
if (root === null) {
	throw new Error("No root element found");
}
createRoot(root).render(
	<StrictMode>
		<ThemeProvider defaultTheme="dark" storageKey="vite-ui-theme">
			<App />
		</ThemeProvider>
	</StrictMode>,
);
