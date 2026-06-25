<script lang="ts">
	import { sim } from '$lib/stores/simulator.svelte';

	type View = 'state' | 'cloud';
	let view = $state<View>('state');

	const counts = $derived(sim.snapshot?.counts);

	// Reconstruct a readable terraform.tfstate-shaped document from the parsed
	// snapshot (we surface attributes keyed by terraform address).
	const stateDoc = $derived.by(() => {
		const resources: Record<string, Record<string, string>> = {};
		for (const r of sim.snapshot?.state ?? []) resources[r.address] = r.attributes;
		return { version: 4, terraform_version: sim.snapshot?.meta.terraformVersion ?? '', resources };
	});

	const cloudDoc = $derived.by(() => {
		const resources: Record<string, Record<string, string>> = {};
		for (const r of sim.snapshot?.cloud ?? [])
			resources[`${r.type}.${r.id}`] = { id: r.id, ...r.attributes };
		return { provider: 'nimbus', resources };
	});

	const text = $derived(JSON.stringify(view === 'state' ? stateDoc : cloudDoc, null, 2));
</script>

<section class="border-base-300 bg-base-200 flex min-h-0 flex-col rounded-lg border">
	<div class="border-base-300 flex items-center justify-between border-b px-3 py-2">
		<div class="flex items-center gap-2">
			<button
				class="font-mono text-sm font-medium {view === 'state'
					? 'text-base-content'
					: 'text-base-content/40'}"
				onclick={() => (view = 'state')}
			>
				terraform.tfstate
			</button>
			<span class="text-base-content/20">·</span>
			<button
				class="font-mono text-sm font-medium {view === 'cloud'
					? 'text-base-content'
					: 'text-base-content/40'}"
				onclick={() => (view = 'cloud')}
			>
				nimbus cloud
			</button>
		</div>
		<div class="flex items-center gap-2 text-[11px]">
			<span
				class="inline-flex items-center gap-1.5 rounded border border-emerald-500/30 bg-emerald-500/10 px-1.5 py-0.5 text-emerald-400"
			>
				<span class="h-1.5 w-1.5 rounded-full bg-emerald-400"></span>
				{counts?.tracked ?? 0} tracked
			</span>
			<span
				class="inline-flex items-center gap-1.5 rounded border px-1.5 py-0.5 {(counts?.problems ??
					0) > 0
					? 'border-rose-500/30 bg-rose-500/10 text-rose-400'
					: 'border-base-300 text-base-content/40'}"
			>
				<span
					class="h-1.5 w-1.5 rounded-full {(counts?.problems ?? 0) > 0
						? 'bg-rose-400'
						: 'bg-base-content/30'}"
				></span>
				{counts?.problems ?? 0} problem{(counts?.problems ?? 0) === 1 ? '' : 's'}
			</span>
		</div>
	</div>

	<pre
		class="text-base-content/80 min-h-0 flex-1 overflow-auto p-3 font-mono text-xs leading-relaxed">{text}</pre>
</section>
