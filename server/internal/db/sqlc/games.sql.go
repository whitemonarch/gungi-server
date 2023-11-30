// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.23.0
// source: games.sql

package db

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

const changeUndo = `-- name: ChangeUndo :one
UPDATE undo_request
SET status = $1
WHERE receiver_id = $2 AND game_id = $3
RETURNING undo_request.sender_id
`

type ChangeUndoParams struct {
	Status     string    `json:"status"`
	ReceiverID uuid.UUID `json:"receiver_id"`
	GameID     uuid.UUID `json:"game_id"`
}

func (q *Queries) ChangeUndo(ctx context.Context, arg ChangeUndoParams) (uuid.UUID, error) {
	row := q.db.QueryRow(ctx, changeUndo, arg.Status, arg.ReceiverID, arg.GameID)
	var sender_id uuid.UUID
	err := row.Scan(&sender_id)
	return sender_id, err
}

const changeUsername = `-- name: ChangeUsername :exec
UPDATE profiles
SET username = $1
WHERE profiles.id = $2
`

type ChangeUsernameParams struct {
	Username string    `json:"username"`
	ID       uuid.UUID `json:"id"`
}

func (q *Queries) ChangeUsername(ctx context.Context, arg ChangeUsernameParams) error {
	_, err := q.db.Exec(ctx, changeUsername, arg.Username, arg.ID)
	return err
}

const createGame = `-- name: CreateGame :one
INSERT INTO games (current_state, ruleset, type, user_1, user_2)
VALUES ($1, $2, $3, $4, $5)
RETURNING id
`

type CreateGameParams struct {
	CurrentState string    `json:"current_state"`
	Ruleset      string    `json:"ruleset"`
	Type         string    `json:"type"`
	User1        uuid.UUID `json:"user_1"`
	User2        uuid.UUID `json:"user_2"`
}

func (q *Queries) CreateGame(ctx context.Context, arg CreateGameParams) (uuid.UUID, error) {
	row := q.db.QueryRow(ctx, createGame,
		arg.CurrentState,
		arg.Ruleset,
		arg.Type,
		arg.User1,
		arg.User2,
	)
	var id uuid.UUID
	err := row.Scan(&id)
	return id, err
}

const createRoom = `-- name: CreateRoom :exec
INSERT INTO public.room_list (host_id, description, rules, type, color)
VALUES ($1, $2, $3, $4, $5)
`

type CreateRoomParams struct {
	HostID      uuid.UUID `json:"host_id"`
	Description string    `json:"description"`
	Rules       string    `json:"rules"`
	Type        string    `json:"type"`
	Color       string    `json:"color"`
}

func (q *Queries) CreateRoom(ctx context.Context, arg CreateRoomParams) error {
	_, err := q.db.Exec(ctx, createRoom,
		arg.HostID,
		arg.Description,
		arg.Rules,
		arg.Type,
		arg.Color,
	)
	return err
}

const createUndo = `-- name: CreateUndo :one
INSERT INTO undo_request (game_id, sender_id, receiver_id)
VALUES (
    $1,
    $2,
    (
        SELECT
            CASE
                WHEN games.user_1 = $2 THEN games.user_2
                ELSE games.user_1
            END
        FROM
            games
        WHERE
            games.id = $1
    )
)
RETURNING receiver_id
`

type CreateUndoParams struct {
	GameID   uuid.UUID `json:"game_id"`
	SenderID uuid.UUID `json:"sender_id"`
}

func (q *Queries) CreateUndo(ctx context.Context, arg CreateUndoParams) (uuid.UUID, error) {
	row := q.db.QueryRow(ctx, createUndo, arg.GameID, arg.SenderID)
	var receiver_id uuid.UUID
	err := row.Scan(&receiver_id)
	return receiver_id, err
}

const deleteRoom = `-- name: DeleteRoom :one
WITH deleted_room AS (
    DELETE FROM public.room_list
    WHERE room_list.id = $1
    AND EXISTS (SELECT FROM profiles WHERE profiles.id = room_list.host_id)
    RETURNING id, host_id, description, rules, type, color, created_at
)
SELECT deleted_room.id, deleted_room.host_id, deleted_room.description, deleted_room.rules, deleted_room.type, deleted_room.color, deleted_room.created_at, profiles.username as host
FROM deleted_room
JOIN profiles ON deleted_room.host_id = profiles.id
`

type DeleteRoomRow struct {
	ID          uuid.UUID          `json:"id"`
	HostID      uuid.UUID          `json:"host_id"`
	Description string             `json:"description"`
	Rules       string             `json:"rules"`
	Type        string             `json:"type"`
	Color       string             `json:"color"`
	CreatedAt   pgtype.Timestamptz `json:"created_at"`
	Host        string             `json:"host"`
}

func (q *Queries) DeleteRoom(ctx context.Context, id uuid.UUID) (DeleteRoomRow, error) {
	row := q.db.QueryRow(ctx, deleteRoom, id)
	var i DeleteRoomRow
	err := row.Scan(
		&i.ID,
		&i.HostID,
		&i.Description,
		&i.Rules,
		&i.Type,
		&i.Color,
		&i.CreatedAt,
		&i.Host,
	)
	return i, err
}

const deleteRoomSafe = `-- name: DeleteRoomSafe :one
WITH deleted_room AS (
    DELETE FROM public.room_list
    WHERE room_list.id = $1 AND room_list.host_id = $2
    AND EXISTS (SELECT FROM profiles WHERE profiles.id = room_list.host_id)
    RETURNING id, host_id, description, rules, type, color, created_at
)
SELECT deleted_room.id, deleted_room.host_id, deleted_room.description, deleted_room.rules, deleted_room.type, deleted_room.color, deleted_room.created_at, profiles.username as host
FROM deleted_room
JOIN profiles ON deleted_room.host_id = profiles.id
`

type DeleteRoomSafeParams struct {
	ID     uuid.UUID `json:"id"`
	HostID uuid.UUID `json:"host_id"`
}

type DeleteRoomSafeRow struct {
	ID          uuid.UUID          `json:"id"`
	HostID      uuid.UUID          `json:"host_id"`
	Description string             `json:"description"`
	Rules       string             `json:"rules"`
	Type        string             `json:"type"`
	Color       string             `json:"color"`
	CreatedAt   pgtype.Timestamptz `json:"created_at"`
	Host        string             `json:"host"`
}

func (q *Queries) DeleteRoomSafe(ctx context.Context, arg DeleteRoomSafeParams) (DeleteRoomSafeRow, error) {
	row := q.db.QueryRow(ctx, deleteRoomSafe, arg.ID, arg.HostID)
	var i DeleteRoomSafeRow
	err := row.Scan(
		&i.ID,
		&i.HostID,
		&i.Description,
		&i.Rules,
		&i.Type,
		&i.Color,
		&i.CreatedAt,
		&i.Host,
	)
	return i, err
}

const getGame = `-- name: GetGame :one
SELECT 
    games.id,
    games.fen,
    games.history,
    games.completed,
    games.date_started,
    games.date_finished,
    games.current_state,
    games.ruleset,
    games.type,
    games.user_1,
    games.user_2,
    user1.username AS player1,
    user2.username AS player2
FROM 
    games
JOIN 
    profiles AS user1 ON user1.id = games.user_1
JOIN 
    profiles AS user2 ON user2.id = games.user_2
WHERE
    games.id = $1
`

type GetGameRow struct {
	ID           uuid.UUID          `json:"id"`
	Fen          pgtype.Text        `json:"fen"`
	History      pgtype.Text        `json:"history"`
	Completed    bool               `json:"completed"`
	DateStarted  pgtype.Timestamptz `json:"date_started"`
	DateFinished pgtype.Timestamptz `json:"date_finished"`
	CurrentState string             `json:"current_state"`
	Ruleset      string             `json:"ruleset"`
	Type         string             `json:"type"`
	User1        uuid.UUID          `json:"user_1"`
	User2        uuid.UUID          `json:"user_2"`
	Player1      string             `json:"player1"`
	Player2      string             `json:"player2"`
}

func (q *Queries) GetGame(ctx context.Context, id uuid.UUID) (GetGameRow, error) {
	row := q.db.QueryRow(ctx, getGame, id)
	var i GetGameRow
	err := row.Scan(
		&i.ID,
		&i.Fen,
		&i.History,
		&i.Completed,
		&i.DateStarted,
		&i.DateFinished,
		&i.CurrentState,
		&i.Ruleset,
		&i.Type,
		&i.User1,
		&i.User2,
		&i.Player1,
		&i.Player2,
	)
	return i, err
}

const getGameWithUndo = `-- name: GetGameWithUndo :one
SELECT 
    games.id,
    games.fen,
    games.history,
    games.completed,
    games.date_started,
    games.date_finished,
    games.current_state,
    games.ruleset,
    games.type,
    games.user_1,
    games.user_2,
    user1.username AS player1,
    user2.username AS player2,
    CASE
        WHEN COUNT(undo_request.game_id) > 0 THEN
            json_agg(
                json_build_object(
                    'sender_username', sender.username,
                    'receiver_username', receiver.username,
                    'status', undo_request.status
                )
            )
        ELSE
            '[]'::json
    END AS undo_requests
FROM 
    games
JOIN 
    profiles AS user1 ON user1.id = games.user_1
JOIN 
    profiles AS user2 ON user2.id = games.user_2
LEFT JOIN
    undo_request ON undo_request.game_id = games.id
LEFT JOIN 
    profiles AS sender ON sender.id = undo_request.sender_id
LEFT JOIN 
    profiles AS receiver ON receiver.id = undo_request.receiver_id
WHERE
    games.id = $1
GROUP BY
    games.id, user1.username, user2.username
`

type GetGameWithUndoRow struct {
	ID           uuid.UUID          `json:"id"`
	Fen          pgtype.Text        `json:"fen"`
	History      pgtype.Text        `json:"history"`
	Completed    bool               `json:"completed"`
	DateStarted  pgtype.Timestamptz `json:"date_started"`
	DateFinished pgtype.Timestamptz `json:"date_finished"`
	CurrentState string             `json:"current_state"`
	Ruleset      string             `json:"ruleset"`
	Type         string             `json:"type"`
	User1        uuid.UUID          `json:"user_1"`
	User2        uuid.UUID          `json:"user_2"`
	Player1      string             `json:"player1"`
	Player2      string             `json:"player2"`
	UndoRequests []byte             `json:"undo_requests"`
}

func (q *Queries) GetGameWithUndo(ctx context.Context, id uuid.UUID) (GetGameWithUndoRow, error) {
	row := q.db.QueryRow(ctx, getGameWithUndo, id)
	var i GetGameWithUndoRow
	err := row.Scan(
		&i.ID,
		&i.Fen,
		&i.History,
		&i.Completed,
		&i.DateStarted,
		&i.DateFinished,
		&i.CurrentState,
		&i.Ruleset,
		&i.Type,
		&i.User1,
		&i.User2,
		&i.Player1,
		&i.Player2,
		&i.UndoRequests,
	)
	return i, err
}

const getIdFromUsername = `-- name: GetIdFromUsername :one
SELECT id FROM profiles
WHERE profiles.username = $1
`

func (q *Queries) GetIdFromUsername(ctx context.Context, username string) (uuid.UUID, error) {
	row := q.db.QueryRow(ctx, getIdFromUsername, username)
	var id uuid.UUID
	err := row.Scan(&id)
	return id, err
}

const getOnboarding = `-- name: GetOnboarding :one
SELECT is_username_onboard_complete FROM profiles
WHERE profiles.id = $1
`

func (q *Queries) GetOnboarding(ctx context.Context, id uuid.UUID) (bool, error) {
	row := q.db.QueryRow(ctx, getOnboarding, id)
	var is_username_onboard_complete bool
	err := row.Scan(&is_username_onboard_complete)
	return is_username_onboard_complete, err
}

const getOngoingGames = `-- name: GetOngoingGames :many
SELECT 
    games.id, 
    games.fen, 
    games.completed, 
    games.date_started, 
    games.current_state, 
    users1.username as username1, 
    users2.username as username2
FROM 
    games
JOIN 
    profiles as users1 ON games.user_1 = users1.id
JOIN 
    profiles as users2 ON games.user_2 = users2.id
WHERE
    (users1.id = $1 OR users2.id = $1) AND games.completed=false
`

type GetOngoingGamesRow struct {
	ID           uuid.UUID          `json:"id"`
	Fen          pgtype.Text        `json:"fen"`
	Completed    bool               `json:"completed"`
	DateStarted  pgtype.Timestamptz `json:"date_started"`
	CurrentState string             `json:"current_state"`
	Username1    string             `json:"username1"`
	Username2    string             `json:"username2"`
}

func (q *Queries) GetOngoingGames(ctx context.Context, id uuid.UUID) ([]GetOngoingGamesRow, error) {
	rows, err := q.db.Query(ctx, getOngoingGames, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetOngoingGamesRow
	for rows.Next() {
		var i GetOngoingGamesRow
		if err := rows.Scan(
			&i.ID,
			&i.Fen,
			&i.Completed,
			&i.DateStarted,
			&i.CurrentState,
			&i.Username1,
			&i.Username2,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getRoomList = `-- name: GetRoomList :many
SELECT
    room_list.description,
    room_list.rules,
    room_list.type,
    room_list.color,
    profiles.username AS host
FROM
    public.room_list
JOIN
    public.profiles ON room_list.host_id = profiles.id
`

type GetRoomListRow struct {
	Description string `json:"description"`
	Rules       string `json:"rules"`
	Type        string `json:"type"`
	Color       string `json:"color"`
	Host        string `json:"host"`
}

func (q *Queries) GetRoomList(ctx context.Context) ([]GetRoomListRow, error) {
	rows, err := q.db.Query(ctx, getRoomList)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetRoomListRow
	for rows.Next() {
		var i GetRoomListRow
		if err := rows.Scan(
			&i.Description,
			&i.Rules,
			&i.Type,
			&i.Color,
			&i.Host,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getUsernameFromId = `-- name: GetUsernameFromId :one
SELECT username FROM profiles
WHERE profiles.id = $1
`

func (q *Queries) GetUsernameFromId(ctx context.Context, id uuid.UUID) (string, error) {
	row := q.db.QueryRow(ctx, getUsernameFromId, id)
	var username string
	err := row.Scan(&username)
	return username, err
}

const makeMove = `-- name: MakeMove :exec
UPDATE games
SET current_state = $2, history = $3
WHERE id = $1
`

type MakeMoveParams struct {
	ID           uuid.UUID   `json:"id"`
	CurrentState string      `json:"current_state"`
	History      pgtype.Text `json:"history"`
}

func (q *Queries) MakeMove(ctx context.Context, arg MakeMoveParams) error {
	_, err := q.db.Exec(ctx, makeMove, arg.ID, arg.CurrentState, arg.History)
	return err
}

const removeUndo = `-- name: RemoveUndo :exec
DELETE FROM undo_request
WHERE sender_id = $1 AND game_id = $2
`

type RemoveUndoParams struct {
	SenderID uuid.UUID `json:"sender_id"`
	GameID   uuid.UUID `json:"game_id"`
}

func (q *Queries) RemoveUndo(ctx context.Context, arg RemoveUndoParams) error {
	_, err := q.db.Exec(ctx, removeUndo, arg.SenderID, arg.GameID)
	return err
}

const updateOnboarding = `-- name: UpdateOnboarding :exec
UPDATE profiles
SET is_username_onboard_complete = true
WHERE profiles.id = $1
`

func (q *Queries) UpdateOnboarding(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.Exec(ctx, updateOnboarding, id)
	return err
}
