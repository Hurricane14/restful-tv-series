package domain

import (
	"context"
	"errors"
)

type SeriesID string

func (id SeriesID) String() string {
	return string(id)
}

var (
	ErrSeriesNotFound = errors.New("series not found")
)

type (
	SeriesRepository interface {
		Create(context.Context, Series) (Series, error)
		FindByTitle(context.Context, string) ([]Series, error)
		FindByID(context.Context, SeriesID) (Series, error)
	}

	Series struct {
		id                 SeriesID
		title              string
		description        string
		episodes           int
		beginYear, endYear int
		creator            string
	}
)

func NewSeries(
	ID SeriesID,
	title, description string,
	episodes int,
	beginYear, endYear int,
	creator string,
) Series {
	return Series{
		id:          ID,
		title:       title,
		description: description,
		episodes:    episodes,
		beginYear:   beginYear,
		endYear:     endYear,
		creator:     creator,
	}
}

func (s *Series) ID() SeriesID {
	return s.id
}

func (s *Series) Title() string {
	return s.title
}

func (s *Series) Description() string {
	return s.description
}

func (s *Series) Episodes() int {
	return s.episodes
}

func (s *Series) BeginYear() int {
	return s.beginYear
}

func (s *Series) EndYear() int {
	return s.endYear
}

func (s *Series) Creator() string {
	return s.creator
}
