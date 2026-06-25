<script lang="ts">
	import { sim } from '$lib/stores/simulator.svelte';
	import { statusMeta, LEGEND_ORDER } from '$lib/sim/status';
	import type { Status } from '$lib/sim/types';

	const counts = $derived(sim.snapshot?.counts);
	const meta = $derived(sim.snapshot?.meta);

	// The most "interesting" non-sync status present, for the leading dot.
	const headline = $derived.by<{ status: Status; n: number } | null>(() => {
		if (!counts) return null;
		for (const s of LEGEND_ORDER) {
			if (s !== 'in_sync' && counts.byStatus[s] > 0) return { status: s, n: counts.byStatus[s] };
		}
		return null;
	});

	const plural = (n: number, w: string) => `${n} ${w}${n === 1 ? '' : 's'}`;
</script>

<footer
	class="border-base-300 bg-base-200 text-base-content/60 flex items-center justify-between border-t px-4 py-1.5 text-xs"
>
	<div class="flex items-center gap-4">
		{#if headline}
			<span class="inline-flex items-center gap-1.5">
				<span class="h-1.5 w-1.5 rounded-full {statusMeta(headline.status).dot}"></span>
				{statusMeta(headline.status).label} in state
			</span>
		{:else}
			<span class="inline-flex items-center gap-1.5">
				<span class="h-1.5 w-1.5 rounded-full bg-emerald-400"></span>
				all in sync
			</span>
		{/if}
		<span>{plural(counts?.problems ?? 0, 'problem')}</span>
		<span>{plural(counts?.total ?? 0, 'resource')}</span>
		{#if sim.lastEvent}
			<span class="text-base-content/40 font-mono">↯ {sim.lastEvent}</span>
		{/if}
	</div>
	<div class="flex items-center gap-4 font-mono">
		<span>backend: {meta?.backend ?? 'local'}</span>
		<span>terraform {meta?.terraformVersion || '—'}</span>
	</div>
</footer>
