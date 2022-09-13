package validator

type Validator interface {
	Validate(s any) error
	Messages(err error) []string
}
