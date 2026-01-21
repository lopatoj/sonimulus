import type { LayoutServerLoad } from './$types';
import { PUBLIC_API_URL, PUBLIC_API_PORT, PUBLIC_API_ROUTE } from '$env/static/public';

export const load: LayoutServerLoad<{ authenticated: boolean }> = async ({ cookies, fetch }) => {
	const sessionId = cookies.get('SESSION_ID');
	if (!sessionId) {
		return { authenticated: false };
	}

	const res = await fetch(`${PUBLIC_API_URL}:${PUBLIC_API_PORT}${PUBLIC_API_ROUTE}/auth/validate`);
	return { authenticated: res.ok };
};
