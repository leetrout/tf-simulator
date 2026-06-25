<script lang="ts">
	import { sim } from '$lib/stores/simulator.svelte';

	const meta = $derived(sim.snapshot?.meta);
	let copied = $state<string | null>(null);

	// Plan/Apply/Refresh are teaching aids: the learner runs terraform in their own
	// shell. The buttons copy the exact command; "Refresh" re-reads the dashboard.
	async function copy(cmd: string) {
		try {
			await navigator.clipboard.writeText(cmd);
			copied = cmd;
			setTimeout(() => (copied === cmd ? (copied = null) : null), 1200);
		} catch {
			copied = null;
		}
	}
</script>

<header class="border-base-300 bg-base-200 flex items-center justify-between border-b px-4 py-2.5">
	<div class="flex items-center gap-2.5">
		<div
			class="bg-primary text-primary-content flex h-7 w-7 items-center justify-center rounded-md"
		>
			<svg viewBox="0 0 24 24" fill="none" class="h-4 w-4" stroke="currentColor" stroke-width="2">
				<path stroke-linecap="round" stroke-linejoin="round" d="M12 3 3 7.5 12 12l9-4.5L12 3Z" />
				<path stroke-linecap="round" stroke-linejoin="round" d="m3 12 9 4.5L21 12" />
				<path stroke-linecap="round" stroke-linejoin="round" d="m3 16.5 9 4.5 9-4.5" />
			</svg>
		</div>
		<div class="flex items-baseline gap-2">
			<span class="text-sm font-semibold tracking-tight">statesim</span>
			<span class="text-base-content/50 text-xs">terraform state simulator</span>
		</div>
	</div>

	<div class="flex items-center gap-2 text-xs">
		<span
			class="border-base-300 text-base-content/70 inline-flex items-center gap-1.5 rounded-md border px-2.5 py-1"
		>
			<span class="bg-primary h-1.5 w-1.5 rounded-full"></span>
			{meta?.provider ?? 'Nimbus Cloud'}
		</span>
		<span class="border-base-300 text-base-content/70 rounded-md border px-2.5 py-1 font-mono">
			{meta?.region ?? 'us-west-1'}
		</span>
		<button
			class="btn btn-ghost btn-sm border-base-300 gap-1.5 font-normal"
			onclick={() => sim.refresh()}
			title="Re-read the dashboard (does not run terraform)"
		>
			<svg
				viewBox="0 0 24 24"
				fill="none"
				class="h-3.5 w-3.5"
				stroke="currentColor"
				stroke-width="2"
			>
				<path stroke-linecap="round" stroke-linejoin="round" d="M4 4v5h5M20 20v-5h-5" />
				<path
					stroke-linecap="round"
					stroke-linejoin="round"
					d="M19 9a7 7 0 0 0-12-3L4 9m16 6a7 7 0 0 1-12 3l-3-3"
				/>
			</svg>
			Refresh
		</button>
		<button
			class="btn btn-ghost btn-sm border-base-300 font-normal"
			onclick={() => copy('terraform plan')}
			title="Copy — run it in your shell"
		>
			{copied === 'terraform plan' ? 'copied!' : 'Plan'}
		</button>
		<button
			class="btn btn-primary btn-sm font-normal"
			onclick={() => copy('terraform apply')}
			title="Copy — run it in your shell"
		>
			{copied === 'terraform apply' ? 'copied!' : 'Apply'}
		</button>
	</div>
</header>
