package domain

import (
	"context"
	"errors"
)

type AuthorID string

func (id AuthorID) String() string {
	return string(id)
}

type ReviewID string

func (id ReviewID) String() string {
	return string(id)
}

var (
	ErrReviewNotFound  = errors.New("review not found")
	ErrAlreadyReviewed = errors.New("this author already reviewed this series")
)

type (
	ReviewRepository interface {
		Create(context.Context, Review) (Review, error)
		FindBySeries(context.Context, SeriesID) ([]Review, error)
		Reviewed(context.Context, SeriesID, AuthorID) error
		WithTransaction(context.Context, func(context.Context) error) error
	}

	Review struct {
		id       ReviewID
		seriesID SeriesID
		author   AuthorID
		text     string
	}
)

func NewReview(
	ID ReviewID,
	seriesID SeriesID,
	authorID AuthorID,
	text string,
) Review {
	return Review{
		id:       ID,
		seriesID: seriesID,
		author:   authorID,
		text:     text,
	}
}

func (r *Review) ID() ReviewID {
	return r.id
}

func (r *Review) SeriesID() SeriesID {
	return r.seriesID
}

func (r *Review) AuthorID() AuthorID {
	return r.author
}

func (r *Review) Text() string {
	return r.text
}
