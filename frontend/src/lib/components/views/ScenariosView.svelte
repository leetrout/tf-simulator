<script lang="ts">
	import { sim } from '$lib/stores/simulator.svelte';
	import { statusMeta } from '$lib/sim/status';
</script>

<div class="h-full overflow-auto p-4">
	<div class="mb-4 flex items-center justify-between">
		<div>
			<h1 class="text-lg font-semibold">Scenario Library</h1>
			<p class="text-base-content/50 text-sm">
				{sim.scenarios.length} state conditions · load one into the simulator
			</p>
		</div>
		<button class="btn btn-outline btn-sm border-base-300 font-normal" onclick={() => sim.reset()}>
			Reset to baseline
		</button>
	</div>

	<div class="grid grid-cols-1 gap-3 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
		{#each sim.scenarios as scenario (scenario.id)}
			{@const meta = statusMeta(scenario.status)}
			<div class="border-base-300 bg-base-200 flex flex-col rounded-lg border p-3">
				<div class="flex items-start justify-between gap-2">
					<div class="flex items-center gap-2">
						<span class="h-2 w-2 rounded-full {meta.dot}"></span>
						<h3 class="text-sm font-semibold">{scenario.title}</h3>
					</div>
					<span
						class="rounded border px-1.5 py-0.5 text-[10px] tracking-wide uppercase {meta.badge}"
					>
						{scenario.badge}
					</span>
				</div>

				<p class="text-base-content/60 mt-2 text-xs leading-relaxed">{scenario.description}</p>

				<div class="text-base-content/40 mt-3 text-[10px] tracking-wide uppercase">Affects</div>
				<div class="text-base-content/80 font-mono text-xs">{scenario.affects}</div>

				<code
					class="bg-base-300/50 text-base-content/70 mt-2 block overflow-x-auto rounded p-2 font-mono text-[11px] whitespace-nowrap"
				>
					{scenario.command}
				</code>

				<div class="flex-1"></div>

				{#if scenario.loaded}
					<button class="btn btn-primary btn-sm mt-3 w-full gap-1.5 font-normal" disabled>
						<span class="bg-primary-content h-1.5 w-1.5 rounded-full"></span>
						Loaded
					</button>
				{:else}
					<button
						class="btn btn-outline btn-sm border-base-300 mt-3 w-full font-normal"
						disabled={sim.busyScenario !== null}
						onclick={() => sim.loadScenario(scenario.id)}
					>
						{sim.busyScenario === scenario.id ? 'Loading…' : 'Load in Sim'}
					</button>
				{/if}
			</div>
		{/each}
	</div>
</div>
