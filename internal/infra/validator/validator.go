package validator

import (
	playground "github.com/go-playground/validator/v10"
)

type Validator struct {
	validate *playground.Validate
}

func New() *Validator {
	return &Validator{
		validate: playground.New(playground.WithRequiredStructEnabled()),
	}
}

func (vl *Validator) Validate(value any) error {
	return vl.validate.Struct(value)
}
