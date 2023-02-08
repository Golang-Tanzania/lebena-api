package main

import (
	"net/http"

	"soka.hopertz.me/internal/data"
	"soka.hopertz.me/internal/validator"
)

func (app *application) showMatchHandler(w http.ResponseWriter, r *http.Request) {

	uuid, err := app.readIDParam(r)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	match, err := app.models.Matches.Get(uuid)

	if err != nil {
		app.notFoundResponse(w, r)
	}

	app.writeJSON(w, http.StatusOK, envelope{"match": match}, nil)

}

func (app *application) listMatchesHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		team         string
		date         string
		stadium_name string
		data.Filters
	}

	qs := r.URL.Query()

	v := validator.New()
	input.team = app.readString(qs, "team", "")
	input.date = app.readString(qs, "date", "")
	input.stadium_name = app.readString(qs, "stadium_name", "")
	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)

	input.Filters.Sort = app.readString(qs, "sort", "date_time")
	input.Filters.SortSafeList = []string{"date_time", "-date_time", "name", "-name", "stadium_name", "-stadium_name"}

	if data.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	matches, metadata, err := app.models.Matches.GetAll(input.team, input.date, input.stadium_name, input.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"metadata": metadata, "matches": matches}, nil)

	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}
