// import { redirect } from '@sveltejs/kit';
import { matrixService } from '$lib';
import type { LayoutLoad } from './$types';

export const load: LayoutLoad = async () => {
	if (!matrixService.isAuthenticated()) {
		// throw redirect(302, '/login'); // TODO: descomentar essa linha depois
	}

	return {};
};
