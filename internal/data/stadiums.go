package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type Stadium struct {
	IdStadium string `json:"id_stadium"`
	FullName  string `json:"full_name"`
	Location  string `json:"location"`
	Capacity  string `json:"capacity"`
}

type StadiumModel struct {
	DB *sql.DB
}

var capacityNull sql.NullString

func (s StadiumModel) Get(uuid string) (*Stadium, error) {

	query := `SELECT id_stadium,full_name,location,capacity
	          FROM stadium
			  WHERE 
			  id_stadium = $1`

	var stadium Stadium
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)

	defer cancel()

	err := s.DB.QueryRowContext(ctx, query, uuid).Scan(
		&stadium.IdStadium,
		&stadium.FullName,
		&stadium.Location,
		&capacityNull,
	)

	if err != nil {

		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}

	}

	capacity := ""

	if capacityNull.Valid {
		capacity = capacityNull.String
	}

	stadium.Capacity = capacity
	return &stadium, nil
}

func (s StadiumModel) GetAll(full_name string, location string, filters Filters) ([]*Stadium, Metadata, error) {

	query := fmt.Sprintf(`
	          SELECT count(*) OVER(), id_stadium, full_name, location, capacity
	          FROM stadium
			  WHERE (to_tsvector('simple',full_name) @@ plainto_tsquery('simple',$1) OR $1 = '')
              AND (to_tsvector('simple', location) @@ plainto_tsquery('simple', $2) OR $2 = '')
			  ORDER BY %s %s
			  LIMIT $3  OFFSET $4`, filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)

	defer cancel()

	args := []interface{}{full_name, location, filters.limit(), filters.offset()}
	rows, err := s.DB.QueryContext(ctx, query, args...)

	if err != nil {
		return nil, Metadata{}, err
	}

	// Importantly ensure that resultset is closed before GETAll() returns
	defer rows.Close()

	totalRecords := 0

	stadiums := []*Stadium{}

	for rows.Next() {
		var stadium Stadium
		var capacityNull sql.NullString

		err := rows.Scan(
			&totalRecords,
			&stadium.IdStadium,
			&stadium.FullName,
			&stadium.Location,
			&capacityNull,
		)

		if err != nil {
			return nil, Metadata{}, err
		}

		capacity := ""

		if capacityNull.Valid {
			capacity = capacityNull.String
		}

		stadium.Capacity = capacity

		stadiums = append(stadiums, &stadium)

	}

	// Retreive any error when encountered during iteration
	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}
	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return stadiums, metadata, nil
}
