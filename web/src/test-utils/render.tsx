import { MantineProvider } from "@mantine/core";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { render as testingLibraryRender } from "@testing-library/react";

const queryClient = new QueryClient();

export function render(ui: React.ReactNode) {
	return testingLibraryRender(ui, {
		wrapper: ({ children }: { children: React.ReactNode }) => (
			<MantineProvider>
				<QueryClientProvider client={queryClient}>
					{children}
				</QueryClientProvider>
			</MantineProvider>
		),
	});
}
