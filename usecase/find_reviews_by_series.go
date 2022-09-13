package usecase

import (
	"context"
	"series/domain"
	"time"
)

type (
	FindReviewsBySeriesUseCase interface {
		Execute(context.Context, domain.SeriesID) (FindReviewsBySeriesOutput, error)
	}

	FindReviewsBySeriesReview struct {
		ID       string `json:"id"`
		AuthorID string `json:"author_id"`
		Text     string `json:"text"`
	}

	FindReviewsBySeriesOutput struct {
		Reviews []FindReviewsBySeriesReview `json:"reviews"`
	}

	FindReviewsBySeriesPresenter interface {
		Output([]domain.Review) FindReviewsBySeriesOutput
	}

	findReviewsBySeriesInteractor struct {
		series    domain.SeriesRepository
		reviews   domain.ReviewRepository
		presenter FindReviewsBySeriesPresenter
		timeout   time.Duration
	}
)

func NewFindReviewsBySeriesInteractor(
	series domain.SeriesRepository,
	reviews domain.ReviewRepository,
	presenter FindReviewsBySeriesPresenter,
	timeout time.Duration,
) FindReviewsBySeriesUseCase {
	return findReviewsBySeriesInteractor{
		series:    series,
		reviews:   reviews,
		presenter: presenter,
		timeout:   timeout,
	}
}

func (i findReviewsBySeriesInteractor) Execute(
	ctx context.Context,
	seriesID domain.SeriesID,
) (FindReviewsBySeriesOutput, error) {
	ctx, cancel := context.WithTimeout(ctx, i.timeout)
	defer cancel()

	var err error
	var reviews []domain.Review

	err = i.reviews.WithTransaction(ctx, func(ctx context.Context) error {
		_, err = i.series.FindByID(ctx, seriesID)
		if err != nil {
			return err
		}

		reviews, err = i.reviews.FindBySeries(ctx, seriesID)
		return err
	})

	if err != nil {
		return i.presenter.Output(nil), err
	}
	return i.presenter.Output(reviews), nil
}
