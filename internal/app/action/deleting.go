package action

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"shipping/internal/domain/repository"
)

func DeleteCustomer(deleter repository.CustomerDeleter, loader repository.CustomerLoader) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Params.ByName("id")
		if customer, _ := loader.LoadCustomerById(ctx, id); customer == nil {
			ctx.JSON(http.StatusNotFound, gin.H{"id": id, "deleted": false, "message": fmt.Sprintf("Customer info not found by id: %s", id)})
			return
		}

		if err := deleter.DeleteCustomer(ctx, id); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"id": id, "message": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"id": id, "deleted": true, "message": "Customer info successfully deleted"})
	}
}
