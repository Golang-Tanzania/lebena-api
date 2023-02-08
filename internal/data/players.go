package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type Player struct {
	IdPlayer     string `json:"id_player"`
	IdTeam       string `json:"id_team"`
	Team         string `json:"team"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Kit          int16  `json:"kit"`
	Position     string `json:"position"`
	Region       string `json:"region"`
	PlayerPhoto  string `json:"player_photo"`
	PreviousClub string `json:"previous_club"`
}
type PlayerModel struct {
	DB *sql.DB
}

var regionNull sql.NullString
var playerPhotoNull sql.NullString
var previousClubNull sql.NullString

func (p PlayerModel) Get(uuid string) (*Player, error) {

	query := `SELECT id_player,player.id_team,club,first_name,last_name,kit,position,region,player_photo,previous_club 
	FROM player 
	INNER JOIN team ON player.id_team = team.id_team
	WHERE id_player = $1`

	var player Player
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)

	defer cancel()

	err := p.DB.QueryRowContext(ctx, query, uuid).Scan(
		&player.IdPlayer,
		&player.IdTeam,
		&player.Team,
		&player.FirstName,
		&player.LastName,
		&player.Kit,
		&player.Position,
		&regionNull,
		&playerPhotoNull,
		&previousClubNull,
	)

	if err != nil {

		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}

	}

	region := ""
	playePhoto := ""
	previousClub := ""

	if regionNull.Valid {
		region = regionNull.String
	}

	player.Region = region

	if playerPhotoNull.Valid {
		playePhoto = playerPhotoNull.String
	}

	player.PlayerPhoto = playePhoto

	if previousClubNull.Valid {
		previousClub = previousClubNull.String
	}

	player.PreviousClub = previousClub

	return &player, nil
}

func (p PlayerModel) GetAll(club string, first_name string, last_name string, position string, region string, previous_club string, filters Filters) ([]*Player, Metadata, error) {

	query := fmt.Sprintf(`
	          SELECT count(*) OVER(), id_player,player.id_team,club,first_name,last_name,kit,position,region,player_photo,previous_club 
	          FROM player
			  INNER JOIN team ON player.id_team = team.id_team
			  WHERE (to_tsvector('simple',club) @@ plainto_tsquery('simple',$1) OR $1 = '')
              AND (to_tsvector('simple', first_name) @@ plainto_tsquery('simple', $2) OR $2 = '')
			  AND (to_tsvector('simple', last_name) @@ plainto_tsquery('simple', $3) OR $3 = '')
			  AND (to_tsvector('simple', position) @@ plainto_tsquery('simple', $4) OR $4 = '')
			  AND (to_tsvector('simple', region) @@ plainto_tsquery('simple', $5) OR $5 = '')
			  AND (to_tsvector('simple', previous_club) @@ plainto_tsquery('simple', $6) OR $6 = '')
			  ORDER BY %s %s
			  LIMIT $7  OFFSET $8`, filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)

	defer cancel()

	args := []interface{}{club, first_name, last_name, position, region, previous_club, filters.limit(), filters.offset()}
	rows, err := p.DB.QueryContext(ctx, query, args...)

	if err != nil {
		return nil, Metadata{}, err
	}

	// Importantly ensure that resultset is closed before GETAll() returns
	defer rows.Close()

	totalRecords := 0

	players := []*Player{}

	for rows.Next() {
		var player Player
		var regionNull sql.NullString
		var playerPhotoNull sql.NullString
		var previousClubNull sql.NullString

		err := rows.Scan(
			&totalRecords,
			&player.IdPlayer,
			&player.IdTeam,
			&player.Team,
			&player.FirstName,
			&player.LastName,
			&player.Kit,
			&player.Position,
			&regionNull,
			&playerPhotoNull,
			&previousClubNull,
		)

		if err != nil {
			return nil, Metadata{}, err
		}

		region := ""
		playePhoto := ""
		previousClub := ""

		if regionNull.Valid {
			region = regionNull.String
		}

		player.Region = region

		if playerPhotoNull.Valid {
			playePhoto = playerPhotoNull.String
		}

		player.PlayerPhoto = playePhoto

		if previousClubNull.Valid {
			previousClub = previousClubNull.String
		}

		player.PreviousClub = previousClub

		players = append(players, &player)

	}

	// Retreive any error when encountered during iteration
	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}
	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return players, metadata, nil
}
