package usecase

import (
	"context"
	"series/domain"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type mockFindSeriesByIDSeriesRepo struct {
	domain.SeriesRepository
	series domain.Series
	err    error
}

func (r mockFindSeriesByIDSeriesRepo) FindByID(
	_ context.Context,
	_ domain.SeriesID,
) (domain.Series, error) {
	return r.series, r.err
}

type mockFindSeriesByIDReviewRepo struct {
	domain.ReviewRepository
	reviews []domain.Review
	err     error
}

func (r mockFindSeriesByIDReviewRepo) WithTransaction(
	ctx context.Context,
	fn func(ctx context.Context) error,
) error {
	return fn(ctx)
}

func (r mockFindSeriesByIDReviewRepo) FindBySeries(
	_ context.Context,
	_ domain.SeriesID,
) ([]domain.Review, error) {
	return r.reviews, r.err
}

type mockFindSeriesByIDPresenter struct {
	output FindSeriesByIDOutput
}

func (p mockFindSeriesByIDPresenter) Output(
	_ domain.Series,
	_ []domain.Review,
) FindSeriesByIDOutput {
	return p.output
}

func TestFindSeriesByIDInteractor(t *testing.T) {
	t.Parallel()

	type Test struct {
		Description string
		Series      domain.SeriesRepository
		Reviews     domain.ReviewRepository
		Presenter   FindSeriesByIDPresenter
		Expected    FindSeriesByIDOutput
		ExpectedErr error
	}
	tests := []Test{
		{
			Description: "Successful search",
			Series: mockFindSeriesByIDSeriesRepo{
				series: domain.NewSeries(
					"ID",
					"Title",
					"Description",
					20,
					1980,
					1990,
					"Creator",
				),
				err: nil,
			},
			Reviews: mockFindSeriesByIDReviewRepo{
				reviews: []domain.Review{
					domain.NewReview(
						"ID",
						"ID",
						"AuthorID",
						"Text",
					),
				},
				err: nil,
			},
			Presenter: mockFindSeriesByIDPresenter{
				output: FindSeriesByIDOutput{
					ID:          "ID",
					Title:       "Title",
					Description: "Description",
					Episodes:    20,
					BeginYear:   1980,
					EndYear:     1990,
					Creator:     "Creator",
					Reviews: []FindSeriesByIDReview{
						{
							ID:       "ID",
							AuthorID: "AuthorID",
							Text:     "Text",
						},
					},
				},
			},
			Expected: FindSeriesByIDOutput{
				ID:          "ID",
				Title:       "Title",
				Description: "Description",
				Episodes:    20,
				BeginYear:   1980,
				EndYear:     1990,
				Creator:     "Creator",
				Reviews: []FindSeriesByIDReview{
					{
						ID:       "ID",
						AuthorID: "AuthorID",
						Text:     "Text",
					},
				},
			},
			ExpectedErr: nil,
		},

		{
			Description: "No reviews means empty slice, not nil",
			Series: mockFindSeriesByIDSeriesRepo{
				series: domain.NewSeries(
					"ID",
					"Title",
					"Description",
					20,
					1980,
					1990,
					"Creator",
				),
				err: nil,
			},
			Reviews: mockFindSeriesByIDReviewRepo{
				reviews: nil,
				err:     nil,
			},
			Presenter: mockFindSeriesByIDPresenter{
				output: FindSeriesByIDOutput{
					ID:          "ID",
					Title:       "Title",
					Description: "Description",
					Episodes:    20,
					BeginYear:   1980,
					EndYear:     1990,
					Creator:     "Creator",
					Reviews:     []FindSeriesByIDReview{},
				},
			},
			Expected: FindSeriesByIDOutput{
				ID:          "ID",
				Title:       "Title",
				Description: "Description",
				Episodes:    20,
				BeginYear:   1980,
				EndYear:     1990,
				Creator:     "Creator",
				Reviews:     []FindSeriesByIDReview{},
			},
			ExpectedErr: nil,
		},

		{
			Description: "Series with provided id does not exist",
			Series: mockFindSeriesByIDSeriesRepo{
				err: domain.ErrSeriesNotFound,
			},
			Reviews:     mockFindSeriesByIDReviewRepo{},
			Presenter:   mockFindSeriesByIDPresenter{},
			Expected:    FindSeriesByIDOutput{},
			ExpectedErr: domain.ErrSeriesNotFound,
		},
	}

	for _, test := range tests {
		t.Run(test.Description, func(t *testing.T) {
			assert := assert.New(t)
			uc := NewFindSeriesByIDInteractor(
				test.Series,
				test.Reviews,
				test.Presenter,
				1*time.Second,
			)
			got, err := uc.Execute(context.TODO(), "")
			assert.Equal(test.ExpectedErr, err)
			assert.Equal(test.Expected, got)
		})
	}
}
