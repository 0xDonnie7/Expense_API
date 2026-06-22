package main

import (
	"net/http"

	"github.com/0xdonnie7/Expense_API/validator"
)

func (app *application) signupHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string
		Password string
	}

	if err := readJSON(w, r, &input); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	v := validator.New()
	v.Check(input.Email != "", "email", "must not be empty")
}
