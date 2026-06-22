<script lang="ts">
	import type { Entity, EntityID } from '$lib/resources';
	export let entity: Entity;
	export let remove: (id: EntityID) => void;
</script>

<div class="m-2 rounded-sm border-2 border-dotted border-black bg-slate-200 p-3">
	<div>
		{entity.Name}
		<button type="button" class="button text-xs" on:click={() => remove(entity.ID)}
			>❌ Remove</button
		>
	</div>
	<div class="flex text-sm">
		<div class="flex-1">{entity.TurboEncabulationRate}</div>
		<div class="flex-1">{entity.RefractionRate}</div>
	</div>
	<div class="text-xs">{entity.ID}</div>
</div>
