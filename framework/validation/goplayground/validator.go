package goplayground

import playground "github.com/go-playground/validator/v10"

type validator struct {
	validator *playground.Validate
}

func NewValidator() *validator {
	return &validator{
		validator: playground.New(),
	}
}

func (v *validator) Validate(s any) error {
	return v.validator.Struct(s)
}

func (v *validator) Messages(err error) []string {
	if err == nil {
		return nil
	}
	verrs, ok := err.(playground.ValidationErrors)
	if !ok {
		return nil
	}
	msgs := make([]string, len(verrs))
	for i, verr := range verrs {
		msgs[i] = verr.Error()
	}
	return msgs
}
