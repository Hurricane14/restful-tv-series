package action

import (
	"encoding/json"
	"errors"
	"net/http"
	"series/adapter/api/response"
	"series/adapter/validator"
	"series/domain"
	"series/usecase"
)

type CreateReviewAction struct {
	uc        usecase.CreateReviewUseCase
	validator validator.Validator
}

func NewCreateReviewAction(
	uc usecase.CreateReviewUseCase,
	validator validator.Validator,
) CreateReviewAction {
	return CreateReviewAction{
		uc:        uc,
		validator: validator,
	}
}

func (a CreateReviewAction) Execute(w http.ResponseWriter, r *http.Request) {
	var res response.Response
	defer func() { res.Send(w) }()

	input := usecase.CreateReviewInput{}
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
	case errors.Is(err, domain.ErrSeriesNotFound):
		res = response.NewError(http.StatusBadRequest, err.Error())
	case err != nil:
		res = response.NewError(http.StatusInternalServerError)
	default:
		res = response.NewSuccess(http.StatusCreated, output)
	}
}
