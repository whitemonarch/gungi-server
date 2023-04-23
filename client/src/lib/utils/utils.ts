export function FenToBoard(fen: string): number[][] {
	const newBoard: number[][] = new Array(81).fill([]);
	const fields = fen.split(' ');
	const split = fields[0].split('/');
	if (split.length != 9) {
		//err
	}

	split.forEach((row, index) => {
		const split2 = row.split(',');
		let fileIndex = 0;
		split2.forEach((column) => {
			const newStack: number[] = [];
			const decodeStack = column.split('');
			decodeStack.forEach((piece) => {
				const charCode = piece.charCodeAt(0);
				if (charCode >= '1'.charCodeAt(0) && charCode <= '9'.charCodeAt(0)) {
					const skipNumber = Number(piece);
					fileIndex += skipNumber - 1;
				} else {
					newStack.push(DecodePiece(piece));
				}
			});
			if (newStack.length > 0) {
				newBoard[CoordsToIndex(fileIndex, index)] = newStack;
			}
			fileIndex += 1;
		});
	});
	return newBoard;
}

export function GetPieceColor(piece: number): string {
	if (piece < 13) {
		return 'w';
	} else {
		return 'b';
	}
}

export function GetTopStack(stack: number[]): number {
	return stack[stack.length - 1];
}

export function CoordsToIndex(file: number, rank: number): number {
	return file + rank * 9;
}

type DecodePieceEnums = {
	[key: string]: number;
};

export function DecodePiece(encodedPiece: string): number {
	const pieceEnums: DecodePieceEnums = {
		P: 0,
		L: 1,
		S: 2,
		G: 3,
		F: 4,
		K: 5,
		Y: 6,
		B: 7,
		W: 8,
		C: 9,
		N: 10,
		T: 11,
		M: 12,
		p: 13,
		l: 14,
		s: 15,
		g: 16,
		f: 17,
		k: 18,
		y: 19,
		b: 20,
		w: 21,
		c: 22,
		n: 23,
		t: 24,
		m: 25,
	};

	return pieceEnums[encodedPiece];
}

type EncodePieceEnums = {
	[key: string]: string;
};

export function EncodePiece(encodedPiece: number): string {
	const pieceEnums: EncodePieceEnums = {
		0: 'P',
		1: 'L',
		2: 'S',
		3: 'G',
		4: 'F',
		5: 'K',
		6: 'Y',
		7: 'B',
		8: 'W',
		9: 'C',
		10: 'N',
		11: 'T',
		12: 'M',
		13: 'p',
		14: 'l',
		15: 's',
		16: 'g',
		17: 'f',
		18: 'k',
		19: 'y',
		20: 'b',
		21: 'w',
		22: 'c',
		23: 'n',
		24: 't',
		25: 'm',
	};

	return pieceEnums[encodedPiece];
}
