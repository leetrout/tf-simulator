<script lang="ts">
	import { sim } from '$lib/stores/simulator.svelte';
	import type { SettingRow } from '$lib/sim/types';

	const sections = $derived<{ title: string; rows: SettingRow[] }[]>([
		{ title: 'Environment', rows: sim.settings?.environment ?? [] },
		{ title: 'Simulator', rows: sim.settings?.simulator ?? [] }
	]);
</script>

<div class="mx-auto h-full max-w-3xl overflow-auto p-4">
	<div class="mb-4">
		<h1 class="text-lg font-semibold">Settings</h1>
		<p class="text-base-content/50 text-sm">simulator configuration</p>
	</div>

	<div
		class="text-base-content/70 mb-6 flex items-center gap-2.5 rounded-lg border border-sky-500/30 bg-sky-500/5 p-3 text-sm"
	>
		<span
			class="flex h-5 w-5 shrink-0 items-center justify-center rounded-full bg-sky-500/20 text-xs font-bold text-sky-400"
		>
			i
		</span>
		<span>These values are read from the running simulator and are read-only for now.</span>
	</div>

	{#each sections as section (section.title)}
		<div class="mb-6">
			<h2 class="text-base-content/40 mb-2 text-[11px] font-semibold tracking-wider uppercase">
				{section.title}
			</h2>
			<div
				class="divide-base-300 border-base-300 bg-base-200 divide-y overflow-hidden rounded-lg border"
			>
				{#each section.rows as row (row.label)}
					<div class="flex items-center justify-between gap-4 px-4 py-3">
						<div>
							<div class="text-sm">{row.label}</div>
							<div class="text-base-content/50 text-xs">{row.description}</div>
						</div>

						{#if row.kind === 'toggle'}
							<input
								type="checkbox"
								class="toggle toggle-primary toggle-sm"
								checked={row.value === true}
								disabled
							/>
						{:else if row.kind === 'color'}
							<span
								class="border-base-300 inline-flex items-center gap-2 rounded-md border px-2 py-1 font-mono text-xs"
							>
								<span class="h-3 w-3 rounded-sm" style="background-color: {row.value}"></span>
								{row.value}
							</span>
						{:else}
							<span
								class="border-base-300 text-base-content/80 rounded-md border px-2.5 py-1 text-xs {row.mono
									? 'font-mono'
									: ''}"
							>
								{row.value}
							</span>
						{/if}
					</div>
				{/each}
			</div>
		</div>
	{/each}
</div>
