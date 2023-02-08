package main

import (
	"errors"
	"net/http"

	"soka.hopertz.me/internal/data"
	"soka.hopertz.me/internal/validator"
)

func (app *application) showStadiumHandler(w http.ResponseWriter, r *http.Request) {

	uuid, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
		//return when you're done processing, to prevent further processing
	}

	stadium, err := app.models.Stadiums.Get(uuid)

	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)

		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"stadium": stadium}, nil)

	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *application) listStadiumsHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		stadiumName string
		location    string
		data.Filters
	}

	qs := r.URL.Query()

	v := validator.New()
	input.stadiumName = app.readString(qs, "name", "")
	input.location = app.readString(qs, "location", "")
	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)

	input.Filters.Sort = app.readString(qs, "sort", "full_name")
	input.Filters.SortSafeList = []string{"full_name", "-full_name"}

	if data.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	teams, metadata, err := app.models.Stadiums.GetAll(input.stadiumName, input.location, input.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"metadata": metadata, "teams": teams}, nil)

	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}
