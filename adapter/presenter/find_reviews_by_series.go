package presenter

import (
	"series/domain"
	"series/usecase"
)

type findReviewsBySeriesPresenter struct{}

func NewFindReviewsBySeriesPresenter() usecase.FindReviewsBySeriesPresenter {
	return findReviewsBySeriesPresenter{}
}

func (findReviewsBySeriesPresenter) Output(
	reviews []domain.Review,
) usecase.FindReviewsBySeriesOutput {
	output := usecase.FindReviewsBySeriesOutput{
		Reviews: make([]usecase.FindReviewsBySeriesReview, len(reviews)),
	}

	for i, review := range reviews {
		output.Reviews[i] = usecase.FindReviewsBySeriesReview{
			ID:       review.ID().String(),
			AuthorID: review.AuthorID().String(),
			Text:     review.Text(),
		}
	}
	return output
}
