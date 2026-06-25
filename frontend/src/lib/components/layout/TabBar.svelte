<script lang="ts">
	import { tabs, type TabId } from '$lib/sim/mock';
	import { sim } from '$lib/stores/simulator.svelte';
	import { statusMeta } from '$lib/sim/status';

	let { active, onselect }: { active: TabId; onselect: (id: TabId) => void } = $props();

	// Live badges: scenario count and problem count.
	const badgeFor = (id: TabId): string | undefined => {
		if (id === 'scenarios') return sim.scenarios.length ? String(sim.scenarios.length) : undefined;
		if (id === 'sim') {
			const p = sim.snapshot?.counts.problems ?? 0;
			return p > 0 ? String(p) : undefined;
		}
		return undefined;
	};

	const loaded = $derived(sim.scenarios.find((s) => s.loaded));
</script>

<nav class="border-base-300 bg-base-100 flex items-center justify-between border-b px-4">
	<div class="flex items-center">
		{#each tabs as tab (tab.id)}
			<button
				class="hover:text-base-content relative flex items-center gap-1.5 border-b-2 px-3 py-2.5 text-sm transition-colors {active ===
				tab.id
					? 'border-primary text-base-content'
					: 'text-base-content/50 border-transparent'}"
				onclick={() => onselect(tab.id)}
			>
				{tab.label}
				{#if badgeFor(tab.id)}
					<span
						class="bg-base-300 text-base-content/60 rounded px-1.5 py-0.5 text-[10px] leading-none"
					>
						{badgeFor(tab.id)}
					</span>
				{/if}
			</button>
		{/each}
	</div>

	<div class="flex items-center gap-3 text-xs">
		<span class="text-base-content/40 text-[10px] tracking-wide uppercase">Loaded</span>
		{#if loaded}
			<span
				class="inline-flex items-center gap-1.5 rounded-md border px-2 py-1 {statusMeta(
					loaded.status
				).badge}"
			>
				<span class="h-1.5 w-1.5 rounded-full {statusMeta(loaded.status).dot}"></span>
				{loaded.title}
			</span>
		{:else}
			<span class="text-base-content/40 inline-flex items-center gap-1.5 px-2 py-1">
				— none —
			</span>
		{/if}
		<span
			class="inline-flex items-center gap-1.5 rounded-md border px-2 py-1 {sim.connected
				? 'border-emerald-500/30 bg-emerald-500/10 text-emerald-400'
				: 'border-base-300 text-base-content/40'}"
		>
			<span
				class="h-1.5 w-1.5 rounded-full {sim.connected ? 'bg-emerald-400' : 'bg-base-content/30'}"
			></span>
			{sim.connected ? 'Live' : 'Offline'}
		</span>
	</div>
</nav>
