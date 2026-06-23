package validator

import "regexp"

type Validator struct {
	Errors map[string]string
}

var (
	EmailRegex         = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9-]+\\.[a-zA-Z]{2,}$")
	PasswordStrengthRx = regexp.MustCompile(`^(?=.*[a-z])(?=.*[A-Z])(?=.*[0-9])(?=.*[^a-zA-Z0-9]).+$`)
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
	v.Check(Matches(password, PasswordStrengthRx), "password", "must contain at least one uppercase letter, lowercase letter, number, and special character")
}

func (v *Validator) ValidateUser(email, password string) {
	v.ValidateEmail(email)
	v.ValidatePasswordPlaintext(password)
}
