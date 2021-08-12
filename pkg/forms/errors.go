package forms

// Define a new errors type
// Holds the validation error messages for forms.
// Name of the form field will be used as the key
type errors map[string][]string

// Add() - add eroor messages for a given field
func (e errors) Add(field, message string){
	e[field] = append(e[field], message)
}

// Get - retreive the first error message for a given field
func (e errors) Get(field string) string{
	if len(e[field]) > 0 {
		return e[field][0]
	}
	return ""
}