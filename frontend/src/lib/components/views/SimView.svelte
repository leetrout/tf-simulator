<script lang="ts">
	import GraphPane from '$lib/components/sim/GraphPane.svelte';
	import StatePane from '$lib/components/sim/StatePane.svelte';
</script>

<div class="flex h-full flex-col gap-3 p-3">
	<!-- problem banner -->
	<div
		class="border-base-300 flex items-start justify-between gap-4 rounded-lg border border-l-2 border-l-amber-400 bg-amber-500/5 p-3"
	>
		<div class="flex gap-3">
			<div
				class="flex h-6 w-6 shrink-0 items-center justify-center rounded bg-amber-500/15 text-amber-400"
			>
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
				<div class="text-sm font-semibold">State references a resource that no longer exists</div>
				<div class="text-base-content/60 mt-0.5 text-xs">
					nimbus_storage_bucket.assets is in state but was deleted in the cloud. A plan will try to
					recreate it.
				</div>
			</div>
		</div>
		<div class="shrink-0 text-right">
			<div class="text-base-content/40 text-[10px] tracking-wide uppercase">Suggested fix</div>
			<code class="font-mono text-xs text-amber-400">
				terraform state rm nimbus_storage_bucket.assets
			</code>
		</div>
	</div>

	<!-- two-pane workspace -->
	<div class="grid min-h-0 flex-1 grid-cols-2 gap-3">
		<GraphPane />
		<StatePane />
	</div>
</div>
