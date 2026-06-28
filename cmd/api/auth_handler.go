package main

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/0xdonnie7/Expense_API/internal/auth"
	"github.com/0xdonnie7/Expense_API/internal/database"
	"github.com/0xdonnie7/Expense_API/validator"
	"github.com/google/uuid"
	"github.com/lib/pq"
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
		return
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
		var pqErr *pq.Error
		switch {
		case errors.As(err, &pqErr) && pqErr.Code == "23505":
			app.badRequestResponse(w, "a user with this email already exists")
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	type userResponse struct {
		ID        string `json:"id"`
		Email     string `json:"email"`
		CreatedAt string `json:"created_at"`
	}

	err = writeJSON(w, http.StatusCreated, userResponse{
		ID:        user.ID.String(),
		Email:     user.Email,
		CreatedAt: user.CreatedAt.Time.Format(time.RFC3339),
	})
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
			app.invalidCredentialsResponse(w, r, "email not found")
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
		app.invalidCredentialsResponse(w, r, "password does not match")
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
