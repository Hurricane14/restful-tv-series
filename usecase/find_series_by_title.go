package usecase

import (
	"context"
	"series/domain"
	"time"
)

type (
	FindSeriesByTitleUseCase interface {
		Execute(context.Context, string) (FindSeriesByTitleOutput, error)
	}

	FindSeriesByTitleSeries struct {
		ID        string `json:"id"`
		Title     string `json:"title"`
		BeginYear int    `json:"begin_year"`
		EndYear   int    `json:"end_year"`
		Creator   string `json:"creator"`
	}

	FindSeriesByTitleOutput struct {
		Series []FindSeriesByTitleSeries `json:"series"`
	}

	FindSeriesByTitlePresenter interface {
		Output([]domain.Series) FindSeriesByTitleOutput
	}

	findSeriesByTitleInteractor struct {
		repo      domain.SeriesRepository
		presenter FindSeriesByTitlePresenter
		timeout   time.Duration
	}
)

func NewFindSeriesByTitleInteractor(
	series domain.SeriesRepository,
	presenter FindSeriesByTitlePresenter,
	timeout time.Duration,
) FindSeriesByTitleUseCase {
	return findSeriesByTitleInteractor{
		repo:      series,
		presenter: presenter,
		timeout:   timeout,
	}
}

func (s findSeriesByTitleInteractor) Execute(
	ctx context.Context, query string,
) (FindSeriesByTitleOutput, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	series, err := s.repo.FindByTitle(ctx, query)
	if err != nil {
		return s.presenter.Output(nil), err
	}
	return s.presenter.Output(series), nil
}
