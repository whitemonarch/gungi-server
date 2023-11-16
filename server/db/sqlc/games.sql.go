// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.17.2
// source: games.sql

package db

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

const createGame = `-- name: CreateGame :one
INSERT INTO games (current_state, ruleset, type)
VALUES ($1, $2, $3)
RETURNING id
`

type CreateGameParams struct {
	CurrentState string `json:"current_state"`
	Ruleset      string `json:"ruleset"`
	Type         string `json:"type"`
}

func (q *Queries) CreateGame(ctx context.Context, arg CreateGameParams) (uuid.UUID, error) {
	row := q.db.QueryRowContext(ctx, createGame, arg.CurrentState, arg.Ruleset, arg.Type)
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
	_, err := q.db.ExecContext(ctx, createRoom,
		arg.HostID,
		arg.Description,
		arg.Rules,
		arg.Type,
		arg.Color,
	)
	return err
}

const createUndo = `-- name: CreateUndo :one
INSERT INTO undo_request (game_id, for_user, from_user)
VALUES ($1, $2, $3)
RETURNING id
`

type CreateUndoParams struct {
	GameID   uuid.UUID `json:"game_id"`
	ForUser  uuid.UUID `json:"for_user"`
	FromUser uuid.UUID `json:"from_user"`
}

func (q *Queries) CreateUndo(ctx context.Context, arg CreateUndoParams) (int64, error) {
	row := q.db.QueryRowContext(ctx, createUndo, arg.GameID, arg.ForUser, arg.FromUser)
	var id int64
	err := row.Scan(&id)
	return id, err
}

const deleteRoom = `-- name: DeleteRoom :exec
DELETE FROM public.room_list
WHERE id = $1
`

func (q *Queries) DeleteRoom(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, deleteRoom, id)
	return err
}

const gameJunction = `-- name: GameJunction :exec
INSERT INTO player_games (user_id, game_id, color)
VALUES ($1, $2, $3)
`

type GameJunctionParams struct {
	UserID uuid.UUID `json:"user_id"`
	GameID uuid.UUID `json:"game_id"`
	Color  string    `json:"color"`
}

func (q *Queries) GameJunction(ctx context.Context, arg GameJunctionParams) error {
	_, err := q.db.ExecContext(ctx, gameJunction, arg.UserID, arg.GameID, arg.Color)
	return err
}

const getGame = `-- name: GetGame :one
SELECT 
  games.id, games.fen, games.history, games.completed, games.date_started, games.date_finished, games.current_state, games.ruleset, games.type, 
  user1.username AS player1,
  user2.username AS player2
FROM games 
JOIN player_games AS player_games_1 ON games.id = player_games_1.game_id AND player_games_1.color = 'w' 
JOIN player_games AS player_games_2 ON games.id = player_games_2.game_id AND player_games_2.color = 'b' 
JOIN profiles AS user1 ON user1.id = player_games_1.user_id
JOIN profiles AS user2 ON user2.id = player_games_2.user_id
WHERE games.id = $1
`

type GetGameRow struct {
	ID           uuid.UUID      `json:"id"`
	Fen          sql.NullString `json:"fen"`
	History      sql.NullString `json:"history"`
	Completed    bool           `json:"completed"`
	DateStarted  time.Time      `json:"date_started"`
	DateFinished sql.NullTime   `json:"date_finished"`
	CurrentState string         `json:"current_state"`
	Ruleset      string         `json:"ruleset"`
	Type         string         `json:"type"`
	Player1      string         `json:"player1"`
	Player2      string         `json:"player2"`
}

func (q *Queries) GetGame(ctx context.Context, id uuid.UUID) (GetGameRow, error) {
	row := q.db.QueryRowContext(ctx, getGame, id)
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
		&i.Player1,
		&i.Player2,
	)
	return i, err
}

const getIdFromUsername = `-- name: GetIdFromUsername :one
SELECT id FROM profiles
WHERE profiles.username = $1
`

func (q *Queries) GetIdFromUsername(ctx context.Context, username string) (uuid.UUID, error) {
	row := q.db.QueryRowContext(ctx, getIdFromUsername, username)
	var id uuid.UUID
	err := row.Scan(&id)
	return id, err
}

const getOngoingGames = `-- name: GetOngoingGames :many
SELECT games.id, games.fen, games.completed, games.date_started, games.current_state, users1.username as username1, users2.username as username2
FROM games
JOIN player_games j ON games.id = j.game_id
JOIN profiles as users1 ON j.user_id = users1.id AND j.color = 'w'
JOIN player_games j2 ON games.id = j2.game_id AND j2.user_id != j.user_id AND j2.color ='b'
JOIN profiles as users2 ON j2.user_id = users2.id
WHERE (users1.id = $1 OR users2.id = $1) AND games.completed=false
`

type GetOngoingGamesRow struct {
	ID           uuid.UUID      `json:"id"`
	Fen          sql.NullString `json:"fen"`
	Completed    bool           `json:"completed"`
	DateStarted  time.Time      `json:"date_started"`
	CurrentState string         `json:"current_state"`
	Username1    string         `json:"username1"`
	Username2    string         `json:"username2"`
}

func (q *Queries) GetOngoingGames(ctx context.Context, id uuid.UUID) ([]GetOngoingGamesRow, error) {
	rows, err := q.db.QueryContext(ctx, getOngoingGames, id)
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
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getRoomList = `-- name: GetRoomList :many
SELECT
    room_list.id, room_list.host_id, room_list.description, room_list.rules, room_list.type, room_list.color, room_list.created_at,
    profiles.username AS host
FROM
    public.room_list
JOIN
    public.profiles ON room_list.host_id = profiles.id
`

type GetRoomListRow struct {
	ID          uuid.UUID `json:"id"`
	HostID      uuid.UUID `json:"host_id"`
	Description string    `json:"description"`
	Rules       string    `json:"rules"`
	Type        string    `json:"type"`
	Color       string    `json:"color"`
	CreatedAt   time.Time `json:"created_at"`
	Host        string    `json:"host"`
}

func (q *Queries) GetRoomList(ctx context.Context) ([]GetRoomListRow, error) {
	rows, err := q.db.QueryContext(ctx, getRoomList)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetRoomListRow
	for rows.Next() {
		var i GetRoomListRow
		if err := rows.Scan(
			&i.ID,
			&i.HostID,
			&i.Description,
			&i.Rules,
			&i.Type,
			&i.Color,
			&i.CreatedAt,
			&i.Host,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getUndos = `-- name: GetUndos :many
SELECT id, game_id, for_user, from_user FROM undo_request
WHERE game_id = $1
`

func (q *Queries) GetUndos(ctx context.Context, gameID uuid.UUID) ([]UndoRequest, error) {
	rows, err := q.db.QueryContext(ctx, getUndos, gameID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []UndoRequest
	for rows.Next() {
		var i UndoRequest
		if err := rows.Scan(
			&i.ID,
			&i.GameID,
			&i.ForUser,
			&i.FromUser,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
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
	row := q.db.QueryRowContext(ctx, getUsernameFromId, id)
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
	ID           uuid.UUID      `json:"id"`
	CurrentState string         `json:"current_state"`
	History      sql.NullString `json:"history"`
}

func (q *Queries) MakeMove(ctx context.Context, arg MakeMoveParams) error {
	_, err := q.db.ExecContext(ctx, makeMove, arg.ID, arg.CurrentState, arg.History)
	return err
}

const removeUndo = `-- name: RemoveUndo :exec
DELETE FROM undo_request
WHERE id = $1
`

func (q *Queries) RemoveUndo(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, removeUndo, id)
	return err
}
