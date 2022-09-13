package presenter

import (
	"series/domain"
	"series/usecase"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateSeriesPresenterOutput(t *testing.T) {
	t.Parallel()

	type Test struct {
		Description string
		Input       domain.Series
		Want        usecase.CreateSeriesOutput
	}
	tests := []Test{
		{
			Description: "Sanity check",
			Input: domain.NewSeries(
				domain.SeriesID("1be9775b-8d32-4710-9ce6-7ece88e30f01"),
				"Title",
				"Description",
				20,
				1980,
				1990,
				"Creator",
			),
			Want: usecase.CreateSeriesOutput{
				ID:          "1be9775b-8d32-4710-9ce6-7ece88e30f01",
				Title:       "Title",
				Description: "Description",
				Episodes:    20,
				BeginYear:   1980,
				EndYear:     1990,
				Creator:     "Creator",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.Description, func(t *testing.T) {
			assert := assert.New(t)
			presenter := NewCreateSeriesPresenter()
			got := presenter.Output(test.Input)
			assert.Equal(test.Want, got)
		})
	}
}
