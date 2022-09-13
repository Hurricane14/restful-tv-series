package action

import (
	"net/http"
	"series/adapter/api/response"
	"series/usecase"
)

type FindSeriesByTitleAction struct {
	uc usecase.FindSeriesByTitleUseCase
}

func NewFindSeriesByTitleAction(uc usecase.FindSeriesByTitleUseCase) FindSeriesByTitleAction {
	return FindSeriesByTitleAction{
		uc: uc,
	}
}

func (a FindSeriesByTitleAction) Execute(w http.ResponseWriter, r *http.Request) {
	var res response.Response
	defer func() { res.Send(w) }()

	query, ok := r.Context().Value(CtxKeyTitleQuery).(string)
	if !ok {
		res = response.NewError(http.StatusBadRequest, "invalid or missing series id")
		return
	}

	output, err := a.uc.Execute(r.Context(), query)
	switch {
	case err != nil:
		res = response.NewError(http.StatusInternalServerError)
	default:
		res = response.NewSuccess(http.StatusOK, output)
	}
}
