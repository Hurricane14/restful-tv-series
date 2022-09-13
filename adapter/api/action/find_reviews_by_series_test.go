package action

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"series/domain"
	"series/usecase"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockFindReviewsBySeriesUseCase struct {
	output usecase.FindReviewsBySeriesOutput
	err    error
}

func (uc mockFindReviewsBySeriesUseCase) Execute(
	context.Context,
	domain.SeriesID,
) (usecase.FindReviewsBySeriesOutput, error) {
	return uc.output, uc.err
}

func TestFindReviewsBySeriesAction(t *testing.T) {
	t.Parallel()

	type Test struct {
		Description  string
		UC           usecase.FindReviewsBySeriesUseCase
		ExpectedCode int
		ExpectedBody any
	}
	tests := []Test{
		{
			Description: "Successful search",
			UC: mockFindReviewsBySeriesUseCase{
				output: usecase.FindReviewsBySeriesOutput{
					Reviews: []usecase.FindReviewsBySeriesReview{
						{
							ID:       "ID",
							AuthorID: "AuthorID",
							Text:     "Text",
						},
					},
				},
				err: nil,
			},
			ExpectedCode: http.StatusOK,
			ExpectedBody: usecase.FindReviewsBySeriesOutput{
				Reviews: []usecase.FindReviewsBySeriesReview{
					{
						ID:       "ID",
						AuthorID: "AuthorID",
						Text:     "Text",
					},
				},
			},
		},

		{
			Description: "Empty slice",
			UC: mockFindReviewsBySeriesUseCase{
				output: usecase.FindReviewsBySeriesOutput{
					Reviews: []usecase.FindReviewsBySeriesReview{},
				},
				err: nil,
			},
			ExpectedCode: http.StatusOK,
			ExpectedBody: usecase.FindReviewsBySeriesOutput{
				Reviews: []usecase.FindReviewsBySeriesReview{},
			},
		},

		{
			Description: "Series not found",
			UC: mockFindReviewsBySeriesUseCase{
				output: usecase.FindReviewsBySeriesOutput{
					Reviews: []usecase.FindReviewsBySeriesReview{},
				},
				err: domain.ErrSeriesNotFound,
			},
			ExpectedCode: http.StatusNotFound,
			ExpectedBody: errorResponse{
				Errors: []string{domain.ErrSeriesNotFound.Error()},
			},
		},

		{
			Description: "Generic error",
			UC: mockFindReviewsBySeriesUseCase{
				err: errors.New("error"),
			},
			ExpectedCode: http.StatusInternalServerError,
			ExpectedBody: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.Description, func(t *testing.T) {
			assert := assert.New(t)

			input, err := json.Marshal(usecase.CreateSeriesInput{})
			assert.Nil(err)
			req, err := http.NewRequest(http.MethodGet, "", bytes.NewReader(input))
			assert.Nil(err)
			ctx := context.WithValue(
				req.Context(),
				CtxKeySeriesID,
				"1be9775b-8d32-4710-9ce6-7ece88e30f01",
			)
			req = req.WithContext(ctx)

			recorder := httptest.NewRecorder()

			action := NewFindReviewsBySeriesAction(test.UC)
			action.Execute(recorder, req)

			assert.Equal(test.ExpectedCode, recorder.Code)
			if recorder.Code >= http.StatusInternalServerError {
				assert.Equal(test.ExpectedCode, recorder.Code)
			} else if recorder.Code >= http.StatusBadRequest {
				output := errorResponse{}
				assert.Nil(json.Unmarshal(recorder.Body.Bytes(), &output))
				assert.Equal(test.ExpectedBody, output)
			} else {
				output := usecase.FindReviewsBySeriesOutput{}
				assert.Nil(json.Unmarshal(recorder.Body.Bytes(), &output))
				assert.Equal(test.ExpectedBody, output)
			}
		})
	}
}
