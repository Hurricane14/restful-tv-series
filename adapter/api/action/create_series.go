package action

import (
	"encoding/json"
	"net/http"
	"series/adapter/api/response"
	"series/adapter/validator"
	"series/usecase"
)

type CreateSeriesAction struct {
	uc        usecase.CreateSeriesUseCase
	validator validator.Validator
}

func NewCreateSeriesAction(
	uc usecase.CreateSeriesUseCase,
	validator validator.Validator,
) CreateSeriesAction {
	return CreateSeriesAction{
		uc:        uc,
		validator: validator,
	}
}

func (a CreateSeriesAction) Execute(w http.ResponseWriter, r *http.Request) {
	var res response.Response
	defer func() { res.Send(w) }()

	input := usecase.CreateSeriesInput{}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		res = response.NewError(http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	if err := a.validator.Validate(input); err != nil {
		res = response.NewError(http.StatusBadRequest, a.validator.Messages(err)...)
		return
	}

	output, err := a.uc.Execute(r.Context(), input)
	switch {
	case err != nil:
		res = response.NewError(http.StatusInternalServerError)
	default:
		res = response.NewSuccess(http.StatusCreated, output)
	}
}
