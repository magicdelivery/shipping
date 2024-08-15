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

	r.GET("/ping", func(c *gin.Context) { action.Ping(c) })
	r.GET("/customers", func(c *gin.Context) { action.LoadAllCustomers(c, loader) })
	r.GET("/customer/:id", func(c *gin.Context) { action.LoadCustomerById(c, loader) })
	r.GET("/customer-parcel-lockers/:id", func(c *gin.Context) { action.FindParcelLockersByCustomerId(c, loader, plClient) })
	r.POST("/customer", func(c *gin.Context) { action.SaveCustomer(c, saver) })
	r.DELETE("/customer/:id", func(c *gin.Context) { action.DeleteCustomer(c, deleter, loader) })
	return r
}
