package gorilla

import (
	"context"
	"net/http"
	"series/adapter/api/action"
	"series/adapter/api/middleware"
	"series/adapter/logger"
	"series/adapter/presenter"
	"series/adapter/repository"
	"series/adapter/validator"
	"series/usecase"
	"time"

	"github.com/gorilla/mux"
)

type service struct {
	repo      repository.Repository
	logger    logger.Logger
	validator validator.Validator
	dbTimeout time.Duration
	router    *mux.Router
}

func NewHandler(
	repo repository.Repository,
	logger logger.Logger,
	validator validator.Validator,
	dbTimeout time.Duration,
) http.Handler {
	service := &service{
		repo:      repo,
		logger:    logger,
		validator: validator,
		dbTimeout: dbTimeout,
		router:    &mux.Router{},
	}

	router := mux.NewRouter()
	api := router.PathPrefix("/v1").Subrouter()

	api.Use(middleware.Logging(logger))
	api.Use(middleware.CORS)

	api.Handle("/series", service.buildCreateSeriesAction()).Methods(http.MethodPost)
	api.Handle("/series", service.buildFindSeriesByTitleAction()).
		Queries("q", "{query}").
		Methods(http.MethodGet)
	api.Handle("/series/{id}", service.buildFindSeriesByIDAction()).
		Methods(http.MethodGet)
	api.Handle("/series/{id}/reviews", service.buildReviewsBySeriesAction()).
		Methods(http.MethodGet)
	api.Handle("/reviews", service.buildCreateReviewAction()).Methods(http.MethodPost)

	service.router = router
	return service
}

func (s *service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *service) buildCreateSeriesAction() http.Handler {
	f := func(w http.ResponseWriter, r *http.Request) {
		uc := usecase.NewCreateSeriesInteractor(
			s.repo.NewSeriesRepository(),
			presenter.NewCreateSeriesPresenter(),
			s.dbTimeout,
		)
		action := action.NewCreateSeriesAction(uc, s.validator)
		action.Execute(w, r)
	}
	return http.HandlerFunc(f)
}

func (s *service) buildCreateReviewAction() http.Handler {
	f := func(w http.ResponseWriter, r *http.Request) {
		uc := usecase.NewCreateReviewInteractor(
			s.repo.NewSeriesRepository(),
			s.repo.NewReviewRepository(),
			presenter.NewCreateReviewPresenter(),
			s.dbTimeout,
		)
		action := action.NewCreateReviewAction(uc, s.validator)
		action.Execute(w, r)
	}
	return http.HandlerFunc(f)
}

func (s *service) buildFindSeriesByTitleAction() http.Handler {
	f := func(w http.ResponseWriter, r *http.Request) {
		query := mux.Vars(r)["query"]
		r = r.WithContext(
			context.WithValue(r.Context(), action.CtxKeyTitleQuery, query),
		)
		uc := usecase.NewFindSeriesByTitleInteractor(
			s.repo.NewSeriesRepository(),
			presenter.NewFindSeriesByTitlePresenter(),
			s.dbTimeout,
		)
		action := action.NewFindSeriesByTitleAction(uc)
		action.Execute(w, r)
	}
	return http.HandlerFunc(f)
}

func (s *service) buildFindSeriesByIDAction() http.Handler {
	f := func(w http.ResponseWriter, r *http.Request) {
		seriesID := mux.Vars(r)["id"]
		r = r.WithContext(
			context.WithValue(r.Context(), action.CtxKeySeriesID, seriesID),
		)
		uc := usecase.NewFindSeriesByIDInteractor(
			s.repo.NewSeriesRepository(),
			s.repo.NewReviewRepository(),
			presenter.NewFindSeriesByIDPresenter(),
			s.dbTimeout,
		)
		action := action.NewFindSeriesByIDAction(uc)
		action.Execute(w, r)
	}
	return http.HandlerFunc(f)
}

func (s *service) buildReviewsBySeriesAction() http.Handler {
	f := func(w http.ResponseWriter, r *http.Request) {
		seriesID := mux.Vars(r)["id"]
		r = r.WithContext(
			context.WithValue(r.Context(), action.CtxKeySeriesID, seriesID),
		)
		uc := usecase.NewFindReviewsBySeriesInteractor(
			s.repo.NewSeriesRepository(),
			s.repo.NewReviewRepository(),
			presenter.NewFindReviewsBySeriesPresenter(),
			s.dbTimeout,
		)
		action := action.NewFindReviewsBySeriesAction(uc)
		action.Execute(w, r)
	}
	return http.HandlerFunc(f)
}
