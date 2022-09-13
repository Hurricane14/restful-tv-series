package presenter

import (
	"series/domain"
	"series/usecase"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindSeriesByIDPresenter(t *testing.T) {
	t.Parallel()

	type Input struct {
		Series  domain.Series
		Reviews []domain.Review
	}
	type Test struct {
		Description string
		Input       Input
		Want        usecase.FindSeriesByIDOutput
	}
	tests := []Test{
		{
			Description: "Sanity check",
			Input: Input{
				Series: domain.NewSeries(
					domain.SeriesID("1be9775b-8d32-4710-9ce6-7ece88e30f01"),
					"Title",
					"Description",
					20,
					1980,
					1990,
					"Creator",
				),
				Reviews: []domain.Review{
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
			},
			Want: usecase.FindSeriesByIDOutput{
				ID:          "1be9775b-8d32-4710-9ce6-7ece88e30f01",
				Title:       "Title",
				Description: "Description",
				Episodes:    20,
				BeginYear:   1980,
				EndYear:     1990,
				Creator:     "Creator",
				Reviews: []usecase.FindSeriesByIDReview{
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
			Input: Input{
				Series: domain.NewSeries(
					domain.SeriesID("1be9775b-8d32-4710-9ce6-7ece88e30f01"),
					"Title",
					"Description",
					20,
					1980,
					1990,
					"Creator",
				),
				Reviews: nil,
			},
			Want: usecase.FindSeriesByIDOutput{
				ID:          "1be9775b-8d32-4710-9ce6-7ece88e30f01",
				Title:       "Title",
				Description: "Description",
				Episodes:    20,
				BeginYear:   1980,
				EndYear:     1990,
				Creator:     "Creator",
				Reviews:     []usecase.FindSeriesByIDReview{},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.Description, func(t *testing.T) {
			assert := assert.New(t)
			presenter := NewFindSeriesByIDPresenter()
			got := presenter.Output(test.Input.Series, test.Input.Reviews)
			assert.Equal(test.Want, got)
		})
	}
}
