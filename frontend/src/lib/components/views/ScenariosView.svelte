<script lang="ts">
	import { scenarios } from '$lib/sim/mock';
</script>

<div class="h-full overflow-auto p-4">
	<div class="mb-4">
		<h1 class="text-lg font-semibold">Scenario Library</h1>
		<p class="text-base-content/50 text-sm">5 state conditions · load one into the simulator</p>
	</div>

	<div class="grid grid-cols-1 gap-3 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
		{#each scenarios as scenario (scenario.title)}
			<div class="border-base-300 bg-base-200 flex flex-col rounded-lg border p-3">
				<div class="flex items-start justify-between gap-2">
					<div class="flex items-center gap-2">
						<span class="h-2 w-2 rounded-full {scenario.dot}"></span>
						<h3 class="text-sm font-semibold">{scenario.title}</h3>
					</div>
					<span
						class="rounded border px-1.5 py-0.5 text-[10px] tracking-wide uppercase {scenario.badgeClass}"
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
					<button class="btn btn-primary btn-sm mt-3 w-full gap-1.5 font-normal">
						<span class="bg-primary-content h-1.5 w-1.5 rounded-full"></span>
						Loaded
					</button>
				{:else}
					<button class="btn btn-outline btn-sm border-base-300 mt-3 w-full font-normal">
						Load in Sim
					</button>
				{/if}
			</div>
		{/each}
	</div>
</div>
