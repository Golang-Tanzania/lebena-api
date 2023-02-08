package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type Team struct {
	IdTeam    string `json:"id_team"`
	TeamPhoto string `json:"team_photo"`
	IdStadium string `json:"id_stadium"`
	Club      string `json:"club"`
}
type TeamModel struct {
	DB *sql.DB
}

var idStadiumNull sql.NullString

func (t TeamModel) Get(uuid string) (*Team, error) {

	query := `SELECT id_team, team_photo, id_stadium, club 
	          FROM team
			  WHERE 
			  id_team = $1`

	var team Team
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)

	defer cancel()

	err := t.DB.QueryRowContext(ctx, query, uuid).Scan(
		&team.IdTeam,
		&team.TeamPhoto,
		&idStadiumNull,
		&team.Club,
	)

	if err != nil {

		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}

	}

	id_stadium := ""

	if idStadiumNull.Valid {
		id_stadium = idStadiumNull.String
	}

	team.IdStadium = id_stadium
	return &team, nil
}

func (t TeamModel) GetAll(club string, filters Filters) ([]*Team, Metadata, error) {
	const yanga = "yanga"

	if yanga == club {
		club = "young"
	}
	query := fmt.Sprintf(`
	          SELECT count(*) OVER(), id_team, team_photo, id_stadium, club 
	          FROM team
			  WHERE (to_tsvector('simple',club) @@ plainto_tsquery('simple',$1) OR $1 = '')
			  ORDER BY %s %s
			  LIMIT $2  OFFSET $3`, filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)

	defer cancel()

	args := []interface{}{club, filters.limit(), filters.offset()}
	rows, err := t.DB.QueryContext(ctx, query, args...)

	if err != nil {
		return nil, Metadata{}, err
	}

	// Importantly ensure that resultset is closed before GETAll() returns
	defer rows.Close()

	totalRecords := 0

	teams := []*Team{}

	for rows.Next() {
		var team Team
		var someTimesNull sql.NullString

		err := rows.Scan(
			&totalRecords,
			&team.IdTeam,
			&team.TeamPhoto,
			&someTimesNull,
			&team.Club,
		)

		if err != nil {
			return nil, Metadata{}, err
		}

		id_stadium := ""

		if someTimesNull.Valid {
			id_stadium = someTimesNull.String
		}

		team.IdStadium = id_stadium

		teams = append(teams, &team)

	}

	// Retreive any error when encountered during iteration
	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}
	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return teams, metadata, nil
}
