<script lang="ts">
	import { createForm } from 'felte';
	import { validator } from '@felte/validator-yup';
	import { browser } from '$app/environment';

	import { EntitySchema } from '$lib/resources';
	import type { CreateEntity } from '$lib/resources';

	let baseURL: string;

	// FIXME: don't duplicate
	if (browser) {
		baseURL = `${window.location.protocol}//${window.location.host}`;
	}

	const { form, errors, reset } = createForm({
		extend: validator<CreateEntity>({ schema: EntitySchema }),
		onSubmit: async (values) => {
			const payload = JSON.stringify(values);

			let resp = fetch(`${baseURL}/api/entity`, {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json'
				},
				body: payload
			});

			// TODO: do something with the response
			resp.then((x) => console.log(x));

			if ((await resp).status == 201) {
				reset();
			}
		}
	});
</script>

<div class="mb-8 border p-4">
	<form class="form" use:form on:submit|preventDefault>
		<div class="mb-4">
			<label class="mb-2 block text-sm font-bold text-gray-700" for="Name"> Entity Name </label>
			<input
				class="focus:shadow-outline w-full appearance-none rounded border px-3 py-2 leading-tight text-gray-700 shadow focus:outline-none"
				id="Name"
				type="text"
				name="Name"
				placeholder="Entity Name"
			/>
			{#if $errors.Name}
				<span class="error">{$errors.Name}</span>
			{/if}
		</div>

		<div class="mb-4">
			<label class="mb-2 block text-sm font-bold text-gray-700" for="TurboEncabulationRate">
				Turbo Encabulation Rate
			</label>
			<input
				class="focus:shadow-outline w-full appearance-none rounded border px-3 py-2 leading-tight text-gray-700 shadow focus:outline-none"
				id="TurboEncabulationRate"
				type="number"
				name="TurboEncabulationRate"
				placeholder="1"
			/>
			{#if $errors.TurboEncabulationRate}
				<span class="error">{$errors.TurboEncabulationRate}</span>
			{/if}
		</div>
		<div class="mb-4">
			<label class="mb-2 block text-sm font-bold text-gray-700" for="RefractionRate">
				Refraction Rate
			</label>
			<input
				class="focus:shadow-outline w-full appearance-none rounded border px-3 py-2 leading-tight text-gray-700 shadow focus:outline-none"
				id="RefractionRate"
				type="number"
				name="RefractionRate"
				placeholder="Optional"
			/>
			{#if $errors.RefractionRate}
				<span class="error">{$errors.RefractionRate}</span>
			{/if}
		</div>
		<button
			class="focus:shadow-outline rounded bg-blue-500 px-4 py-2 font-bold text-white hover:bg-blue-700 focus:outline-none"
			type="submit"
		>
			Add Entity
		</button>
	</form>
</div>
