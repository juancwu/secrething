export async function post(path: string, body: unknown) {
	const url = getApiUrl(path);
	const res = await fetch(url, {
		method: "POST",
		body: JSON.stringify(body),
		headers: {
			"Content-Type": "application/json",
			"X-Application-Name": "secrething",
			"X-Application-ID": "secrething@12352312",
		},
		credentials: "include",
	});
	return res;
}

function getApiUrl(path: string): string {
	//biome-ignore lint/style/noParameterAssign: only reassigning one parameter so its fine
	if (path.startsWith("/")) path = path.substring(1);
	return `${import.meta.env.VITE_BACKEND_BASE_URL}/${path}`;
}
