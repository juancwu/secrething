// ./test-utils/render.tsx
import { render as testingLibraryRender } from "@testing-library/react";
import { MantineProvider } from "@mantine/core";
import { ClerkProvider } from "@clerk/clerk-react";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { BrowserRouter } from "react-router";

const PUBLISHABLE_KEY = import.meta.env.VITE_CLERK_PUBLISHABLE_KEY;

const queryClient = new QueryClient();

export function render(ui: React.ReactNode) {
	return testingLibraryRender(<>{ui}</>, {
		wrapper: ({ children }: { children: React.ReactNode }) => (
			<BrowserRouter>
				<ClerkProvider publishableKey={PUBLISHABLE_KEY} afterSignOutUrl="/">
					<MantineProvider>
						<QueryClientProvider client={queryClient}>
							{children}
						</QueryClientProvider>
					</MantineProvider>
				</ClerkProvider>
			</BrowserRouter>
		),
	});
}
