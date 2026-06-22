<script lang="ts">
	import { stateJson } from '$lib/sim/mock';
</script>

<section class="border-base-300 bg-base-200 flex min-h-0 flex-col rounded-lg border">
	<div class="border-base-300 flex items-center justify-between border-b px-3 py-2">
		<span class="font-mono text-sm font-medium">terraform.tfstate</span>
		<div class="flex items-center gap-2 text-[11px]">
			<span
				class="inline-flex items-center gap-1.5 rounded border border-emerald-500/30 bg-emerald-500/10 px-1.5 py-0.5 text-emerald-400"
			>
				<span class="h-1.5 w-1.5 rounded-full bg-emerald-400"></span>
				7 tracked
			</span>
			<span
				class="inline-flex items-center gap-1.5 rounded border border-rose-500/30 bg-rose-500/10 px-1.5 py-0.5 text-rose-400"
			>
				<span class="h-1.5 w-1.5 rounded-full bg-rose-400"></span>
				1 problem
			</span>
		</div>
	</div>

	<pre
		class="text-base-content/80 min-h-0 flex-1 overflow-auto p-3 font-mono text-xs leading-relaxed">{stateJson}</pre>
</section>
