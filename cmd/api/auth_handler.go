package main

import (
	"database/sql"
	"errors"
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
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := readJSON(w, r, &input)
	if err != nil {
		app.invalidJSONResponse(w, err)
		return
	}

	user, err := app.db.GetUserByEmail(r.Context(), input.Email)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			app.invalidCredentialsResponse(w, r)
			return
		default:
			app.serverErrorResponse(w, r, err)
			return
		}
	}

	password := auth.Password{Hash: user.PasswordHash}
	matches, err := password.Matches(input.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if !matches {
		app.invalidCredentialsResponse(w, r)
		return
	}

	token, err := auth.GenerateJWT(user.ID, user.Email, app.config.jwt.secret)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = writeJSON(w, http.StatusOK, envelope{
		"authentication_token": token,
	})
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
