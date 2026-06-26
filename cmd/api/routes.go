package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (app *application) routes() http.Handler {
	r := chi.NewRouter()

	r.Group(func(r chi.Router) {
		r.Post("auth/signup", app.signupHandler)
		r.Post("auth/login", app.loginHandler)
	})

	r.Group(func(r chi.Router) {
		r.Use(app.authRequired)

		r.Route("api/expense", func(r chi.Router) {
			r.Post("/create", app.createExpense)
			// r.Delete("/remove", app.removeExpense)
		})

		r.Route("api/expense/{id}", func(r chi.Router) {
			r.Get("/", app.getExpenseByID)
			r.Put("/", app.updateExpenseByID)
			r.Delete("/", app.removeExpenseByID)
		})
	})

	return r
}
