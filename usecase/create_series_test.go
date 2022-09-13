package usecase

import (
	"context"
	"errors"
	"series/domain"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type mockCreateSeriesRepository struct {
	domain.SeriesRepository
	result domain.Series
	err    error
}

func (r mockCreateSeriesRepository) Create(
	_ context.Context,
	_ domain.Series,
) (domain.Series, error) {
	return r.result, r.err
}

type mockCreateSeriesPresenter struct {
	result CreateSeriesOutput
}

func (p mockCreateSeriesPresenter) Output(domain.Series) CreateSeriesOutput {
	return p.result
}

func TestCreateSeriesInteractor(t *testing.T) {
	t.Parallel()

	testSeries := domain.NewSeries(
		domain.SeriesID("1be9775b-8d32-4710-9ce6-7ece88e30f01"),
		"Title",
		"Description",
		20,
		1980,
		1990,
		"Creator",
	)
	testInput := CreateSeriesInput{
		Title:       "",
		Description: "",
		Episodes:    0,
		BeginYear:   0,
		EndYear:     0,
		Creator:     "",
	}
	testOutput := CreateSeriesOutput{
		ID:          testSeries.ID().String(),
		Title:       testSeries.Title(),
		Description: testSeries.Description(),
		Episodes:    testSeries.Episodes(),
		BeginYear:   testSeries.BeginYear(),
		EndYear:     testSeries.EndYear(),
		Creator:     testSeries.Creator(),
	}
	testErr := errors.New("Error")

	type Test struct {
		Description string
		Repo        domain.SeriesRepository
		Presenter   CreateSeriesPresenter
		Input       CreateSeriesInput
		Expected    CreateSeriesOutput
		ExpectedErr any
	}
	tests := []Test{
		{
			Description: "Successful creation",
			Repo: mockCreateSeriesRepository{
				result: testSeries,
				err:    nil,
			},
			Presenter: mockCreateSeriesPresenter{
				result: testOutput,
			},
			Input:       testInput,
			Expected:    testOutput,
			ExpectedErr: nil,
		},
		{
			Description: "Some error",
			Repo: mockCreateSeriesRepository{
				result: domain.Series{},
				err:    testErr,
			},
			Presenter: mockCreateSeriesPresenter{
				result: CreateSeriesOutput{},
			},
			Input:       testInput,
			Expected:    CreateSeriesOutput{},
			ExpectedErr: testErr,
		},
	}

	for _, test := range tests {
		t.Run(test.Description, func(t *testing.T) {
			assert := assert.New(t)
			uc := NewCreateSeriesInteractor(
				test.Repo,
				test.Presenter,
				1*time.Second,
			)
			output, err := uc.Execute(context.TODO(), test.Input)
			assert.Equal(test.ExpectedErr, err)
			assert.Equal(test.Expected, output)
		})
	}
}
