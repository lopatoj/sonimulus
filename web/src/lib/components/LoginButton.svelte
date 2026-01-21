<script lang="ts">
	import Soundcloud from '$lib/components/Soundcloud.svelte';
	import type { HTMLButtonAttributes } from 'svelte/elements';

	const { class: cls, ...rest }: HTMLButtonAttributes = $props();
</script>

<button
	class="flex flex-row items-center gap-1 rounded-sm border bg-white p-2 px-3 font-semibold text-black hover:text-[rgba(18,18,18,0.4)] active:opacity-40 {cls}"
	{...rest}
>
	Login with SoundCloud
	<Soundcloud />
</button>
