package action

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"shipping/internal/domain/model"
	"shipping/internal/domain/repository/mocks"
)

func TestLoadCustomerById(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name             string
		id               string
		mockReturn       *model.Customer
		mockErr          error
		expectedStatus   int
		expectedResponse map[string]interface{}
	}{
		{
			name:             "Happy Path",
			id:               "1",
			mockReturn:       &model.Customer{Id: "1", Name: "user_1", CreatedAt: 10000001},
			mockErr:          nil,
			expectedStatus:   http.StatusOK,
			expectedResponse: map[string]interface{}{"customer": map[string]interface{}{"address": nil, "id": "1", "name": "user_1", "created_at": float64(10000001)}},
		},
		{
			name:             "Not Found",
			id:               "2",
			mockReturn:       nil,
			mockErr:          nil,
			expectedStatus:   http.StatusNotFound,
			expectedResponse: map[string]interface{}{"id": "2", "message": "Customer not found by id: 2"},
		},
		{
			name:             "Internal Server Error",
			id:               "3",
			mockReturn:       nil,
			mockErr:          fmt.Errorf("some error"),
			expectedStatus:   http.StatusInternalServerError,
			expectedResponse: map[string]interface{}{"id": "3", "message": "some error"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			r := gin.Default()

			customerLoaderMock := mocks.NewCustomerLoader(t)
			customerLoaderMock.
				On("LoadCustomerById", mock.Anything, tt.id).
				Return(tt.mockReturn, tt.mockErr)

			req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/customer/%s", tt.id), nil)
			w := httptest.NewRecorder()

			r.GET("/customer/:id", func(c *gin.Context) {
				LoadCustomerById(c, customerLoaderMock)
			})

			// Act
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tt.expectedStatus, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			for key, expectedValue := range tt.expectedResponse {
				assert.Equal(t, expectedValue, response[key])
			}

			customerLoaderMock.AssertExpectations(t)
		})
	}
}
