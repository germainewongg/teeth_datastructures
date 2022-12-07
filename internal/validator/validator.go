package validator

import (
	"regexp"
	"strings"
)

type Validator struct {
	FieldErrors    map[string]string
	NonFieldErrors []string
}

func (v *Validator) IsBlank(value, field string) {
	if v.FieldErrors == nil {
		v.FieldErrors = make(map[string]string)
	}
	if strings.TrimSpace(value) == "" {
		v.FieldErrors[field] = "This field cannot be blank"
	}
}

func (v *Validator) Valid() bool {

	return len(v.FieldErrors) == 0
}

func (v *Validator) ValidEmail(email string) {
	var EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	if !EmailRX.MatchString(email) {
		v.FieldErrors["email"] = "Invalid email"
	}
}

func (v *Validator) AddNonFieldErrors(error string) {
	v.NonFieldErrors = append(v.NonFieldErrors, error)
}
