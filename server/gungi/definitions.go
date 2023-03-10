package gungi

const BOARD_SQUARE_NUM = 180

type Enum = int

// Piece Enums
const (
	BPA Enum = iota
	BLG
	BMG
	BGE
	BFO
	BKN
	BAR
	BMU
	BSA
	BCN
	BSP
	BCP
	BMA
	WPA
	WLG
	WMG
	WGE
	WFO
	WKN
	WAR
	WMU
	WSA
	WCN
	WSP
	WCP
	WMA
)

// board.Stockpile[0] = 9
// board.Stockpile[1] = 4
// board.Stockpile[2] = 4
// board.Stockpile[3] = 6
// board.Stockpile[4] = 2
// board.Stockpile[5] = 2
// board.Stockpile[6] = 2
// board.Stockpile[7] = 1
// board.Stockpile[8] = 2
// board.Stockpile[9] = 2
// board.Stockpile[10] = 2
// board.Stockpile[11] = 1
// board.Stockpile[12] = 1
// board.Stockpile[13] = 9
// board.Stockpile[14] = 4
// board.Stockpile[15] = 4
// board.Stockpile[16] = 6
// board.Stockpile[17] = 2
// board.Stockpile[18] = 2
// board.Stockpile[19] = 2
// board.Stockpile[20] = 1
// board.Stockpile[21] = 2
// board.Stockpile[22] = 2
// board.Stockpile[23] = 2
// board.Stockpile[24] = 1
// board.Stockpile[25] = 1

// File Enums
const (
	FILE_A Enum = iota
	FILE_B
	FILE_C
	FILE_D
	FILE_E
	FILE_F
	FILE_G
	FILE_H
	FILE_i
	FILE_NONE
)

// Rank Enums
const (
	RANK_1 Enum = iota
	RANK_2
	RANK_3
	RANK_4
	RANK_5
	RANK_6
	RANK_7
	RANK_8
	RANK_9
	RANK_NONE
)

// Color Enums
const (
	BLACK Enum = iota
	WHITE
	BOTH
)

// Coordinate Enums
const (
	A9 Enum = iota + 37
	A8
	A7
	A6
	A5
	A4
	A3
	A2
	A1
	B9 Enum = iota + 37 + 3
	B8
	B7
	B6
	B5
	B4
	B3
	B2
	B1
	C9 Enum = iota + 37 + 6
	C8
	C7
	C6
	C5
	C4
	C3
	C2
	C1
	D9 Enum = iota + 37 + 9
	D8
	D7
	D6
	D5
	D4
	D3
	D2
	D1
	E9 Enum = iota + 37 + 12
	E8
	E7
	E6
	E5
	E4
	E3
	E2
	E1
	F9 Enum = iota + 37 + 15
	F8
	F7
	F6
	F5
	F4
	F3
	F2
	F1
	G9 Enum = iota + 37 + 18
	G8
	G7
	G6
	G5
	G4
	G3
	G2
	G1
	H9 Enum = iota + 37 + 21
	H8
	H7
	H6
	H5
	H4
	H3
	H2
	H1
	I9 Enum = iota + 37 + 24
	I8
	I7
	I6
	I5
	I4
	I3
	I2
	I1
	NO_SQ
)

// True or False
const (
	FALSE Enum = iota
	TRUE
)

// Move Type
const (
	MOVE int = iota
	STACK
	ATTACK
	PLACE
	READY
)
