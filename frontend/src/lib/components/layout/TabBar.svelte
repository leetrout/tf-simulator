<script lang="ts">
	import { tabs, type TabId } from '$lib/sim/mock';

	let { active, onselect }: { active: TabId; onselect: (id: TabId) => void } = $props();
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
				{#if tab.badge}
					<span
						class="bg-base-300 text-base-content/60 rounded px-1.5 py-0.5 text-[10px] leading-none"
					>
						{tab.badge}
					</span>
				{/if}
			</button>
		{/each}
	</div>

	<div class="flex items-center gap-3 text-xs">
		<span class="text-base-content/40 text-[10px] tracking-wide uppercase">Loaded</span>
		<span
			class="inline-flex items-center gap-1.5 rounded-md border border-amber-500/30 bg-amber-500/10 px-2 py-1 text-amber-400"
		>
			<span class="h-1.5 w-1.5 rounded-full bg-amber-400"></span>
			Phantom in State
		</span>
		<span
			class="inline-flex items-center gap-1.5 rounded-md border border-emerald-500/30 bg-emerald-500/10 px-2 py-1 text-emerald-400"
		>
			<span class="h-1.5 w-1.5 rounded-full bg-emerald-400"></span>
			Simulation
		</span>
	</div>
</nav>
