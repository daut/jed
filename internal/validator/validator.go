package validator

import "strings"

type Validator struct {
	Errors map[string]string
}

func New() *Validator {
	return &Validator{Errors: make(map[string]string)}
}

// HasErrors returns true if the validator has any errors
func (v *Validator) HasErrors() bool {
	return len(v.Errors) > 0
}

// AddError adds an error message to the validator
func (v *Validator) AddError(key, message string) {
	if _, exists := v.Errors[key]; !exists {
		v.Errors[key] = message
	}
}

// IsNotEmpty checks if a string is not empty or whitespace
func (v *Validator) IsNotEmpty(val string, key, message string) {
	if strings.TrimSpace(val) == "" {
		v.AddError("value", message)
	}
}

// Check adds an error message to the validator if the condition is false
func (v *Validator) Check(ok bool, key, message string) {
	if !ok {
		v.AddError(key, message)
	}
}
