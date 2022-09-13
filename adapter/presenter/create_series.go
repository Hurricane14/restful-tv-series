package presenter

import (
	"series/domain"
	"series/usecase"
)

type createSeriesPresenter struct{}

func NewCreateSeriesPresenter() usecase.CreateSeriesPresenter {
	return createSeriesPresenter{}
}

func (createSeriesPresenter) Output(series domain.Series) usecase.CreateSeriesOutput {
	return usecase.CreateSeriesOutput{
		ID:          series.ID().String(),
		Title:       series.Title(),
		Description: series.Description(),
		Episodes:    series.Episodes(),
		BeginYear:   series.BeginYear(),
		EndYear:     series.EndYear(),
		Creator:     series.Creator(),
	}
}
