package presenter

import (
	"series/domain"
	"series/usecase"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindReviewsBySeriesPresenter(t *testing.T) {
	t.Parallel()

	type Test struct {
		Description string
		Input       []domain.Review
		Want        usecase.FindReviewsBySeriesOutput
	}
	tests := []Test{
		{
			Description: "Sanity check",
			Input: []domain.Review{
				domain.NewReview(
					domain.ReviewID("1be9775b-8d32-4710-9ce6-7ece88e30f01"),
					domain.SeriesID("1be9775b-8d32-4710-9ce6-7ece88e30f01"),
					domain.AuthorID("1be9775b-8d32-4710-9ce6-7ece88e30f01"),
					"Review text",
				),
				domain.NewReview(
					domain.ReviewID("1be9775b-8d32-4710-9ce6-7ece88e30f01"),
					domain.SeriesID("1be9775b-8d32-4710-9ce6-7ece88e30f01"),
					domain.AuthorID("1be9775b-8d32-4710-9ce6-7ece88e30f01"),
					"Review text",
				),
			},
			Want: usecase.FindReviewsBySeriesOutput{
				Reviews: []usecase.FindReviewsBySeriesReview{
					{
						ID:       "1be9775b-8d32-4710-9ce6-7ece88e30f01",
						AuthorID: "1be9775b-8d32-4710-9ce6-7ece88e30f01",
						Text:     "Review text",
					},
					{
						ID:       "1be9775b-8d32-4710-9ce6-7ece88e30f01",
						AuthorID: "1be9775b-8d32-4710-9ce6-7ece88e30f01",
						Text:     "Review text",
					},
				},
			},
		},
		{
			Description: "No reviews means empty slice, not nil",
			Input:       nil,
			Want: usecase.FindReviewsBySeriesOutput{
				Reviews: []usecase.FindReviewsBySeriesReview{},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.Description, func(t *testing.T) {
			assert := assert.New(t)
			presenter := NewFindReviewsBySeriesPresenter()
			got := presenter.Output(test.Input)
			assert.Equal(test.Want, got)
		})
	}
}
