<script lang="ts">
	import { goto, invalidate } from '$app/navigation';
	import { onMount } from 'svelte';
	import Navbar from '$lib/components/Navbar.svelte';
	import Notifications from '$lib/components/Notifications.svelte';
	import TopNotification from '$lib/components/TopNotification.svelte';
	import { setTopNotificationStore, setNotificationStore } from '$lib/store/notificationStore.svelte';
	import { setWebsocketStore } from '$lib/store/websocketStore.svelte.js';

	// z-index order:
	// 0 - base
	// 1 - piece under
	// 2 - piece
	// 3 - name
	// 4 - navbar
	// 5 - top notification
	// 6 - tooltip
	// 10 - notifications

	let { data, children } = $props();

	let notificationStore = setNotificationStore();
	let topNotificationStore = setTopNotificationStore();
	let websocketStore = setWebsocketStore();

	$effect(() => {
		if (websocketStore.state === 'connected' && data.session && websocketStore.isAuthenticated === false) {
			websocketStore.authenticate(data.session.access_token);
		}
	});

	onMount(() => {
		const {
			data: { subscription },
		} = data.supabase.auth.onAuthStateChange((_, newSession) => {
			if (!newSession) {
				/**
				 * Queue this as a task so the navigation won't prevent the
				 * triggering function from completing
				 */
				setTimeout(() => {
					goto('/', { invalidateAll: true });
				});
			}
			if (newSession?.expires_at !== data.session?.expires_at) {
				invalidate('supabase:auth');
			}
		});

		return () => {
			subscription.unsubscribe;
		};
	});
</script>

<svelte:head>
	<title>White Monarch Server</title>
</svelte:head>
<Notifications />
<TopNotification />
<Navbar session={data.session} />
{@render children()}

<style lang="scss" global>
	@import '../main.scss';
	@import '../normalize.css';
</style>
