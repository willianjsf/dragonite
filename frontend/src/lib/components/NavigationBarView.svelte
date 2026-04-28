<script>
	import { HouseIcon, MessageSquare, SettingsIcon } from '@lucide/svelte';
	import { Navigation } from '@skeletonlabs/skeleton-svelte';
	const links = [
		{ label: 'Home', href: '/#', icon: HouseIcon },
		{ label: 'Rooms', href: '/#', icon: MessageSquare },
		{ label: 'Settings', href: '/#', icon: SettingsIcon }
	];
</script>

<Navigation layout="bar">
	<Navigation.Menu class="grid grid-cols-3 gap-2">
		{#each links as link (link)}
			{@const Icon = link.icon}
			<Navigation.TriggerAnchor href={link.href}>
				<Icon class="size-5" />
				<Navigation.TriggerText>{link.label}</Navigation.TriggerText>
			</Navigation.TriggerAnchor>
		{/each}
	</Navigation.Menu>
</Navigation>
