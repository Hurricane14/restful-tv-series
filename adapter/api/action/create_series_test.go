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

type mockCreateSeriesUseCase struct {
	output usecase.CreateSeriesOutput
	err    error
}

func (uc mockCreateSeriesUseCase) Execute(
	context.Context,
	usecase.CreateSeriesInput,
) (usecase.CreateSeriesOutput, error) {
	return uc.output, uc.err
}

func TestCreateSeriesAction(t *testing.T) {
	t.Parallel()

	type Test struct {
		Description  string
		UC           usecase.CreateSeriesUseCase
		ExpectedCode int
		ExpectedBody any
	}
	tests := []Test{
		{
			Description: "Successfull creation",
			UC: mockCreateSeriesUseCase{
				output: usecase.CreateSeriesOutput{
					ID:          "ID",
					Title:       "Title",
					Description: "Description",
					Episodes:    20,
					BeginYear:   1970,
					EndYear:     1980,
					Creator:     "Creator",
				},
				err: nil,
			},
			ExpectedCode: http.StatusCreated,
			ExpectedBody: usecase.CreateSeriesOutput{
				ID:          "ID",
				Title:       "Title",
				Description: "Description",
				Episodes:    20,
				BeginYear:   1970,
				EndYear:     1980,
				Creator:     "Creator",
			},
		},

		{
			Description: "Generic error",
			UC: mockCreateSeriesUseCase{
				output: usecase.CreateSeriesOutput{},
				err:    errors.New("error"),
			},
			ExpectedCode: http.StatusInternalServerError,
			ExpectedBody: errorResponse{
				Errors: []string{"error"},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.Description, func(t *testing.T) {
			assert := assert.New(t)

			input, err := json.Marshal(usecase.CreateSeriesInput{})
			assert.Nil(err)
			req, err := http.NewRequest(http.MethodPost, "", bytes.NewReader(input))
			assert.Nil(err)
			recorder := httptest.NewRecorder()

			action := NewCreateSeriesAction(test.UC, mockValidator{})
			action.Execute(recorder, req)

			assert.Equal(test.ExpectedCode, recorder.Code)
			if recorder.Code >= http.StatusInternalServerError {
				assert.Equal(test.ExpectedCode, recorder.Code)
			} else if recorder.Code >= http.StatusBadRequest {
				output := errorResponse{}
				assert.Nil(json.Unmarshal(recorder.Body.Bytes(), &output))
				assert.Equal(test.ExpectedBody, output)
			} else {
				output := usecase.CreateSeriesOutput{}
				assert.Nil(json.Unmarshal(recorder.Body.Bytes(), &output))
				assert.Equal(test.ExpectedBody, output)
			}
		})
	}
}
