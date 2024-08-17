package action

import (
	"fmt"
	"net/http"
	"strconv"

	"shipping/internal/domain/repository"
	"shipping/internal/infra/parcel_locker"

	"github.com/gin-gonic/gin"
)

func LoadAllCustomers(ctx *gin.Context, loader repository.CustomerLoader) {
	if customers, err := loader.LoadAllCustomers(ctx); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
	} else {
		response := gin.H{"customers": customers}
		ctx.JSON(http.StatusOK, response)
	}
}

func LoadCustomerById(ctx *gin.Context, loader repository.CustomerLoader) {
	id := ctx.Params.ByName("id")
	customer, err := loader.LoadCustomerById(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"id": id, "message": err.Error()})
	} else if customer == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"id": id, "message": fmt.Sprintf("Customer not found by id: %s", id)})
	} else {
		ctx.JSON(http.StatusOK, gin.H{"customer": customer})
	}
}

func FindParcelLockersByCustomerId(
	ctx *gin.Context,
	loader repository.CustomerLoader,
	plClient parcel_locker.ParcelLockerClient,
) {
	id := ctx.Params.ByName("id")
	distanceStr := ctx.DefaultQuery("distance", "10")
	distance, err := strconv.ParseFloat(distanceStr, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid distance value"})
		// c.Abort()
		return
	}

	customer, err := loader.LoadCustomerById(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"id": id, "message": err.Error()})
	} else if customer == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"id": id, "message": fmt.Sprintf("Customer not found by id: %s", id)})
	} else {
		parcel_lockers, err := plClient.FindParcelLockersNear(ctx, customer, distance)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"id": id, "message": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"parcel_lockers": parcel_lockers})
	}
}
