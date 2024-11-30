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

func TestDeleteCustomer_NotFound(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)

	r := gin.Default()
	customerLoaderMock := mocks.NewCustomerLoader(t)
	customerLoaderMock.
		On("LoadCustomerById", mock.Anything, "1").
		Return(nil, nil)

	customerDeleterMock := mocks.NewCustomerDeleter(t)

	req, _ := http.NewRequest(http.MethodDelete, "/customer/1", nil)
	w := httptest.NewRecorder()

	r.DELETE("/customer/:id", func(c *gin.Context) {
		DeleteCustomer(customerDeleterMock, customerLoaderMock)(c)
	})
	// Act
	r.ServeHTTP(w, req)
	// Assert
	assert.Equal(t, http.StatusNotFound, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	expectedResponse := map[string]interface{}{
		"id":      "1",
		"deleted": false,
		"message": "Customer info not found by id: 1",
	}
	assert.Equal(t, expectedResponse, response)

	customerLoaderMock.AssertExpectations(t)

	// Do not check call of DeleteCustomer, since we don't call it in this scenario
}

func TestDeleteCustomer_Found(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name             string
		id               string
		mockDeleteErr    error
		expectedStatus   int
		expectedResponse map[string]interface{}
	}{
		{
			name:             "Error Deleting Customer",
			id:               "2",
			mockDeleteErr:    fmt.Errorf("delete error"),
			expectedStatus:   http.StatusInternalServerError,
			expectedResponse: map[string]interface{}{"id": "2", "message": "delete error"},
		},
		{
			name:             "Successful Deletion",
			id:               "3",
			mockDeleteErr:    nil,
			expectedStatus:   http.StatusOK,
			expectedResponse: map[string]interface{}{"id": "3", "deleted": true, "message": "Customer info successfully deleted"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			r := gin.Default()

			customerLoaderMock := mocks.NewCustomerLoader(t)
			customerLoaderMock.
				On("LoadCustomerById", mock.Anything, tt.id).
				Return(&model.Customer{ID: tt.id}, nil)

			customerDeleterMock := mocks.NewCustomerDeleter(t)
			customerDeleterMock.
				On("DeleteCustomer", mock.Anything, tt.id).
				Return(tt.mockDeleteErr)

			req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/customer/%s", tt.id), nil)
			w := httptest.NewRecorder()

			r.DELETE("/customer/:id", func(c *gin.Context) {
				DeleteCustomer(customerDeleterMock, customerLoaderMock)(c)
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
			customerDeleterMock.AssertExpectations(t)
		})
	}
}
