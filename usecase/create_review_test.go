package usecase

import (
	"context"
	"series/domain"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type mockCreateReviewSeriesRepo struct {
	domain.SeriesRepository
	series domain.Series
	err    error
}

func (r mockCreateReviewSeriesRepo) FindByID(
	_ context.Context,
	_ domain.SeriesID,
) (domain.Series, error) {
	return r.series, r.err
}

type mockCreateReviewReviewRepo struct {
	domain.ReviewRepository
	reviewedErr error
	review      domain.Review
	createErr   error
}

func (r mockCreateReviewReviewRepo) WithTransaction(
	ctx context.Context,
	fn func(ctx context.Context) error,
) error {
	return fn(ctx)
}

func (r mockCreateReviewReviewRepo) Create(
	_ context.Context,
	_ domain.Review,
) (domain.Review, error) {
	return r.review, r.createErr
}

func (r mockCreateReviewReviewRepo) Reviewed(
	_ context.Context,
	_ domain.SeriesID,
	_ domain.AuthorID,
) error {
	return r.reviewedErr
}

type mockCreateReviewPresenter struct {
	output CreateReviewOutput
}

func (p mockCreateReviewPresenter) Output(
	_ domain.Review,
) CreateReviewOutput {
	return p.output
}

func TestCreateReviewInteractor(t *testing.T) {
	t.Parallel()

	type Test struct {
		Description string
		Series      domain.SeriesRepository
		Reviews     domain.ReviewRepository
		Presenter   CreateReviewPresenter
		Expected    CreateReviewOutput
		ExpectedErr error
	}
	tests := []Test{
		{
			Description: "Successful creation",
			Series: mockCreateReviewSeriesRepo{
				series: domain.Series{},
				err:    nil,
			},
			Reviews: mockCreateReviewReviewRepo{
				reviewedErr: nil,
				review: domain.NewReview(
					"ID",
					"SeriesID",
					"AuthorID",
					"Text",
				),
				createErr: nil,
			},
			Presenter: mockCreateReviewPresenter{
				output: CreateReviewOutput{
					ID:       "ID",
					SeriesID: "SeriesID",
					AuthorID: "AuthorID",
					Text:     "Text",
				},
			},
			Expected: CreateReviewOutput{
				ID:       "ID",
				SeriesID: "SeriesID",
				AuthorID: "AuthorID",
				Text:     "Text",
			},
			ExpectedErr: nil,
		},

		{
			Description: "Creating review for series that does not exist",
			Series: mockCreateReviewSeriesRepo{
				series: domain.Series{},
				err:    domain.ErrSeriesNotFound,
			},
			Reviews:     mockCreateReviewReviewRepo{},
			Presenter:   mockCreateReviewPresenter{},
			Expected:    CreateReviewOutput{},
			ExpectedErr: domain.ErrSeriesNotFound,
		},

		{
			Description: "Creating second review of series by one author",
			Series:      mockCreateReviewSeriesRepo{},
			Reviews: mockCreateReviewReviewRepo{
				reviewedErr: domain.ErrAlreadyReviewed,
			},
			Presenter:   mockCreateReviewPresenter{},
			Expected:    CreateReviewOutput{},
			ExpectedErr: domain.ErrAlreadyReviewed,
		},
	}

	for _, test := range tests {
		t.Run(test.Description, func(t *testing.T) {
			assert := assert.New(t)
			uc := NewCreateReviewInteractor(
				test.Series,
				test.Reviews,
				test.Presenter,
				1*time.Second,
			)
			got, err := uc.Execute(context.TODO(), CreateReviewInput{})
			assert.Equal(test.ExpectedErr, err)
			assert.Equal(test.Expected, got)
		})
	}
}
