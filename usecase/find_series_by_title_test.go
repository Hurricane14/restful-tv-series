package usecase

import (
	"context"
	"series/domain"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type mockFindSeriesByTitleRepo struct {
	domain.SeriesRepository
	series []domain.Series
	err    error
}

func (r mockFindSeriesByTitleRepo) FindByTitle(
	_ context.Context,
	title string,
) ([]domain.Series, error) {
	return r.series, r.err
}

type mockFindSeriesByTitlePresenter struct {
	output FindSeriesByTitleOutput
}

func (p mockFindSeriesByTitlePresenter) Output(
	_ []domain.Series,
) FindSeriesByTitleOutput {
	return p.output
}

func TestFindSeriesByTitleInteractor(t *testing.T) {
	t.Parallel()

	type Test struct {
		Description string
		Repo        domain.SeriesRepository
		Presenter   FindSeriesByTitlePresenter
		Expected    FindSeriesByTitleOutput
		ExpectedErr error
	}
	tests := []Test{
		{
			Description: "Successful search",
			Repo: mockFindSeriesByTitleRepo{
				series: []domain.Series{
					domain.NewSeries(
						"ID1",
						"Title1",
						"Description",
						20,
						1980,
						1990,
						"Creator1",
					),
					domain.NewSeries(
						"ID2",
						"Title2",
						"Description",
						13,
						1970,
						1978,
						"Creator2",
					),
				},
				err: nil,
			},
			Presenter: mockFindSeriesByTitlePresenter{
				output: FindSeriesByTitleOutput{
					Series: []FindSeriesByTitleSeries{
						{
							ID:        "ID1",
							Title:     "Title1",
							BeginYear: 1980,
							EndYear:   1990,
							Creator:   "Creator1",
						},
						{
							ID:        "ID2",
							Title:     "Title2",
							BeginYear: 1970,
							EndYear:   1975,
							Creator:   "Creator2",
						},
					},
				},
			},
			Expected: FindSeriesByTitleOutput{
				Series: []FindSeriesByTitleSeries{
					{
						ID:        "ID1",
						Title:     "Title1",
						BeginYear: 1980,
						EndYear:   1990,
						Creator:   "Creator1",
					},
					{
						ID:        "ID2",
						Title:     "Title2",
						BeginYear: 1970,
						EndYear:   1975,
						Creator:   "Creator2",
					},
				},
			},
			ExpectedErr: nil,
		},

		{
			Description: "No series means empty slice, not nil",
			Repo:        mockFindSeriesByTitleRepo{},
			Presenter: mockFindSeriesByTitlePresenter{
				output: FindSeriesByTitleOutput{
					Series: []FindSeriesByTitleSeries{},
				},
			},
			Expected: FindSeriesByTitleOutput{
				Series: []FindSeriesByTitleSeries{},
			},
			ExpectedErr: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.Description, func(t *testing.T) {
			assert := assert.New(t)
			uc := NewFindSeriesByTitleInteractor(
				test.Repo,
				test.Presenter,
				1*time.Second,
			)
			got, err := uc.Execute(context.TODO(), "")
			assert.Equal(test.ExpectedErr, err)
			assert.Equal(test.Expected, got)
		})
	}
}
