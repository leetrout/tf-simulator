<script lang="ts">
	import { onMount, tick } from 'svelte';
	import { sim } from '$lib/stores/simulator.svelte';
	import { statusMeta, LEGEND_ORDER } from '$lib/sim/status';
	import type { Status } from '$lib/sim/types';
	// Vendored jsongraph (see $lib/jsongraph/PATCHES.md). Types come from the .d.ts.
	import type { Diagram } from '$lib/jsongraph/diagram.js';

	let canvas: HTMLCanvasElement;
	let diagram: Diagram | null = null;
	let jg: typeof import('$lib/jsongraph/index.js') | null = null;

	const region = $derived(sim.snapshot?.meta.region ?? 'us-west-1');

	// A signature that changes whenever the node set or any status changes — used to
	// avoid rebuilding the canvas on unrelated refreshes.
	const signature = $derived(
		(sim.graph?.nodes ?? [])
			.map((n) => `${n.id}:${n.status}`)
			.sort()
			.join('|') + `#${(sim.graph?.edges ?? []).length}`
	);

	// Map the backend graph into a jsongraph { nodes, edges } source. Each node
	// carries statusColor (read by the patched renderer) and a single "status" row.
	function buildSource() {
		const nodes = (sim.graph?.nodes ?? []).map((n) => ({
			id: n.id,
			label: n.address || n.label,
			status: statusMeta(n.status as Status).label,
			statusColor: statusMeta(n.status as Status).hex
		}));
		const edges = (sim.graph?.edges ?? []).map((e) => ({
			source: e.source,
			target: e.target,
			label: e.label ?? ''
		}));
		return { nodes, edges };
	}

	async function render() {
		if (!jg || !canvas) return;
		const nodes = sim.graph?.nodes ?? [];
		if (nodes.length === 0) {
			diagram = null;
			return;
		}
		const opts = { theme: 'dark' as const, dir: 'TB' as const, graph: { fields: ['status'] } };
		if (!diagram) {
			diagram = jg.createGraph(canvas, buildSource(), opts).diagram;
		} else {
			// Re-layout and swap the model, keeping the camera; the patched renderer
			// re-rasterizes (setModel clears the bitmap cache) so colors update.
			const model = jg.graphJSON(buildSource(), opts.graph);
			jg.layout(model, { dir: 'TB' });
			diagram.setModel(model, { keepCamera: true });
		}
	}

	onMount(() => {
		let ro: ResizeObserver | undefined;
		(async () => {
			jg = await import('$lib/jsongraph/index.js');
			await render();
			ro = new ResizeObserver(() => diagram?.resize());
			ro.observe(canvas);
		})();
		return () => ro?.disconnect();
	});

	// Rebuild/update whenever the graph signature changes.
	$effect(() => {
		void signature;
		tick().then(render);
	});
</script>

<section class="border-base-300 bg-base-200 flex min-h-0 flex-col rounded-lg border">
	<div class="border-base-300 flex items-center justify-between gap-3 border-b px-3 py-2">
		<div class="flex shrink-0 items-baseline gap-2">
			<span class="text-sm font-medium whitespace-nowrap">Resource Graph</span>
			<span class="text-base-content/40 hidden font-mono text-xs whitespace-nowrap lg:inline">
				nimbus · {sim.graph?.nodes?.length ?? 0} resources
			</span>
		</div>
		<div class="text-base-content/60 flex items-center gap-2.5 text-[11px]">
			{#each LEGEND_ORDER as s (s)}
				<span class="inline-flex items-center gap-1 whitespace-nowrap">
					<span class="h-1.5 w-1.5 rounded-full {statusMeta(s).dot}"></span>
					{statusMeta(s).label}
				</span>
			{/each}
		</div>
	</div>

	<div class="bg-base-100/40 relative min-h-0 flex-1 overflow-hidden">
		{#if (sim.graph?.nodes?.length ?? 0) === 0}
			<div class="text-base-content/40 flex h-full items-center justify-center text-sm">
				{sim.loading ? 'loading…' : 'No resources. Apply the example or load a scenario.'}
			</div>
		{/if}
		<div
			class="text-base-content/40 pointer-events-none absolute top-2 left-3 z-10 font-mono text-xs"
		>
			region: {region}
		</div>
		<canvas bind:this={canvas} class="h-full w-full"></canvas>
	</div>
</section>
