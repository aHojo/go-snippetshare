package forms

import (
	"fmt"
	"net/url"
	"strings"
	"unicode/utf8"
)

// Form embeds a url.Values object to hold form data and errors
type Form struct {
	url.Values
	Errors errors
}

// Define a new function to initialize a form struct. Takes form data as a param
func New(data url.Values) *Form {
	return &Form{data, errors(map[string][]string{})}
}

// Required checks that specific fields in the form are filled out, present, and not blank
// Adds to the errors field if errors are found.
func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		value := f.Get(field)
		if strings.TrimSpace(value) == "" {
			f.Errors.Add(field, "is required")
		}
	}
}

// MaxLength check that a specific field in the form is less than maximum number of elements
func (f *Form) MaxLength(field string, d int) {
	value := f.Get(field)
	if value == "" {
		return
	}
	if utf8.RuneCountInString(value) > d {
		f.Errors.Add(field, fmt.Sprintf("Field is too long (maximum is %d)", d))
	}
}

// PermittedValues check that a specific field in the form is one of the permitted values
func (f *Form)PermittedValues(field string, opts ...string) {
	value := f.Get(field)
	if value == "" {
		return
	}
	for _, opt := range opts {
		if value == opt {
			return
		}
	}
	f.Errors.Add(field, "This field is invalid")
} 


func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}