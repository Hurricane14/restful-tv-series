package action

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"series/usecase"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockFindSeriesByTitleUseCase struct {
	output usecase.FindSeriesByTitleOutput
	err    error
}

func (uc mockFindSeriesByTitleUseCase) Execute(
	context.Context,
	string,
) (usecase.FindSeriesByTitleOutput, error) {
	return uc.output, uc.err
}

func TestFindSeriesByTitleAction(t *testing.T) {
	t.Parallel()

	type Test struct {
		Description  string
		UC           usecase.FindSeriesByTitleUseCase
		ExpectedCode int
		ExpectedBody any
	}
	tests := []Test{
		{
			Description: "Successful search",
			UC: mockFindSeriesByTitleUseCase{
				output: usecase.FindSeriesByTitleOutput{
					Series: []usecase.FindSeriesByTitleSeries{
						{
							ID:        "SeriesID",
							Title:     "Title",
							BeginYear: 1980,
							EndYear:   1990,
							Creator:   "Creator",
						},
					},
				},
				err: nil,
			},
			ExpectedCode: http.StatusOK,
			ExpectedBody: usecase.FindSeriesByTitleOutput{
				Series: []usecase.FindSeriesByTitleSeries{
					{
						ID:        "SeriesID",
						Title:     "Title",
						BeginYear: 1980,
						EndYear:   1990,
						Creator:   "Creator",
					},
				},
			},
		},

		{
			Description: "Generic error",
			UC: mockFindSeriesByTitleUseCase{
				output: usecase.FindSeriesByTitleOutput{},
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
				CtxKeyTitleQuery,
				"query",
			)
			req = req.WithContext(ctx)

			recorder := httptest.NewRecorder()

			action := NewFindSeriesByTitleAction(test.UC)
			action.Execute(recorder, req)

			assert.Equal(test.ExpectedCode, recorder.Code)
			if recorder.Code >= http.StatusInternalServerError {
				assert.Equal(test.ExpectedCode, recorder.Code)
			} else if recorder.Code >= http.StatusBadRequest {
				output := errorResponse{}
				assert.Nil(json.Unmarshal(recorder.Body.Bytes(), &output))
				assert.Equal(test.ExpectedBody, output)
			} else {
				output := usecase.FindSeriesByTitleOutput{}
				assert.Nil(json.Unmarshal(recorder.Body.Bytes(), &output))
				assert.Equal(test.ExpectedBody, output)
			}
		})
	}
}
