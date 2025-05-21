import type { ApiErrorResponse } from "@/lib/api/types";

class APIError extends Error {
	constructor(message: string, name: string) {
		super(message);
		this.name = `APIError/${name}`;
	}
}

export class AuthError extends APIError {
	public data: ApiErrorResponse;
	constructor(message: string, name: string, data: ApiErrorResponse) {
		super(message, `AuthError/${name}`);
		this.data = data;
	}
}
