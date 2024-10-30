package route

import (
	"shipping/internal/app/action"
	"shipping/internal/infra/config"
	"shipping/internal/infra/parcel_locker"
	"shipping/internal/infra/storage"

	"github.com/gin-gonic/gin"
)

func SetupRouter(config *config.Config) *gin.Engine {
	r := gin.Default()
	loader, saver, deleter := storage.DefaultServices(&config.RedisStorage)
	plClient := parcel_locker.NewParcelLockerClient(config)

	r.GET("/ping", action.Ping())
	r.GET("/customers", action.LoadAllCustomers(loader))
	r.GET("/customer/:id", action.LoadCustomerById(loader))
	r.GET("/customer-parcel-lockers/:id", action.FindParcelLockersByCustomerId(loader, plClient))
	r.POST("/customer", action.SaveCustomer(saver))
	r.DELETE("/customer/:id", action.DeleteCustomer(deleter, loader))
	return r
}
