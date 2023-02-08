package data

import (
	"context"
	"database/sql"
	"time"
)

type PlayerGoal struct {
	IdPlayer  string `json:"id_player"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Goals     int16  `json:"goals"`
}

type PlayerGoalModel struct {
	DB *sql.DB
}

type PlayerAssistModel struct {
	DB *sql.DB
}

type PlayerAssist struct {
	IdPlayer  string `json:"id_player"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Assists   int16  `json:"assist"`
}

func (pg PlayerGoalModel) GetAll(filters Filters) ([]*PlayerGoal, Metadata, error) {
	query := `
	        SELECT count(*) OVER(), p.id_player,p.first_name, p.last_name, SUM(goals) as total_goals FROM score
	        JOIN player AS p ON score.id_player = p.id_player
            GROUP BY p.id_player
            HAVING SUM(goals) > 0
            ORDER BY SUM(goals) DESC
            LIMIT $1  OFFSET $2`

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)

	defer cancel()

	args := []interface{}{filters.limit(), filters.offset()}
	rows, err := pg.DB.QueryContext(ctx, query, args...)

	if err != nil {
		return nil, Metadata{}, err
	}

	// Importantly ensure that resultset is closed before GETAll() returns
	defer rows.Close()

	totalRecords := 0

	playerGoals := []*PlayerGoal{}

	for rows.Next() {
		var playerGoal PlayerGoal

		err := rows.Scan(
			&totalRecords,
			&playerGoal.IdPlayer,
			&playerGoal.FirstName,
			&playerGoal.LastName,
			&playerGoal.Goals,
		)

		if err != nil {
			return nil, Metadata{}, err
		}

		playerGoals = append(playerGoals, &playerGoal)

	}

	// Retreive any error when encountered during iteration
	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}
	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return playerGoals, metadata, nil

}

func (pa *PlayerAssistModel) GetAll(filters Filters) ([]*PlayerAssist, Metadata, error) {
	query := `
	SELECT count(*) OVER(), p.id_player,p.first_name, p.last_name, SUM(assists) as total_goals FROM score
	JOIN player AS p ON score.id_player = p.id_player
	GROUP BY p.id_player
	HAVING SUM(assists) > 0
	ORDER BY SUM(assists) DESC
	LIMIT $1  OFFSET $2`

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)

	defer cancel()

	args := []interface{}{filters.limit(), filters.offset()}
	rows, err := pa.DB.QueryContext(ctx, query, args...)

	if err != nil {
		return nil, Metadata{}, err
	}

	// Importantly ensure that resultset is closed before GETAll() returns
	defer rows.Close()

	totalRecords := 0

	playerAssists := []*PlayerAssist{}

	for rows.Next() {
		var playerAssist PlayerAssist

		err := rows.Scan(
			&totalRecords,
			&playerAssist.IdPlayer,
			&playerAssist.FirstName,
			&playerAssist.LastName,
			&playerAssist.Assists,
		)

		if err != nil {
			return nil, Metadata{}, err
		}

		playerAssists = append(playerAssists, &playerAssist)

	}

	// Retreive any error when encountered during iteration
	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}
	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return playerAssists, metadata, nil

}
