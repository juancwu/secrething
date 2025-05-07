import { render } from "@/test-utils/render";
import { screen } from "@testing-library/react";
import { describe, expect, it } from "vitest";
import App from "./App";

describe("App (component)", () => {
	it("should render", () => {
		render(<App />);
		expect(screen.getByText("Secrething")).toBeInTheDocument();
	});
});
