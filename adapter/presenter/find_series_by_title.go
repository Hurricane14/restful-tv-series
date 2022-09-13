package presenter

import (
	"series/domain"
	"series/usecase"
)

type findSeriesByTitlePresenter struct{}

func NewFindSeriesByTitlePresenter() usecase.FindSeriesByTitlePresenter {
	return findSeriesByTitlePresenter{}
}

func (findSeriesByTitlePresenter) Output(series []domain.Series) usecase.FindSeriesByTitleOutput {
	output := usecase.FindSeriesByTitleOutput{
		Series: make([]usecase.FindSeriesByTitleSeries, len(series)),
	}

	for i, series := range series {
		output.Series[i] = usecase.FindSeriesByTitleSeries{
			ID:        series.ID().String(),
			Title:     series.Title(),
			BeginYear: series.BeginYear(),
			EndYear:   series.EndYear(),
			Creator:   series.Creator(),
		}
	}
	return output
}
