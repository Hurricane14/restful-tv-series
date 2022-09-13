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

type mockCreateReviewUsecase struct {
	output usecase.CreateReviewOutput
	err    error
}

func (uc mockCreateReviewUsecase) Execute(
	context.Context,
	usecase.CreateReviewInput,
) (usecase.CreateReviewOutput, error) {
	return uc.output, uc.err
}

func TestCreateReviewAction(t *testing.T) {
	t.Parallel()

	type Test struct {
		Description  string
		UC           usecase.CreateReviewUseCase
		ExpectedCode int
		ExpectedBody any
	}
	tests := []Test{
		{
			Description: "Successful creation",
			UC: mockCreateReviewUsecase{
				output: usecase.CreateReviewOutput{
					ID:       "ID",
					SeriesID: "SeriesID",
					AuthorID: "AuthorID",
					Text:     "Text",
				},
				err: nil,
			},
			ExpectedCode: http.StatusCreated,
			ExpectedBody: usecase.CreateReviewOutput{
				ID:       "ID",
				SeriesID: "SeriesID",
				AuthorID: "AuthorID",
				Text:     "Text",
			},
		},

		{
			Description: "Create review for series that does not exist",
			UC: mockCreateReviewUsecase{
				output: usecase.CreateReviewOutput{},
				err:    domain.ErrSeriesNotFound,
			},
			ExpectedCode: http.StatusBadRequest,
			ExpectedBody: errorResponse{
				Errors: []string{domain.ErrSeriesNotFound.Error()},
			},
		},

		{
			Description: "Generic error",
			UC: mockCreateReviewUsecase{
				output: usecase.CreateReviewOutput{},
				err:    errors.New("error"),
			},
			ExpectedCode: http.StatusInternalServerError,
			ExpectedBody: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.Description, func(t *testing.T) {
			assert := assert.New(t)

			input, err := json.Marshal(usecase.CreateReviewInput{})
			assert.Nil(err)
			req, err := http.NewRequest(http.MethodPost, "", bytes.NewReader(input))
			assert.Nil(err)
			recorder := httptest.NewRecorder()

			action := NewCreateReviewAction(test.UC, mockValidator{})
			action.Execute(recorder, req)

			assert.Equal(test.ExpectedCode, recorder.Code)
			if recorder.Code >= http.StatusInternalServerError {
				assert.Equal(test.ExpectedCode, recorder.Code)
			} else if recorder.Code >= http.StatusBadRequest {
				output := errorResponse{}
				assert.Nil(json.Unmarshal(recorder.Body.Bytes(), &output))
				assert.Equal(test.ExpectedBody, output)
			} else {
				output := usecase.CreateReviewOutput{}
				assert.Nil(json.Unmarshal(recorder.Body.Bytes(), &output))
				assert.Equal(test.ExpectedBody, output)
			}
		})
	}
}
