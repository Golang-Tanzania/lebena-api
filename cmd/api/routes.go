package main

import (
	"expvar"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)
	router.HandlerFunc(http.MethodGet, "/", app.home)
	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)

	// teams Routes
	router.HandlerFunc(http.MethodGet, "/v1/teams", app.requireActivatedUser(app.listTeamsHandler))
	router.HandlerFunc(http.MethodGet, "/v1/teams/:uuid", app.requireActivatedUser(app.showTeamHandler))

	// stadiums Routes
	router.HandlerFunc(http.MethodGet, "/v1/stadiums", app.requireActivatedUser(app.listStadiumsHandler))
	router.HandlerFunc(http.MethodGet, "/v1/stadiums/:uuid", app.requireActivatedUser(app.showStadiumHandler))

	//player Routes
	router.HandlerFunc(http.MethodGet, "/v1/players", app.requireActivatedUser(app.listPlayersHandler))
	router.HandlerFunc(http.MethodGet, "/v1/players/:uuid", app.requireActivatedUser(app.showPlayerHandler))

	// Matches Routes
	router.HandlerFunc(http.MethodGet, "/v1/matches", app.requireActivatedUser(app.listMatchesHandler))
	router.HandlerFunc(http.MethodGet, "/v1/matches/:uuid", app.requireActivatedUser(app.showMatchHandler))

	// Results Routes
	router.HandlerFunc(http.MethodGet, "/v1/results/:uuid", app.requireActivatedUser(app.showResultHandler))
	router.HandlerFunc(http.MethodGet, "/v1/results", app.requireActivatedUser(app.listResultsHandler))

	//Stats for Players

	router.HandlerFunc(http.MethodGet, "/v1/topscorers", app.requireActivatedUser(app.playerGoals))
	router.HandlerFunc(http.MethodGet, "/v1/topassists", app.requireActivatedUser(app.playerAssists))

	// User Routes

	router.HandlerFunc(http.MethodPost, "/v1/users", app.registerUserHandler)
	router.HandlerFunc(http.MethodPut, "/v1/users/activated", app.activateUserHandler)
	router.HandlerFunc(http.MethodPost, "/v1/tokens/authentication", app.createAuthenticationTokenHandler)

	// Metrics Routes
	router.Handler(http.MethodGet, "/v1/metrics", expvar.Handler())

	return app.metrics(app.recoverPanic(app.enableCORS(app.rateLimit(app.authenticate(router)))))

}
