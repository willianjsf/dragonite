<script lang="ts">
	import { resolve } from '$app/paths';
	import { matrixService } from '$lib';
	import { Compass, Users } from '@lucide/svelte';

	// estado derivado, caso userProfile atualize o displayName atualiza automaticamente
	let displayName = $derived(matrixService.userProfile.displayname);
</script>

<main class="flex h-screen flex-col items-center justify-center gap-8">
	<hgroup>
		<h1 class="text-center h1">
			Bem-vindo, <span class="text-primary-100">{displayName || 'João Silva'}</span>
		</h1>
		<p class="text-center">Escolha uma das opções abaixo para começar</p>
	</hgroup>

	<div class="flex flex-col gap-4 md:flex-row">
		<a class="btn preset-filled-primary-50-950 btn-lg" role="button" href={resolve('/#')}
			><Compass /> Descubra Canais Públicos</a
		>

		<a class="btn preset-filled-primary-50-950 btn-lg" role="button" href={resolve('/#')}>
			<Users /> Crie um Canal em Grupo</a
		>
	</div>
</main>

<style>
</style>
