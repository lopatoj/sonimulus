import createClient from 'openapi-fetch';
import type { paths } from './schema';
import { PUBLIC_API_PORT, PUBLIC_API_URL, PUBLIC_API_ROUTE } from '$env/static/public';

const client = createClient<paths>({
	baseUrl: `${PUBLIC_API_URL}:${PUBLIC_API_PORT}${PUBLIC_API_ROUTE}`,
	credentials: 'include'
});

export const auth = {
	login: () => {
		window.location.href = `${PUBLIC_API_URL}:${PUBLIC_API_PORT}${PUBLIC_API_ROUTE}/auth`;
	}
};
