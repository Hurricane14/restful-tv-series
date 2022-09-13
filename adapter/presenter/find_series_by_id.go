package presenter

import (
	"series/domain"
	"series/usecase"
)

type findSeriesByIDPresenter struct{}

func NewFindSeriesByIDPresenter() usecase.FindSeriesByIDPresenter {
	return findSeriesByIDPresenter{}
}

func (findSeriesByIDPresenter) Output(
	series domain.Series,
	reviews []domain.Review,
) usecase.FindSeriesByIDOutput {
	output := usecase.FindSeriesByIDOutput{
		ID:          series.ID().String(),
		Title:       series.Title(),
		Description: series.Description(),
		Episodes:    series.Episodes(),
		BeginYear:   series.BeginYear(),
		EndYear:     series.EndYear(),
		Creator:     series.Creator(),
		Reviews:     make([]usecase.FindSeriesByIDReview, len(reviews)),
	}

	for i, review := range reviews {
		output.Reviews[i] = usecase.FindSeriesByIDReview{
			ID:       review.ID().String(),
			AuthorID: review.AuthorID().String(),
			Text:     review.Text(),
		}
	}

	return output
}
