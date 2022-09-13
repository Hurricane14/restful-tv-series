package postgres

import (
	"context"
	"fmt"
	"series/domain"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

type DB struct {
	pool *pgxpool.Pool
}

func NewDB(c *Config) (*DB, error) {
	dbURL := fmt.Sprintf(
		"postgresql://%s:%s@%s:%d/%s?sslmode=%s",
		c.user, c.password, c.host, c.port, c.database, c.sslmode,
	)

	withRetries := func(fn func() error, sleep time.Duration, tries int) error {
		var err error
		for {
			if tries == 0 {
				return err
			}
			err = fn()
			if err == nil {
				return nil
			}
			time.Sleep(sleep)
			tries--
		}
	}

	var err error
	var pool *pgxpool.Pool
	withRetries(func() error {
		pool, err = pgxpool.Connect(context.Background(), dbURL)
		return err
	}, 5*time.Second, 5)
	if err != nil {
		return nil, err
	}

	return &DB{pool}, nil
}

func (db *DB) Close() {
	db.pool.Close()
}

func (db *DB) NewSeriesRepository() domain.SeriesRepository {
	return &seriesRepository{
		db: db,
	}
}

func (db *DB) NewReviewRepository() domain.ReviewRepository {
	return &reviewRepository{
		db: db,
	}
}
