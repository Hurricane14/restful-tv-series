package usecase

import (
	"context"
	"series/domain"
	"time"
)

type (
	FindSeriesByIDUseCase interface {
		Execute(context.Context, domain.SeriesID) (FindSeriesByIDOutput, error)
	}

	FindSeriesByIDReview struct {
		ID       string `json:"id"`
		AuthorID string `json:"author_id"`
		Text     string `json:"text"`
	}

	FindSeriesByIDOutput struct {
		ID          string                 `json:"id"`
		Title       string                 `json:"title"`
		Description string                 `json:"description"`
		Episodes    int                    `json:"episodes"`
		BeginYear   int                    `json:"begin_year"`
		EndYear     int                    `json:"end_year"`
		Creator     string                 `json:"creator"`
		Reviews     []FindSeriesByIDReview `json:"reviews"`
	}

	FindSeriesByIDPresenter interface {
		Output(domain.Series, []domain.Review) FindSeriesByIDOutput
	}

	findSeriesByIDInteractor struct {
		series    domain.SeriesRepository
		reviews   domain.ReviewRepository
		presenter FindSeriesByIDPresenter
		timeout   time.Duration
	}
)

func NewFindSeriesByIDInteractor(
	series domain.SeriesRepository,
	reviews domain.ReviewRepository,
	presenter FindSeriesByIDPresenter,
	timeout time.Duration,
) FindSeriesByIDUseCase {

	return findSeriesByIDInteractor{
		series:    series,
		reviews:   reviews,
		presenter: presenter,
		timeout:   timeout,
	}
}

func (i findSeriesByIDInteractor) Execute(
	ctx context.Context, seriesID domain.SeriesID,
) (FindSeriesByIDOutput, error) {

	ctx, cancel := context.WithTimeout(ctx, i.timeout)
	defer cancel()

	var err error
	var series domain.Series
	var reviews []domain.Review

	err = i.reviews.WithTransaction(ctx, func(ctx context.Context) error {
		series, err = i.series.FindByID(ctx, seriesID)
		if err != nil {
			return err
		}
		reviews, err = i.reviews.FindBySeries(ctx, seriesID)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return i.presenter.Output(domain.Series{}, nil), err
	}
	return i.presenter.Output(series, reviews), nil

}
