<script lang="ts">
	import { onMount } from 'svelte';
	import Header from '$lib/components/layout/Header.svelte';
	import TabBar from '$lib/components/layout/TabBar.svelte';
	import StatusBar from '$lib/components/layout/StatusBar.svelte';
	import SimView from '$lib/components/views/SimView.svelte';
	import ScenariosView from '$lib/components/views/ScenariosView.svelte';
	import SettingsView from '$lib/components/views/SettingsView.svelte';
	import { sim } from '$lib/stores/simulator.svelte';
	import type { TabId } from '$lib/sim/mock';

	let activeTab = $state<TabId>('sim');

	onMount(() => {
		sim.start();
		return () => sim.stop();
	});
</script>

<div class="bg-base-100 text-base-content flex h-screen flex-col overflow-hidden">
	<Header />
	<TabBar active={activeTab} onselect={(id) => (activeTab = id)} />

	<main class="min-h-0 flex-1">
		{#if activeTab === 'sim'}
			<SimView />
		{:else if activeTab === 'scenarios'}
			<ScenariosView />
		{:else if activeTab === 'settings'}
			<SettingsView />
		{/if}
	</main>

	<StatusBar />
</div>
