package postgres

import (
	"context"
	"errors"
	"series/domain"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
)

type reviewRepository struct {
	db *DB
}

func (r *reviewRepository) Create(
	ctx context.Context,
	review domain.Review,
) (domain.Review, error) {
	var execer interface {
		Exec(context.Context, string, ...any) (pgconn.CommandTag, error)
	} = r.db.pool

	tx, ok := ctx.Value(CtxKeyTx).(pgx.Tx)
	if ok {
		execer = tx
	}

	const query = `
    INSERT INTO
      reviews(id, series_id, author_id, text)
    VALUES
      ($1, $2, $3, $4)
  `

	_, err := execer.Exec(
		ctx,
		query,
		review.ID(),
		review.SeriesID(),
		review.AuthorID(),
		review.Text(),
	)
	if err != nil {
		return domain.Review{}, err
	}
	return review, nil
}

func (r *reviewRepository) FindBySeries(
	ctx context.Context,
	seriesID domain.SeriesID,
) ([]domain.Review, error) {
	var querier interface {
		Query(context.Context, string, ...any) (pgx.Rows, error)
	} = r.db.pool

	tx, ok := ctx.Value(CtxKeyTx).(pgx.Tx)
	if ok {
		querier = tx
	}

	const query = `
    SELECT
      id, series_id, author_id, text
    FROM reviews
    WHERE
      series_id = $1
  `

	rows, err := querier.Query(ctx, query, seriesID)
	if errors.Is(err, pgx.ErrNoRows); err != nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	reviews := []domain.Review{}
	for rows.Next() {
		var (
			id        string
			series_id string
			author_id string
			text      string
		)
		err := rows.Scan(
			&id,
			&series_id,
			&author_id,
			&text,
		)
		if err != nil {
			return nil, err
		}
		reviews = append(reviews, domain.NewReview(
			domain.ReviewID(id),
			domain.SeriesID(series_id),
			domain.AuthorID(author_id),
			text,
		))
	}
	return reviews, nil
}

func (r *reviewRepository) Reviewed(
	ctx context.Context,
	seriesID domain.SeriesID,
	authorID domain.AuthorID,
) error {
	var querier interface {
		QueryRow(context.Context, string, ...any) pgx.Row
	} = r.db.pool

	tx, ok := ctx.Value(CtxKeyTx).(pgx.Tx)
	if ok {
		querier = tx
	}

	const query = `
    SELECT 
      id
    FROM reviews
    WHERE
      series_id = $1 and author_id = $2
  `

	var id string
	row := querier.QueryRow(ctx, query, seriesID, authorID)
	err := row.Scan(&id)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil
	} else if err != nil {
		return err
	} else {
		return domain.ErrAlreadyReviewed
	}
}

func (r *reviewRepository) WithTransaction(
	ctx context.Context,
	fn func(context.Context) error,
) error {
	tx, err := r.db.pool.Begin(ctx)
	if err != nil {
		return err
	}
	ctx = context.WithValue(ctx, CtxKeyTx, tx)
	if err := fn(ctx); err != nil {
		rbErr := tx.Rollback(ctx)
		// TODO: wrap err into rbErr?
		if rbErr != nil {
			return rbErr
		}
		return err
	}
	return tx.Commit(ctx)
}
