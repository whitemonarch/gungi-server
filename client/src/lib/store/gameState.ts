import { writable, derived } from 'svelte/store';
import { reverseList } from '$lib/helpers';
import { FenToBoard } from '$lib/utils/utils';

export interface BoardState {
	completed: boolean;
	current_state: string;
	date_finished: { Time: string; Valid: boolean };
	date_started: string;
	fen: { String: string; Valid: boolean };
	result: string;
	history: string;
	id: string;
	moveList: { [key: string]: number[] };
	player1: string;
	player2: string;
	ruleset: string;
	type: string;
	// check: string
}

export function createGameStore(initState: BoardState, username: string | null) {
	const gameState = writable<BoardState>(initState);
	const completed = derived(gameState, (data) => {
		const result = {
			completed: data.completed,
			result: data.result,
		};
		return result
	});
	const player1Name = derived(gameState, (data) => data.player1);
	const player2Name = derived(gameState, (data) => data.player2);
	const userColor = derived(gameState, (data) => {
		if (data.player1 === username) return 'w';
		if (data.player2 === username) return 'b';
		return 'spectator';
	});
	const moveHistory = derived(gameState, (data) => data.history.split(' '));
	const manualFlip = writable(false);
	const isViewReversed = derived([player1Name, player2Name, manualFlip], ([player1Name, player2Name, manualFlip]) => {
		if (username !== player1Name && username !== player2Name) return manualFlip;
		const isUserWhite = username === player1Name;
		return manualFlip === isUserWhite;
	});
	const turnColor = derived(gameState, (data) => data.current_state.split(' ')[2]);
	const isPlayer1Ready = derived(gameState, (data) => data.current_state.split(' ')[3][0] === '1');
	const isPlayer2Ready = derived(gameState, (data) => data.current_state.split(' ')[3][1] === '1');
	const isUserTurn = derived([userColor, turnColor], ([userColor, turnColor]) => userColor === turnColor);
	const player1HandList = derived(gameState, (data) => {
		const hands = data.current_state.split(' ')[1];
		const handString = hands.split('/')[0];

		const hand: number[] = [];
		for (let i = 0; i < handString.length; i++) {
			hand.push(Number(handString[i]));
		}
		return hand;
	});
	const player2HandList = derived(gameState, (data) => {
		const hands = data.current_state.split(' ')[1];
		const handString = hands.split('/')[1];

		const hand: number[] = [];
		for (let i = 0; i < handString.length; i++) {
			hand.push(Number(handString[i]));
		}

		return hand;
	});
	const player1ArmyCount = derived(gameState, (data) => {
		const pieces = data.current_state.split(' ')[0];
		const wPieces = pieces.match(/[A-Z]/g);
		return wPieces?.length ?? 0;
	});
	const player2ArmyCount = derived(gameState, (data) => {
		const pieces = data.current_state.split(' ')[0];
		const bPieces = pieces.match(/[a-z]/g);
		return bPieces?.length ?? 0;
	});
	const player1HandCount = derived(player2HandList, (handList) => {
		const count = handList.reduce((a, b) => a + b);
		return count;
	});
	const player2HandCount = derived(player2HandList, (handList) => {
		const count = handList.reduce((a, b) => a + b);
		return count;
	});
	const moveList = derived(gameState, (data) => {
		return data.moveList;
	});
	const moveListUI = derived([moveList, isViewReversed], ([moveList, isViewReversed]) => {
		const transformedMoveList: { [key: number]: number[] } = {};
		if (isViewReversed) {
			for (const key in moveList) {
				const transformedKey = 80 - parseInt(key, 10);
				const transformedValues = moveList[key].map((value) => 80 - value);
				transformedMoveList[transformedKey] = transformedValues;
			}
		} else {
			return moveList;
		}

		return transformedMoveList;
	});
	const boardState = derived(gameState, (data) => FenToBoard(data.current_state));
	const boardUI = derived([boardState, isViewReversed], ([boardState, isViewReversed]) =>
		isViewReversed ? reverseList(boardState) : boardState
	);

	return {
		gameState,
		completed,
		player1Name,
		player2Name,
		userColor,
		moveHistory,
		manualFlip,
		isViewReversed,
		turnColor,
		isPlayer1Ready,
		isPlayer2Ready,
		isUserTurn,
		player1HandList,
		player2HandList,
		player1ArmyCount,
		player2ArmyCount,
		player1HandCount,
		player2HandCount,
		moveList,
		moveListUI,
		boardState,
		boardUI,
	};
}

// export function createGameContext(initState: BoardState, username: string){
// 	const gameStore = createGameStore(initState, username)
// 	const gameState = createService("gameState", gameStore.gameState)
// 	const userColor = createService("userColor", gameStore.userColor)
// 	const manualFlip = createService("manualFlip", gameStore.manualFlip)
// 	const isViewReversed = createService("isViewReversed", gameStore.isViewReversed)
// 	const turnColor = createService("turnColor", gameStore.turnColor)
// 	const isPlayer1Ready = createService("isPlayer1Ready", gameStore.isPlayer1Ready)
// 	const isPlayer2Ready = createService("isPlayer2Ready", gameStore.isPlayer2Ready)
// 	const isUserTurn = createService("isUserTurn", gameStore.isUserTurn)
// 	const player1ArmyCount = createService("player1ArmyCount", gameStore.player1ArmyCount)
// 	const player2ArmyCount = createService("player2ArmyCount", gameStore.player2ArmyCount)
// 	const player1HandCount = createService("player1HandCount", gameStore.player1HandCount)
// 	const player2HandCount = createService("player2HandCount", gameStore.player2HandCount)
// 	const moveListUI = createService("moveListUI", gameStore.moveListUI)
// 	const boardUI = createService("boardUI", gameStore.boardUI)
// 	return {
// 		gameState,
// 		userColor,
// 		manualFlip,
// 		isViewReversed,
// 		turnColor,
// 		isPlayer1Ready,
// 		isPlayer2Ready,
// 		isUserTurn,
// 		player1ArmyCount,
// 		player2ArmyCount,
// 		player1HandCount,
// 		player2HandCount,
// 		moveListUI,
// 		boardUI
// 	}
// }
