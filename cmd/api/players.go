package main

import (
	"errors"
	"net/http"

	"soka.hopertz.me/internal/data"
	"soka.hopertz.me/internal/validator"
)

func (app *application) showPlayerHandler(w http.ResponseWriter, r *http.Request) {

	uuid, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	player, err := app.models.Players.Get(uuid)

	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)

		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"player": player}, nil)

	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *application) listPlayersHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		club         string
		firstName    string
		lastName     string
		position     string
		region       string
		previousClub string
		data.Filters
	}

	qs := r.URL.Query()

	v := validator.New()
	input.club = app.readString(qs, "club", "")
	input.firstName = app.readString(qs, "firstname", "")
	input.lastName = app.readString(qs, "lastname", "")
	input.position = app.readString(qs, "position", "")
	input.region = app.readString(qs, "region", "")
	input.previousClub = app.readString(qs, "prevclub", "")
	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)

	input.Filters.Sort = app.readString(qs, "sort", "first_name")
	input.Filters.SortSafeList = []string{"first_name", "-first_name", "club", "-club", "region", "-region"}

	if data.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	players, metadata, err := app.models.Players.GetAll(input.club, input.firstName, input.lastName, input.position, input.region, input.previousClub, input.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"metadata": metadata, "players": players}, nil)

	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}
