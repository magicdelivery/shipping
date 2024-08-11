package action

import (
	"net/http"

	"shipping/internal/domain/model"
	"shipping/internal/domain/repository"

	"github.com/gin-gonic/gin"
)

func SaveCustomer(ctx *gin.Context, saver repository.CustomerSaver) {
	var customer model.Customer = model.Customer{}
	if err := ctx.BindJSON(&customer); err != nil {
		// ctx.BindJSON sets status code 400, thus the next code definition does not effect
		ctx.JSON(http.StatusInternalServerError, map[string]any{"customer": customer, "created": false, "message": err.Error()})
		return
	}

	if err := saver.SaveCustomer(ctx, customer); err != nil {
		ctx.JSON(http.StatusInternalServerError, map[string]any{"customer": customer, "created": false, "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"customer": customer, "created": true, "message": "Customer info created successfully"})
}
