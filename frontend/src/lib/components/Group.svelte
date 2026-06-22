<script lang="ts">
	import type { Group } from '$lib/resources';
	export let group: Group;
</script>

<div class="">
	<div>{group.Name}</div>
	<div class="flex text-sm"></div>
	<div class="text-xs">{group.ID}</div>
</div>
