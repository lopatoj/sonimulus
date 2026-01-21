<script lang="ts">
	import { auth } from '$lib/client';
	import LoginButton from '$lib/components/LoginButton.svelte';
</script>

<div class="flex h-full w-full flex-col items-center justify-center">
	<p class="mb-4 w-2/3 text-center text-lg">
		Welcome to Sonimulus! In order to get started, please connect your SoundCloud account.
	</p>
	<LoginButton onclick={auth.login} class="hover:cursor-pointer" />
</div>
