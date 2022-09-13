package action

import (
	"errors"
	"net/http"
	"series/adapter/api/response"
	"series/domain"
	"series/usecase"
)

type FindReviewsBySeriesAction struct {
	uc usecase.FindReviewsBySeriesUseCase
}

func NewFindReviewsBySeriesAction(uc usecase.FindReviewsBySeriesUseCase) FindReviewsBySeriesAction {
	return FindReviewsBySeriesAction{
		uc: uc,
	}
}

func (a FindReviewsBySeriesAction) Execute(w http.ResponseWriter, r *http.Request) {
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
		res = response.NewError(http.StatusBadRequest, err.Error())
	case err != nil:
		res = response.NewError(http.StatusInternalServerError)
	default:
		res = response.NewSuccess(http.StatusOK, output)
	}
}
