<script lang="ts">
	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import { matrixService } from '$lib';
	import { toaster } from '$lib/stores/toaster';

	let homeserver = $state('http://localhost:8080');
	let username = $state('');
	let userpassword = $state('');

	async function handleLogin(event: Event) {
		event.preventDefault();
		try {
			// atempts to login
			matrixService.login(homeserver, username, userpassword);
			// redirect to dashboard if login successful
			await goto(resolve('/dashboard'));
			// eslint-disable-next-line @typescript-eslint/no-explicit-any
		} catch (err: any) {
			console.error(err);
			toaster.error({ title: 'Failed to Login', description: err.message });
		}
	}
</script>

<form onsubmit={handleLogin} class="mx-auto w-full max-w-sm space-y-4 rounded-xl border-2 p-6">
	<fieldset class="space-y-4">
		<legend class="text-2xl font-bold">Sign In</legend>

		<label class="label">
			<span class="label-text">Host Server</span>
			<input
				name="host"
				type="text"
				class="input"
				placeholder="matrix.org"
				bind:value={homeserver}
			/>
		</label>

		<label class="label">
			<span class="label-text">Username</span>
			<input
				name="username"
				type="text"
				class="input"
				bind:value={username}
				placeholder="ex.: joao67"
			/>
		</label>

		<label class="label">
			<span class="label-text">Password</span>
			<input
				name="password"
				class="input"
				type="password"
				bind:value={userpassword}
				placeholder="Password"
			/>
		</label>
	</fieldset>
	<fieldset class="flex justify-end">
		<button type="submit" class="btn preset-filled-primary-500">Log in</button>
	</fieldset>
</form>

<style>
	legend {
		margin-left: calc(50% - 35px - 8px);
	}
</style>
