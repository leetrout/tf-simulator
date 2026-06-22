<script lang="ts">
	import { legend, graphNodes } from '$lib/sim/mock';
</script>

<section class="border-base-300 bg-base-200 flex min-h-0 flex-col rounded-lg border">
	<div class="border-base-300 flex items-center justify-between gap-3 border-b px-3 py-2">
		<div class="flex shrink-0 items-baseline gap-2">
			<span class="text-sm font-medium whitespace-nowrap">Resource Graph</span>
			<span class="text-base-content/40 hidden font-mono text-xs whitespace-nowrap lg:inline">
				nimbus · project-alpha
			</span>
		</div>
		<div class="text-base-content/60 flex items-center gap-2.5 text-[11px]">
			{#each legend as item (item.label)}
				<span class="inline-flex items-center gap-1 whitespace-nowrap">
					<span class="h-1.5 w-1.5 rounded-full {item.dot}"></span>
					{item.label}
				</span>
			{/each}
		</div>
	</div>

	<div class="bg-base-100/40 min-h-0 flex-1 overflow-auto p-4">
		<div class="text-base-content/40 mb-3 font-mono text-xs">region: us-west-1</div>
		<div class="grid grid-cols-2 gap-3">
			{#each graphNodes as node (node.name)}
				<div
					class="border-base-300 bg-base-200 rounded-md border p-2.5 {node.wide
						? 'col-span-2'
						: ''}"
				>
					<div class="flex items-center justify-between gap-2">
						<span class="text-base-content/40 text-[10px] tracking-wide uppercase">
							{node.kind}
						</span>
						<span
							class="rounded border px-1.5 py-0.5 text-[10px] tracking-wide uppercase {node.statusClass}"
						>
							{node.status}
						</span>
					</div>
					<div class="text-base-content/90 mt-1.5 truncate font-mono text-sm">{node.name}</div>
				</div>
			{/each}
		</div>
	</div>
</section>
