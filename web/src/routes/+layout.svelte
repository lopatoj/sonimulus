<script lang="ts">
	import './layout.css';
	import favicon from '$lib/assets/favicon.svg';
	let { children } = $props();
</script>

<svelte:head><link rel="icon" href={favicon} /></svelte:head>

<div class="relative h-full w-full">
	<div
		class="absolute top-0 left-0 z-10 m-4 box-border flex h-16 w-[calc(100%-2rem)] flex-row items-center rounded-full border-t border-b border-t-neutral-700 border-b-neutral-800 bg-black p-6 shadow-md shadow-black/30"
	>
		<h1 class="text-2xl font-bold">Sonimulus</h1>
	</div>
	{#await import('$lib/components/Plane.svelte').then((module) => module.default)}
		<p>Loading...</p>
	{:then Component}
		<Component />
	{:catch error}
		<p>Something went wrong: {error.message}</p>
	{/await}
</div>
<div class="h-full w-feed overflow-y-scroll">
	{@render children()}
</div>

<style>
</style>
