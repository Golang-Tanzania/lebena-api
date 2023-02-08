package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
)

type Models struct {
	Teams         TeamModel
	Stadiums      StadiumModel
	Players       PlayerModel
	Matches       MatchModel
	Results       ResultModel
	PlayerAssists PlayerAssistModel
	PlayerGoals   PlayerGoalModel
	Users         UserModel
	Tokens        TokenModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Teams:         TeamModel{DB: db},
		Stadiums:      StadiumModel{DB: db},
		Players:       PlayerModel{DB: db},
		Matches:       MatchModel{DB: db},
		Results:       ResultModel{DB: db},
		PlayerAssists: PlayerAssistModel{DB: db},
		PlayerGoals:   PlayerGoalModel{DB: db},
		Users:         UserModel{DB: db},
		Tokens:        TokenModel{DB: db},
	}
}
