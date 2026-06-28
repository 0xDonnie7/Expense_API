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
)

func (app *application) createExpense(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.AuthenticatedUserID(r)
	if !ok {
		app.invalidAuthenticationTokenResponse(w, r)
		return
	}

	var input struct {
		Amount      string `json:"amount"`
		Category    string `json:"category"`
		Description string `json:"description"`
		Date        string `json:"date"`
	}
	err := readJSON(w, r, &input)
	if err != nil {
		app.invalidJSONResponse(w, err)
		return
	}

	v := validator.New()
	v.ValidateExpense(input.Amount, input.Category, input.Description, input.Date)
	if !v.Valid() {
		app.failedValidationResponse(w, v.Errors)
		return
	}

	date, err := time.Parse("2006-01-02", input.Date)
	if err != nil {
		app.failedValidationResponse(w, map[string]string{"date": "must be a valid date in YYYY-MM-DD format"})
		return
	}

	expense, err := app.db.CreateExpense(r.Context(), database.CreateExpenseParams{
		ID:       uuid.New(),
		UserID:   userID,
		Amount:   input.Amount,
		Category: input.Category,
		Description: sql.NullString{
			String: input.Description,
			Valid:  input.Description != "",
		},
		Date:      date,
		CreatedAt: sql.NullTime{Time: time.Now(), Valid: true},
		UpdatedAt: sql.NullTime{Valid: false},
	})

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	type expenseResponse struct {
		ID          uuid.UUID `json:"id"`
		UserID      uuid.UUID `json:"user_id"`
		Amount      string    `json:"amount"`
		Category    string    `json:"category"`
		Description string    `json:"description"`
		Date        time.Time `json:"date"`
		CreatedAt   time.Time `json:"created_at"`
		UpdatedAt   time.Time `json:"updated_at"`
	}

	err = writeJSON(w, http.StatusCreated, expenseResponse{
		ID:          expense.ID,
		UserID:      expense.UserID,
		Amount:      expense.Amount,
		Category:    expense.Category,
		Description: expense.Description.String,
		Date:        date,
		CreatedAt:   expense.CreatedAt.Time,
	})
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

}

func (app *application) getExpenseByID(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.AuthenticatedUserID(r)
	if !ok {
		app.invalidAuthenticationTokenResponse(w, r)
		return
	}

	expenseID, err := checkIDParam(w, r)
	if err != nil {
		app.invalidParameterResponse(w, "id")
		return
	}

	expense, err := app.db.GetExpenseByID(r.Context(), database.GetExpenseByIDParams{
		ID:     expenseID,
		UserID: userID,
	})
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			app.notFoundResponse(w, "record not found")
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = writeJSON(w, http.StatusOK, expense)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *application) removeExpenseByID(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.AuthenticatedUserID(r)
	if !ok {
		app.invalidAuthenticationTokenResponse(w, r)
		return
	}

	expenseID, err := checkIDParam(w, r)
	if err != nil {
		app.invalidParameterResponse(w, "id")
		return
	}

	_, err = app.db.GetExpenseByID(r.Context(), database.GetExpenseByIDParams{
		ID:     expenseID,
		UserID: userID,
	})
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			app.notFoundResponse(w, "record not found")
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.db.DeleteExpense(r.Context(), database.DeleteExpenseParams{
		ID:     expenseID,
		UserID: userID,
	})
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (app *application) updateExpenseByID(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.AuthenticatedUserID(r)
	if !ok {
		app.invalidAuthenticationTokenResponse(w, r)
		return
	}

	expenseID, err := checkIDParam(w, r)
	if err != nil {
		app.invalidParameterResponse(w, "id")
		return
	}

	var input struct {
		Amount      string `json:"amount"`
		Category    string `json:"category"`
		Description string `json:"description"`
		Date        string `json:"date"`
	}

	err = readJSON(w, r, &input)
	if err != nil {
		app.invalidJSONResponse(w, err)
		return
	}

	v := validator.New()
	v.ValidateExpense(input.Amount, input.Category, input.Description, input.Date)
	if !v.Valid() {
		app.failedValidationResponse(w, v.Errors)
		return
	}

	parseDate, err := time.Parse("2006-01-02", input.Date)
	if err != nil {
		app.failedValidationResponse(w, map[string]string{"date": "must be a valid date in YYYY-MM-DD format"})
		return
	}

	updatedExpense, err := app.db.UpdateExpense(r.Context(), database.UpdateExpenseParams{
		ID:       expenseID,
		UserID:   userID,
		Amount:   input.Amount,
		Category: input.Category,
		Description: sql.NullString{
			String: input.Description,
			Valid:  input.Description != "",
		},
		Date:      parseDate,
		UpdatedAt: sql.NullTime{Time: time.Now(), Valid: true},
	})
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			app.notFoundResponse(w, "record not found")
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = writeJSON(w, http.StatusOK, updatedExpense)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}
