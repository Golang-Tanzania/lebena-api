package main

import (
	"net/http"

	"soka.hopertz.me/internal/data"
	"soka.hopertz.me/internal/validator"
)

func (app *application) playerGoals(w http.ResponseWriter, r *http.Request) {
	var input struct {
		firstname string
		lastname  string
		data.Filters
	}

	qs := r.URL.Query()

	v := validator.New()
	input.firstname = app.readString(qs, "firstname", "")
	input.lastname = app.readString(qs, "lastname", "")
	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)

	input.Filters.Sort = app.readString(qs, "sort", "first_name")
	input.Filters.SortSafeList = []string{"first_name", "-first_name", "last_name", "-last_name"}

	if data.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	result, metadata, err := app.models.PlayerGoals.GetAll(input.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"metadata": metadata, "Topscores": result}, nil)

	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *application) playerAssists(w http.ResponseWriter, r *http.Request) {
	var input struct {
		firstname string
		lastname  string
		data.Filters
	}

	qs := r.URL.Query()

	v := validator.New()
	input.firstname = app.readString(qs, "firstname", "")
	input.lastname = app.readString(qs, "lastname", "")
	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)

	input.Filters.Sort = app.readString(qs, "sort", "first_name")
	input.Filters.SortSafeList = []string{"first_name", "-first_name", "last_name", "-last_name"}

	if data.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	result, metadata, err := app.models.PlayerAssists.GetAll(input.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"metadata": metadata, "Topassists": result}, nil)

	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}
