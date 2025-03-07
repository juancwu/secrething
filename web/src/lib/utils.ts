import { type ClassValue, clsx } from "clsx";
import { twMerge } from "tailwind-merge";
import { getToken } from "@/lib/auth";

// Generic API request function with auth token handling
export async function apiRequest<T>(
	url: string,
	options: RequestInit = {},
	includeToken = true,
): Promise<T> {
	const token = getToken();

	const headers: Record<string, any> = {
		"Content-Type": "application/json",
		...options.headers,
	};

	if (includeToken && token) {
		headers["Authorization"] = `Bearer ${token}`;
	}

	const response = await fetch(url, {
		...options,
		headers,
	});

	if (!response.ok) {
		const errorData = await response.json().catch(() => null);
		throw new Error(
			errorData?.message || `Request failed with status ${response.status}`,
		);
	}

	return await response.json();
}

export function cn(...inputs: ClassValue[]) {
	return twMerge(clsx(inputs));
}
