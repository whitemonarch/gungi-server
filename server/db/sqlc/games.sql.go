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

const getGames = `-- name: GetGames :many
SELECT games.id, games.fen, games.completed, games.date_started, games.current_state, users1.raw_user_meta_data -> 'username' as username1, users2.raw_user_meta_data -> 'username' as username2
FROM games
JOIN player_games j ON games.id = j.game_id
JOIN auth.users users1 ON j.user_id = users1.id
JOIN player_games j2 ON games.id = j2.game_id AND j2.user_id != j.user_id
JOIN auth.users users2 ON j2.user_id = users2.id
WHERE ((users1.id = $1 AND j.color ='w') OR (users2.id = $1 AND j.color ='b')) AND games.completed=false
`

type GetGamesRow struct {
	ID           uuid.UUID      `json:"id"`
	Fen          sql.NullString `json:"fen"`
	Completed    bool           `json:"completed"`
	DateStarted  time.Time      `json:"date_started"`
	CurrentState string         `json:"current_state"`
	Username1    interface{}    `json:"username1"`
	Username2    interface{}    `json:"username2"`
}

func (q *Queries) GetGames(ctx context.Context, id uuid.UUID) ([]GetGamesRow, error) {
	rows, err := q.db.QueryContext(ctx, getGames, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetGamesRow
	for rows.Next() {
		var i GetGamesRow
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
