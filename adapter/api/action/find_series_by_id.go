package action

import (
	"errors"
	"net/http"
	"series/adapter/api/response"
	"series/domain"
	"series/usecase"
)

type FindSeriesByIDAction struct {
	uc usecase.FindSeriesByIDUseCase
}

func NewFindSeriesByIDAction(uc usecase.FindSeriesByIDUseCase) FindSeriesByIDAction {
	return FindSeriesByIDAction{
		uc: uc,
	}
}

func (a FindSeriesByIDAction) Execute(w http.ResponseWriter, r *http.Request) {
	var res response.Response
	defer func() { res.Send(w) }()

	seriesID, ok := r.Context().Value(CtxKeySeriesID).(string)
	if !ok || !domain.IsValidUUID(seriesID) {
		res = response.NewError(http.StatusBadRequest, "invalid or missing series id")
		return
	}

	output, err := a.uc.Execute(r.Context(), domain.SeriesID(seriesID))
	switch {
	case errors.Is(err, domain.ErrSeriesNotFound):
		res = response.NewError(http.StatusNotFound, err.Error())
	case err != nil:
		res = response.NewError(http.StatusInternalServerError)
	default:
		res = response.NewSuccess(http.StatusOK, output)
	}
}
