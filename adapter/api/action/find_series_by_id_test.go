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

type mockFindSeriesByIDUseCase struct {
	output usecase.FindSeriesByIDOutput
	err    error
}

func (uc mockFindSeriesByIDUseCase) Execute(
	context.Context,
	domain.SeriesID,
) (usecase.FindSeriesByIDOutput, error) {
	return uc.output, uc.err
}

func TestFindSeriesByID(t *testing.T) {
	t.Parallel()

	type Test struct {
		Description  string
		UC           usecase.FindSeriesByIDUseCase
		ExpectedCode int
		ExpectedBody any
	}
	tests := []Test{
		{
			Description: "Successful search",
			UC: mockFindSeriesByIDUseCase{
				output: usecase.FindSeriesByIDOutput{
					ID:          "SeriesID",
					Title:       "Title",
					Description: "Description",
					Episodes:    20,
					BeginYear:   1980,
					EndYear:     1990,
					Creator:     "Creator",
					Reviews: []usecase.FindSeriesByIDReview{
						{
							ID:       "ReviewID",
							AuthorID: "AuthorID",
							Text:     "Text",
						},
					},
				},
				err: nil,
			},
			ExpectedCode: http.StatusOK,
			ExpectedBody: usecase.FindSeriesByIDOutput{
				ID:          "SeriesID",
				Title:       "Title",
				Description: "Description",
				Episodes:    20,
				BeginYear:   1980,
				EndYear:     1990,
				Creator:     "Creator",
				Reviews: []usecase.FindSeriesByIDReview{
					{
						ID:       "ReviewID",
						AuthorID: "AuthorID",
						Text:     "Text",
					},
				},
			},
		},

		{
			Description: "No reviews means empty slice, not nil",
			UC: mockFindSeriesByIDUseCase{
				output: usecase.FindSeriesByIDOutput{
					ID:          "SeriesID",
					Title:       "Title",
					Description: "Description",
					Episodes:    20,
					BeginYear:   1980,
					EndYear:     1990,
					Creator:     "Creator",
					Reviews:     []usecase.FindSeriesByIDReview{},
				},
				err: nil,
			},
			ExpectedCode: http.StatusOK,
			ExpectedBody: usecase.FindSeriesByIDOutput{
				ID:          "SeriesID",
				Title:       "Title",
				Description: "Description",
				Episodes:    20,
				BeginYear:   1980,
				EndYear:     1990,
				Creator:     "Creator",
				Reviews:     []usecase.FindSeriesByIDReview{},
			},
		},

		{
			Description: "Searching series that does not exist",
			UC: mockFindSeriesByIDUseCase{
				output: usecase.FindSeriesByIDOutput{},
				err:    domain.ErrSeriesNotFound,
			},
			ExpectedCode: http.StatusNotFound,
			ExpectedBody: errorResponse{
				Errors: []string{domain.ErrSeriesNotFound.Error()},
			},
		},

		{
			Description: "Generic error",
			UC: mockFindSeriesByIDUseCase{
				output: usecase.FindSeriesByIDOutput{},
				err:    errors.New("error"),
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

			action := NewFindSeriesByIDAction(test.UC)
			action.Execute(recorder, req)

			assert.Equal(test.ExpectedCode, recorder.Code)
			if recorder.Code >= http.StatusInternalServerError {
				assert.Equal(test.ExpectedCode, recorder.Code)
			} else if recorder.Code >= http.StatusBadRequest {
				output := errorResponse{}
				assert.Nil(json.Unmarshal(recorder.Body.Bytes(), &output))
				assert.Equal(test.ExpectedBody, output)
			} else {
				output := usecase.FindSeriesByIDOutput{}
				assert.Nil(json.Unmarshal(recorder.Body.Bytes(), &output))
				assert.Equal(test.ExpectedBody, output)
			}
		})
	}
}
