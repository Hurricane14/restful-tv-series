package action

type mockValidator struct{}

func (v mockValidator) Validate(s any) error {
	return nil
}

func (v mockValidator) Messages(err error) []string {
	return nil
}

type errorResponse struct {
	Errors []string `json:"errors"`
}
