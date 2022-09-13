package presenter

import (
	"series/domain"
	"series/usecase"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateReviewPresenterOutput(t *testing.T) {
	t.Parallel()

	type Test struct {
		Description string
		Input       domain.Review
		Want        usecase.CreateReviewOutput
	}
	tests := []Test{
		{
			Description: "Sanity check",
			Input: domain.NewReview(
				domain.ReviewID("1be9775b-8d32-4710-9ce6-7ece88e30f01"),
				domain.SeriesID("1be9775b-8d32-4710-9ce6-7ece88e30f01"),
				domain.AuthorID("1be9775b-8d32-4710-9ce6-7ece88e30f01"),
				"Review text",
			),
			Want: usecase.CreateReviewOutput{
				ID:       "1be9775b-8d32-4710-9ce6-7ece88e30f01",
				SeriesID: "1be9775b-8d32-4710-9ce6-7ece88e30f01",
				AuthorID: "1be9775b-8d32-4710-9ce6-7ece88e30f01",
				Text:     "Review text",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.Description, func(t *testing.T) {
			assert := assert.New(t)
			presenter := NewCreateReviewPresenter()
			got := presenter.Output(test.Input)
			assert.Equal(test.Want, got)
		})
	}
}
