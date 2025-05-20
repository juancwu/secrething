class APIError extends Error {
	constructor(message: string, name: string) {
		super(message);
		this.name = `APIError/${name}`;
	}
}

export class AuthError extends APIError {
	constructor(message: string, name: string) {
		super(message, `AuthError/${name}`);
	}
}
