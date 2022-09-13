package usecase

import (
	"context"
	"series/domain"
	"time"
)

type (
	CreateSeriesUseCase interface {
		Execute(context.Context, CreateSeriesInput) (CreateSeriesOutput, error)
	}

	CreateSeriesInput struct {
		Title       string `json:"title"       validate:"required,max=70"`
		Description string `json:"description" validate:"required,max=200"`
		Episodes    int    `json:"episodes"    validate:"required,min=1"`
		BeginYear   int    `json:"begin_year"  validate:"required,min=1946,max=2030"`
		EndYear     int    `json:"end_year"    validate:"eq=0|gtefield=BeginYear"`
		Creator     string `json:"creator"     validate:"required,min=5,max=30"`
	}

	CreateSeriesOutput struct {
		ID          string `json:"id"`
		Title       string `json:"title"`
		Description string `json:"description"`
		Episodes    int    `json:"episodes"`
		BeginYear   int    `json:"begin_year"`
		EndYear     int    `json:"end_year"`
		Creator     string `json:"creator"`
	}

	CreateSeriesPresenter interface {
		Output(domain.Series) CreateSeriesOutput
	}

	createSeriesInteractor struct {
		repo      domain.SeriesRepository
		presenter CreateSeriesPresenter
		timeout   time.Duration
	}
)

func NewCreateSeriesInteractor(
	repo domain.SeriesRepository,
	presenter CreateSeriesPresenter,
	timeout time.Duration,
) CreateSeriesUseCase {

	return createSeriesInteractor{
		repo:      repo,
		presenter: presenter,
		timeout:   timeout,
	}
}

func (i createSeriesInteractor) Execute(
	ctx context.Context, input CreateSeriesInput,
) (CreateSeriesOutput, error) {

	ctx, cancel := context.WithTimeout(ctx, i.timeout)
	defer cancel()

	series := domain.NewSeries(
		domain.SeriesID(domain.NewUUID()),
		input.Title, input.Description,
		input.Episodes,
		input.BeginYear, input.EndYear,
		input.Creator,
	)

	series, err := i.repo.Create(ctx, series)
	if err != nil {
		return i.presenter.Output(domain.Series{}), err
	}

	return i.presenter.Output(series), nil
}
