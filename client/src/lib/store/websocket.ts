import { readable } from 'svelte/store';
import { topNotification } from './notification';
import { browser, dev } from '$app/environment';

type wsConnStateType = 'connecting' | 'connected' | 'reconnecting' | 'closed' | 'error';

interface msgType {
	type: string;
	payload?: unknown;
}

function createWsStore() {
	if (!browser) return;
	const url = dev ? `ws://${import.meta.env.VITE_API_URL}/ws` : `wss://${import.meta.env.VITE_API_URL}/ws`
	const newSocket = new WebSocket(url);
	const { subscribe } = readable<wsConnStateType>('connecting', (set) => {
		newSocket.addEventListener('open', () => {
			set('connected');
		});
		newSocket.addEventListener('error', () => {
			set('error');
		});
		newSocket.addEventListener('close', () => {
			set('closed');
			topNotification?.set('You are disconnected! Please refresh or try again later.')
		});
		newSocket.addEventListener('message', (event) => {
			try {
				const data = JSON.parse(event.data);
				if (data.type == 'auth') {
					// data.payload == 'success' && set('connected');
				}
			} catch (err) {
				console.error('Error: ', err);
			}
		});
	});

	function send(msg: msgType) {
		newSocket.send(JSON.stringify(msg));
	}

	function authenticate(token: string) {
		const msg = {
			type: 'auth',
			payload: `Bearer ${token}`,
		};
		send(msg);
	}

	function close() {
		newSocket.close();
	}

	function addMsgListener(fn: (event?: MessageEvent) => void) {
		newSocket.addEventListener('message', fn);
		return () => {
			newSocket.removeEventListener('message', fn);
		};
	}

	console.log('create ws store');
	return {
		subscribe,
		send,
		authenticate,
		close,
		addMsgListener,
	};
}

export const ws = createWsStore();
