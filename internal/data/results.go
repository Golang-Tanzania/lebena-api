package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type Result struct {
	IdMatch     string             `json:"id_match"`
	HomeTeam    string             `json:"home_team"`
	AwayTeam    string             `json:"away_team"`
	StadiumName string             `json:"stadium_name"`
	Date        string             `json:"date"`
	Score       map[string][]int16 `json:"score"`
}
type ResultModel struct {
	DB *sql.DB
}

func (r ResultModel) Get(uuid string) (*Result, error) {

	query := `
	        SELECT  id_match, h.club AS "home",a.club AS "away", s.full_name as stadium, date_time, home AS " ", away AS " "
	        FROM match
	        JOIN 
	        	team AS h ON match.id_home = h.id_team
	        JOIN 
	        	team AS a ON match.id_away = a.id_team
	        JOIN 
	        	stadium AS s ON match.id_stadium = s.id_stadium
        
	        JOIN result using (id_match)
	
	        WHERE id_match = $1`

	var result Result
	var home int16
	var away int16

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)

	defer cancel()

	err := r.DB.QueryRowContext(ctx, query, uuid).Scan(
		&result.IdMatch,
		&result.HomeTeam,
		&result.AwayTeam,
		&result.StadiumName,
		&result.Date,
		&home,
		&away,
	)

	if err != nil {

		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}

	}
	scores := []int16{home, away}

	ft := map[string][]int16{"ft": scores}
	result.Score = ft
	return &result, nil
}

func (r ResultModel) GetAll(team string, date string, stadium_name string, filters Filters) ([]*Result, Metadata, error) {

	query := fmt.Sprintf(`
	        SELECT count(*) OVER(), id_match, h.club AS "home",a.club AS "away", s.full_name as stadium, date_time, home AS " ", away AS " "
	        FROM match
	        JOIN 
	        	team AS h ON match.id_home = h.id_team
	        JOIN 
	        	team AS a ON match.id_away = a.id_team
	        JOIN 
	        	stadium AS s ON match.id_stadium = s.id_stadium
        
	        JOIN result using (id_match)

			ORDER BY %s %s
			  
			LIMIT $1  OFFSET $2`, filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)

	defer cancel()

	args := []interface{}{filters.limit(), filters.offset()}
	rows, err := r.DB.QueryContext(ctx, query, args...)

	if err != nil {
		return nil, Metadata{}, err
	}

	// Importantly ensure that resultset is closed before GETAll() returns
	defer rows.Close()

	totalRecords := 0

	results := []*Result{}

	for rows.Next() {
		var result Result
		var home int16
		var away int16

		err := rows.Scan(
			&totalRecords,
			&result.IdMatch,
			&result.HomeTeam,
			&result.AwayTeam,
			&result.StadiumName,
			&result.Date,
			&home,
			&away,
		)

		if err != nil {
			return nil, Metadata{}, err
		}

		scores := []int16{home, away}

		ft := map[string][]int16{"ft": scores}
		result.Score = ft

		results = append(results, &result)

	}

	// Retreive any error when encountered during iteration
	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}
	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return results, metadata, nil
}
