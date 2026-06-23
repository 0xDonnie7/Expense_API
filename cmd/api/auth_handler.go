package main

import (
	"net/http"

	"github.com/0xdonnie7/Expense_API/internal/auth"
	"github.com/0xdonnie7/Expense_API/internal/database"
	"github.com/0xdonnie7/Expense_API/validator"
	"github.com/google/uuid"
)

func (app *application) signupHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := readJSON(w, r, &input); err != nil {
		app.invalidJSONResponse(w, err)
		return
	}

	v := validator.New()
	v.ValidateUser(input.Email, input.Password)
	if !v.Valid() {
		app.failedValidationResponse(w, v.Errors)
	}

	password := auth.Password{}

	err := password.HashPassword(input.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	user, err := app.db.CreateUser(r.Context(), database.CreateUserParams{
		ID:           uuid.New(),
		Email:        input.Email,
		PasswordHash: password.Hash,
	})
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = writeJSON(w, http.StatusOK, user)
	if err != nil {
		app.invalidJSONResponse(w, err)
	}

}

func (app *application) loginHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string
		Password string
	}

	err := readJSON(w, r, &input)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}
