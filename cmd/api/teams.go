package main

import (
	"errors"
	"net/http"

	"soka.hopertz.me/internal/data"
	"soka.hopertz.me/internal/validator"
)

func (app *application) showTeamHandler(w http.ResponseWriter, r *http.Request) {

	uuid, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
		//return when you're done processing, to prevent further processing
	}

	team, err := app.models.Teams.Get(uuid)

	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)

		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"team": team}, nil)

	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *application) listTeamsHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Club string
		data.Filters
	}

	qs := r.URL.Query()

	v := validator.New()
	input.Club = app.readString(qs, "club", "")
	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)

	input.Filters.Sort = app.readString(qs, "sort", "club")
	input.Filters.SortSafeList = []string{"club", "-club"}

	if data.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	teams, metadata, err := app.models.Teams.GetAll(input.Club, input.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"metadata": metadata, "teams": teams}, nil)

	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}
