export const ssr = false;

import { createClient } from 'matrix-js-sdk';
import type { LayoutLoad } from './$types';

export const load: LayoutLoad = async () => {
	const token = localStorage.getItem('matrix_access_token');
	const baseUrl = localStorage.getItem('matrix_homeserver');
	const userId = localStorage.getItem('matrix_user_id');

	if (!token || !baseUrl || !userId) {
		return { mClient: null };
	}

	// Tries to login
	try {
		const client = createClient({ baseUrl, accessToken: token, userId });
		await client.startClient();
		return { mClient: client };
	} catch (error) {
		console.error(error);
		return { mClient: null };
	}
};
