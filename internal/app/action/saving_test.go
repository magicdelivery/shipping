package action

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"shipping/internal/domain/model"
	"shipping/internal/domain/repository/mocks"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSaveCustomer_FailOnBindJson(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	jsonBody := `invalid-json`
	customerSaverMock := mocks.NewCustomerSaver(t)

	req, _ := http.NewRequest(http.MethodPost, "/customer", strings.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.POST("/customer", func(c *gin.Context) {
		SaveCustomer(c, customerSaverMock)
	})

	// Act
	r.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.False(t, response["created"].(bool))
	assert.Equal(t, "invalid character 'i' looking for beginning of value", response["message"].(string))

	// customerSaverMock.AssertExpectations(t)
}

func TestSaveCustomer_FailOnSave(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	customerSaverMock := mocks.NewCustomerSaver(t)
	customerSaverMock.
		On("SaveCustomer", mock.Anything, mock.Anything).
		Times(1).
		Return(fmt.Errorf("some redis error"))

	customerInput := model.Customer{
		ID:        "1",
		Name:      "user_1",
		CreatedAt: 10000001,
		Address: &model.ShippingAddress{
			ID:        "1",
			City:      "Riga",
			Street:    "Brivibas str 1",
			Longitude: 59.0,
			Latitude:  24.0,
		},
	}
	jsonData, _ := json.Marshal(customerInput)
	req, _ := http.NewRequest(http.MethodPost, "/customer", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	r.POST("/customer", func(c *gin.Context) {
		SaveCustomer(c, customerSaverMock)
	})

	// Act
	r.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.False(t, response["created"].(bool))
	assert.Equal(t, "some redis error", response["message"].(string))

	customerSaverMock.AssertExpectations(t)
}

func TestSaveCustomer_HappyPath(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	customerSaverMock := mocks.NewCustomerSaver(t)

	customerInput := model.Customer{
		ID:        "1",
		Name:      "user_1",
		CreatedAt: 10000001,
		Address: &model.ShippingAddress{
			ID:        "1",
			City:      "Riga",
			Street:    "Brivibas str 1",
			Longitude: 59.0,
			Latitude:  24.0,
		},
	}

	customerSaverMock.
		On("SaveCustomer", mock.Anything, customerInput).
		Once().
		Return(nil)

	jsonData, _ := json.Marshal(customerInput)
	req, _ := http.NewRequest(http.MethodPost, "/customer", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	r.POST("/customer", func(c *gin.Context) {
		SaveCustomer(c, customerSaverMock)
	})

	// Act
	r.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusCreated, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.True(t, response["created"].(bool))
	assert.Equal(t, "Customer info created successfully", response["message"].(string))
	assert.NotNil(t, response["customer"])

	customerData, err := json.Marshal(response["customer"])
	assert.NoError(t, err)
	var customerResponse model.Customer
	err = json.Unmarshal(customerData, &customerResponse)
	assert.NoError(t, err)
	assert.True(t, reflect.DeepEqual(customerInput, customerResponse))

	customerSaverMock.AssertExpectations(t)
}
