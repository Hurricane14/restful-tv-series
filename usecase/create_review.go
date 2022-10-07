package usecase

import (
	"context"
	"series/domain"
	"time"
)

type (
	CreateReviewUseCase interface {
		Execute(context.Context, CreateReviewInput) (CreateReviewOutput, error)
	}

	CreateReviewInput struct {
		SeriesID string `json:"series_id" validate:"required,uuid_rfc4122"`
		AuthorID string `json:"author_id" validate:"required,uuid_rfc4122"`
		Text     string `json:"text"      validate:"required,max=500"`
	}

	CreateReviewOutput struct {
		ID       string `json:"id"`
		SeriesID string `json:"series_id"`
		AuthorID string `json:"author_id"`
		Text     string `json:"text"`
	}

	CreateReviewPresenter interface {
		Output(domain.Review) CreateReviewOutput
	}

	createReviewInteractor struct {
		series    domain.SeriesRepository
		reviews   domain.ReviewRepository
		presenter CreateReviewPresenter
		timeout   time.Duration
	}
)

func NewCreateReviewInteractor(
	series domain.SeriesRepository,
	reviews domain.ReviewRepository,
	presenter CreateReviewPresenter,
	timeout time.Duration,
) CreateReviewUseCase {
	return createReviewInteractor{
		series:    series,
		reviews:   reviews,
		presenter: presenter,
		timeout:   timeout,
	}
}

func (i createReviewInteractor) Execute(
	ctx context.Context, input CreateReviewInput,
) (CreateReviewOutput, error) {
	ctx, cancel := context.WithTimeout(ctx, i.timeout)
	defer cancel()

	var (
		err    error
		review domain.Review
	)

	err = i.reviews.WithTransaction(ctx, func(ctx context.Context) error {
		_, err = i.series.FindByID(ctx, domain.SeriesID(input.SeriesID))
		if err != nil {
			return err
		}

		err = i.reviews.Reviewed(
			ctx,
			domain.SeriesID(input.SeriesID),
			domain.AuthorID(input.AuthorID),
		)
		if err != nil {
			return err
		}

		review, err = i.reviews.Create(ctx, domain.NewReview(
			domain.ReviewID(domain.NewUUID()),
			domain.SeriesID(input.SeriesID),
			domain.AuthorID(input.AuthorID),
			input.Text,
		))
		return err
	})

	if err != nil {
		return i.presenter.Output(domain.Review{}), err
	}
	return i.presenter.Output(review), nil
}
