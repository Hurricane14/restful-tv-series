package postgres

import (
	"context"
	"errors"
	"series/domain"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
)

type seriesRepository struct {
	db *DB
}

// Create implements domain.SeriesRepository
func (r *seriesRepository) Create(
	ctx context.Context,
	series domain.Series,
) (domain.Series, error) {
	var execer interface {
		Exec(context.Context, string, ...any) (pgconn.CommandTag, error)
	} = r.db.pool

	tx, ok := ctx.Value(CtxKeyTx).(pgx.Tx)
	if ok {
		execer = tx
	}

	const query = `
  INSERT INTO
    series(id, title, description, episodes, begin_year, end_year, creator)
  VALUES
    ($1, $2, $3, $4, $5, $6, $7)
  `

	_, err := execer.Exec(
		ctx,
		query,
		series.ID(),
		series.Title(),
		series.Description(),
		series.Episodes(),
		series.BeginYear(),
		series.EndYear(),
		series.Creator(),
	)
	if err != nil {
		return domain.Series{}, err
	}
	return series, nil
}

// FindByID implements domain.SeriesRepository
func (r *seriesRepository) FindByID(
	ctx context.Context,
	ID domain.SeriesID,
) (domain.Series, error) {
	var (
		id                 string
		title              string
		description        string
		episodes           int
		beginYear, endYear int
		creator            string
		querier            interface {
			QueryRow(context.Context, string, ...any) pgx.Row
		} = r.db.pool
	)

	tx, ok := ctx.Value(CtxKeyTx).(pgx.Tx)
	if ok {
		querier = tx
	}

	const query = `
    SELECT
      id, title, description, episodes, begin_year, end_year, creator
    FROM series
    WHERE id = $1
  `
	row := querier.QueryRow(ctx, query, ID)
	err := row.Scan(
		&id,
		&title,
		&description,
		&episodes,
		&beginYear, &endYear,
		&creator,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.Series{}, domain.ErrSeriesNotFound
	} else if err != nil {
		return domain.Series{}, err
	}
	return domain.NewSeries(
		domain.SeriesID(id),
		title,
		description,
		episodes,
		beginYear, endYear,
		creator,
	), nil
}

// FindByTitle implements domain.SeriesRepository
func (r *seriesRepository) FindByTitle(ctx context.Context, title string) ([]domain.Series, error) {
	var querier interface {
		Query(context.Context, string, ...any) (pgx.Rows, error)
	} = r.db.pool

	tx, ok := ctx.Value(CtxKeyTx).(pgx.Tx)
	if ok {
		querier = tx
	}

	const query = `
    SELECT
      id, title, description, episodes, begin_year, end_year, creator
    FROM series
    WHERE
      make_tsvector(title, description) @@ to_tsquery($1)
  `
	rows, err := querier.Query(ctx, query, title)
	if errors.Is(err, pgx.ErrNoRows); err != nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	series := []domain.Series{}
	for rows.Next() {
		var (
			id                 string
			title              string
			description        string
			episodes           int
			beginYear, endYear int
			creator            string
		)
		err := rows.Scan(
			&id,
			&title,
			&description,
			&episodes,
			&beginYear, &endYear,
			&creator,
		)
		if err != nil {
			return nil, err
		}
		series = append(series, domain.NewSeries(
			domain.SeriesID(id),
			title,
			description,
			episodes,
			beginYear,
			endYear,
			creator,
		))
	}
	return series, nil
}
