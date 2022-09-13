package repository

import "series/domain"

type Repository interface {
	NewSeriesRepository() domain.SeriesRepository
	NewReviewRepository() domain.ReviewRepository
}
