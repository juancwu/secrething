import { defineWorkspace } from "vitest/config";

export default defineWorkspace([
	{
		extends: "vite.config.ts",
		test: {
			browser: {
				enabled: true,
				headless: true,
				provider: "playwright",
				instances: [
					{ browser: "chromium" },
					{ browser: "firefox" },
					{ browser: "webkit" },
				],
			},
		},
	},
]);
