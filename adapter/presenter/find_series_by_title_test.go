package presenter

import (
	"series/domain"
	"series/usecase"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindSeriesByTitlePresenter(t *testing.T) {
	t.Parallel()

	type Test struct {
		Description string
		Input       []domain.Series
		Want        usecase.FindSeriesByTitleOutput
	}
	tests := []Test{
		{
			Description: "Sanity check",
			Input: []domain.Series{
				domain.NewSeries(
					domain.SeriesID("1be9775b-8d32-4710-9ce6-7ece88e30f01"),
					"Title",
					"Description",
					20,
					1980,
					1990,
					"Creator",
				),
			},
			Want: usecase.FindSeriesByTitleOutput{
				Series: []usecase.FindSeriesByTitleSeries{
					{
						ID:        "1be9775b-8d32-4710-9ce6-7ece88e30f01",
						Title:     "Title",
						BeginYear: 1980,
						EndYear:   1990,
						Creator:   "Creator",
					},
				},
			},
		},
		{
			Description: "No series means empty slice, not nil",
			Input:       nil,
			Want: usecase.FindSeriesByTitleOutput{
				Series: []usecase.FindSeriesByTitleSeries{},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.Description, func(t *testing.T) {
			assert := assert.New(t)
			presenter := NewFindSeriesByTitlePresenter()
			got := presenter.Output(test.Input)
			assert.Equal(test.Want, got, test.Description)
		})
	}
}
