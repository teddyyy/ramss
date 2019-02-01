package main

import "github.com/go-playground/validator"

// Validator ...
type Validator struct {
	validator *validator.Validate
}

// Validate ...
func (v *Validator) Validate(i interface{}) error {
	return v.validator.Struct(i)
}
