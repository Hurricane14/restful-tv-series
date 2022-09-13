package presenter

import (
	"series/domain"
	"series/usecase"
)

type createReviewPresenter struct{}

func NewCreateReviewPresenter() usecase.CreateReviewPresenter {
	return createReviewPresenter{}
}

func (createReviewPresenter) Output(review domain.Review) usecase.CreateReviewOutput {
	return usecase.CreateReviewOutput{
		ID:       review.ID().String(),
		SeriesID: review.SeriesID().String(),
		AuthorID: review.AuthorID().String(),
		Text:     review.Text(),
	}
}
