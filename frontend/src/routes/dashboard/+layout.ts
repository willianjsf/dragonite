// import { redirect } from '@sveltejs/kit';
import type { LayoutLoad } from './$types';

export const load: LayoutLoad = async ({ parent }) => {
	const { mClient } = await parent();

	if (!mClient) {
		// throw redirect(302, '/login'); // TODO: descomentar essa linha depois
	}

	return {};
};
