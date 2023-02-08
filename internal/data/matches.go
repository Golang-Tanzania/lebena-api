package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type Match struct {
	IdMatch     string `json:"id_match"`
	HomeTeam    string `json:"home_team"`
	AwayTeam    string `json:"away_team"`
	StadiumName string `json:"stadium_name"`
	Date        string `json:"date"`
}
type MatchModel struct {
	DB *sql.DB
}

func (m MatchModel) Get(uuid string) (*Match, error) {

	query := `SELECT 
	            id_match, h.club AS "home", a.club AS "away", s.full_name as stadium, date_time
              FROM 
	            match
              JOIN 
	            team AS h ON match.id_home = h.id_team
              JOIN
	            team AS a ON match.id_away=a.id_team
              JOIN 
	            stadium AS s ON match.id_stadium = s.id_stadium
	          WHERE id_match = $1`

	var match Match
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)

	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, uuid).Scan(
		&match.IdMatch,
		&match.HomeTeam,
		&match.AwayTeam,
		&match.StadiumName,
		&match.Date,
	)

	if err != nil {

		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}

	}

	return &match, nil
}

func (m MatchModel) GetAll(team string, date string, stadium_name string, filters Filters) ([]*Match, Metadata, error) {

	query := fmt.Sprintf(`
	          SELECT count(*) OVER(), id_match, h.club AS "home", a.club AS "away", s.full_name as stadium, date_time
	          FROM match
              JOIN 
	            team AS h ON match.id_home = h.id_team
              JOIN
	            team AS a ON match.id_away=a.id_team
              JOIN 
	            stadium AS s ON match.id_stadium = s.id_stadium
			  ORDER BY %s %s
			  LIMIT $1  OFFSET $2`, filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)

	defer cancel()

	args := []interface{}{filters.limit(), filters.offset()}
	rows, err := m.DB.QueryContext(ctx, query, args...)

	if err != nil {
		return nil, Metadata{}, err
	}

	// Importantly ensure that resultset is closed before GETAll() returns
	defer rows.Close()

	totalRecords := 0

	matches := []*Match{}

	for rows.Next() {
		var match Match

		err := rows.Scan(
			&totalRecords,
			&match.IdMatch,
			&match.HomeTeam,
			&match.AwayTeam,
			&match.StadiumName,
			&match.Date,
		)
		if err != nil {
			return nil, Metadata{}, err
		}

		matches = append(matches, &match)

	}

	// Retreive any error when encountered during iteration
	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}
	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return matches, metadata, nil
}
