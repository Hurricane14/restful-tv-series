package usecase

import (
	"context"
	"series/domain"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type mockFindReviewsBySeriesSeriesRepo struct {
	domain.SeriesRepository
	series domain.Series
	err    error
}

func (r mockFindReviewsBySeriesSeriesRepo) FindByID(
	_ context.Context,
	_ domain.SeriesID,
) (domain.Series, error) {
	return r.series, r.err
}

type mockFindReviewsBySeriesReviewRepo struct {
	domain.ReviewRepository
	reviews []domain.Review
	err     error
}

func (r mockFindReviewsBySeriesReviewRepo) WithTransaction(
	ctx context.Context,
	fn func(ctx context.Context) error,
) error {
	return fn(ctx)
}

func (r mockFindReviewsBySeriesReviewRepo) FindBySeries(
	_ context.Context,
	_ domain.SeriesID,
) ([]domain.Review, error) {
	return r.reviews, r.err
}

type mockFindReviewBySeriesPresenter struct {
	output FindReviewsBySeriesOutput
}

func (p mockFindReviewBySeriesPresenter) Output(_ []domain.Review) FindReviewsBySeriesOutput {
	return p.output
}

func TestFindReviewsBySeriesInteractor(t *testing.T) {
	t.Parallel()

	type Test struct {
		Description string
		Series      domain.SeriesRepository
		Reviews     domain.ReviewRepository
		Presenter   FindReviewsBySeriesPresenter
		Expected    FindReviewsBySeriesOutput
		ExpectedErr error
	}
	tests := []Test{
		{
			Description: "Successful finding of reviews",
			Series: mockFindReviewsBySeriesSeriesRepo{
				series: domain.Series{},
				err:    nil,
			},
			Reviews: mockFindReviewsBySeriesReviewRepo{
				ReviewRepository: nil,
				reviews: []domain.Review{
					domain.NewReview(
						"ID",
						"SeriesID",
						"AuthorID",
						"Text",
					),
				},
				err: nil,
			},
			Presenter: mockFindReviewBySeriesPresenter{
				output: FindReviewsBySeriesOutput{
					Reviews: []FindReviewsBySeriesReview{
						{
							ID:       "ID",
							AuthorID: "AuthorID",
							Text:     "Text",
						},
					},
				},
			},
			Expected: FindReviewsBySeriesOutput{
				Reviews: []FindReviewsBySeriesReview{
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
			Series: mockFindReviewsBySeriesSeriesRepo{
				series: domain.Series{},
				err:    nil,
			},
			Reviews: mockFindReviewsBySeriesReviewRepo{
				reviews: nil,
				err:     nil,
			},
			Presenter: mockFindReviewBySeriesPresenter{
				output: FindReviewsBySeriesOutput{
					Reviews: []FindReviewsBySeriesReview{},
				},
			},
			Expected: FindReviewsBySeriesOutput{
				Reviews: []FindReviewsBySeriesReview{},
			},
			ExpectedErr: nil,
		},

		{
			Description: "Searching reviews for series that does not exist",
			Series: mockFindReviewsBySeriesSeriesRepo{
				series: domain.Series{},
				err:    domain.ErrSeriesNotFound,
			},
			Reviews:     mockFindReviewsBySeriesReviewRepo{},
			Presenter:   mockFindReviewBySeriesPresenter{},
			Expected:    FindReviewsBySeriesOutput{},
			ExpectedErr: domain.ErrSeriesNotFound,
		},
	}

	for _, test := range tests {
		t.Run(test.Description, func(t *testing.T) {
			assert := assert.New(t)
			uc := NewFindReviewsBySeriesInteractor(
				test.Series,
				test.Reviews,
				test.Presenter,
				1*time.Second,
			)
			got, err := uc.Execute(context.TODO(), domain.SeriesID(""))
			assert.Equal(test.ExpectedErr, err)
			assert.Equal(test.Expected, got)
		})
	}
}
