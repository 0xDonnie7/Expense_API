package validator

import (
	"regexp"
	"slices"
	"time"
	"unicode"
)

type Validator struct {
	Errors map[string]string
}

var (
	EmailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9-]+\\.[a-zA-Z]{2,}$")
	MoneyRx    = regexp.MustCompile(`^(0\.(0[1-9]|[1-9]\d)|[1-9]\d{0,7}(\.\d{1,2})?)$`)
)

func New() *Validator {
	return &Validator{
		Errors: make(map[string]string),
	}
}

func (v *Validator) AddError(key, message string) {
	if _, exists := v.Errors[key]; !exists {
		v.Errors[key] = message
	}
}

func (v *Validator) Check(ok bool, key, message string) {
	if !ok {
		v.AddError(key, message)
	}
}

func (v *Validator) Error() map[string]string {
	return v.Errors
}

func (v *Validator) Valid() bool {
	return len(v.Errors) == 0
}

func Matches(values string, rx *regexp.Regexp) bool {
	if rx == nil {
		return false
	}
	return rx.MatchString(values)
}

func (v *Validator) ValidateEmail(email string) {
	v.Check(email != "", "email", "must be provided")
	v.Check(Matches(email, EmailRegex), "email", "invalid email address")
}

func (v *Validator) ValidatePasswordPlaintext(password string) {
	v.Check(password != "", "password", "must be provided")
	v.Check(len(password) >= 8, "password", "must be atleast 8 characters")
	v.Check(len(password) <= 72, "password", "must not be more than 72 bytes long")
	v.Check(StrongPassword(password), "password", "must contain at least one uppercase letter, lowercase letter, number, and special character")
}

func (v *Validator) ValidateUser(email, password string) {
	v.ValidateEmail(email)
	v.ValidatePasswordPlaintext(password)
}

func (v *Validator) ValidateExpense(amount, category, description, date string) {
	v.Check(amount != "", "amount", "must be provided")
	v.Check(Matches(amount, MoneyRx), "amount", "must be a positive amount with up to 2 decimal places")

	v.Check(category != "", "category", "must be provided")
	v.Check(PermittedValue(category, "Groceries", "Leisure", "Electronics", "Utilities", "Clothing", "Health", "Others"), "category", "must be a valid category")

	v.Check(len(description) <= 1000, "description", "must not be more than 1000 characters long")

	v.Check(date != "", "date", "must be provided")
	_, err := time.Parse("2006-01-02", date)
	v.Check(err == nil, "date", "must be a valid date in YYYY-MM-DD format")
}

func PermittedValue[T comparable](value T, permittedValues ...T) bool {
	return slices.Contains(permittedValues, value)
}

func StrongPassword(password string) bool {
	if len(password) < 8 {
		return false
	}

	var hasLower, hasUpper, hasDigit, hasSymbol bool

	for _, r := range password {
		switch {
		case unicode.IsLower(r):
			hasLower = true
		case unicode.IsUpper(r):
			hasUpper = true
		case unicode.IsDigit(r):
			hasDigit = true
		case unicode.IsPunct(r) || unicode.IsSymbol(r):
			hasSymbol = true
		}
	}

	return hasLower && hasUpper && hasDigit && hasSymbol
}
