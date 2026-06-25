<script lang="ts">
	import GraphPane from '$lib/components/sim/GraphPane.svelte';
	import StatePane from '$lib/components/sim/StatePane.svelte';
	import { sim } from '$lib/stores/simulator.svelte';
	import { statusMeta } from '$lib/sim/status';
	import type { ResourceStatus } from '$lib/sim/types';

	const problem = $derived<ResourceStatus | undefined>(
		sim.snapshot?.status.find((s) => s.status !== 'in_sync')
	);
	const loadedScenario = $derived(sim.scenarios.find((s) => s.loaded));

	// Headline + suggested fix for the current problem.
	function suggestedFix(p: ResourceStatus): string {
		if (loadedScenario?.command) return loadedScenario.command;
		const addr = p.address || `${p.type}.${p.name}`;
		switch (p.status) {
			case 'phantom':
				return `terraform state rm ${addr}`;
			case 'drift':
				return 'terraform apply  # reconcile drift';
			case 'untracked':
				return `terraform import ${p.type}.${p.name} ${p.id}`;
			case 'moved':
				return `terraform state mv  # or add a moved block`;
			default:
				return 'terraform plan';
		}
	}

	const headline = $derived.by(() => {
		if (!problem) return null;
		const titles: Record<string, string> = {
			phantom: 'State references a resource that no longer exists',
			drift: 'Cloud resource has drifted from your configuration',
			untracked: 'A cloud resource is not tracked in state',
			moved: 'A resource was renamed and needs a moved block'
		};
		return titles[problem.status] ?? 'State condition detected';
	});
</script>

<div class="flex h-full flex-col gap-3 p-3">
	{#if problem}
		{@const meta = statusMeta(problem.status)}
		<div
			class="border-base-300 flex items-start justify-between gap-4 rounded-lg border border-l-2 bg-amber-500/5 p-3 {meta.text}"
			style="border-left-color: currentColor"
		>
			<div class="flex gap-3">
				<div class="flex h-6 w-6 shrink-0 items-center justify-center rounded {meta.badge}">
					<svg
						viewBox="0 0 24 24"
						fill="none"
						class="h-3.5 w-3.5"
						stroke="currentColor"
						stroke-width="2"
					>
						<path stroke-linecap="round" stroke-linejoin="round" d="M12 9v4m0 4h.01" />
						<path
							stroke-linecap="round"
							stroke-linejoin="round"
							d="M10.3 3.9 1.8 18a2 2 0 0 0 1.7 3h17a2 2 0 0 0 1.7-3L13.7 3.9a2 2 0 0 0-3.4 0Z"
						/>
					</svg>
				</div>
				<div>
					<div class="text-base-content text-sm font-semibold">{headline}</div>
					<div class="text-base-content/60 mt-0.5 text-xs">{problem.detail}</div>
				</div>
			</div>
			<div class="shrink-0 text-right">
				<div class="text-base-content/40 text-[10px] tracking-wide uppercase">Suggested fix</div>
				<code class="font-mono text-xs {meta.text}">{suggestedFix(problem)}</code>
			</div>
		</div>
	{:else}
		<div
			class="border-base-300 flex items-center gap-3 rounded-lg border border-l-2 border-l-emerald-400 bg-emerald-500/5 p-3"
		>
			<span class="h-1.5 w-1.5 rounded-full bg-emerald-400"></span>
			<div class="text-sm font-medium text-emerald-400">In sync</div>
			<div class="text-base-content/60 text-xs">
				All {sim.snapshot?.counts.total ?? 0} resources are tracked and current. terraform plan reports
				no changes.
			</div>
		</div>
	{/if}

	<div class="grid min-h-0 flex-1 grid-cols-2 gap-3">
		<GraphPane />
		<StatePane />
	</div>
</div>
