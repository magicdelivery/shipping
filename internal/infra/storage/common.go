package storage

import (
	"fmt"
	"shipping/internal/infra/config"

	redis "github.com/redis/go-redis/v9"
)

const (
	CUSTOMER_LIST = "customers"
)

func NewRedisClient(config *config.Config) *redis.Client {
	addr := fmt.Sprintf("%s:%d",
		config.RedisStorage.Host,
		config.RedisStorage.Port)
	return redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
		DB:       0,
	})
}

func DefaultServices(config *config.Config) (*RedisLoader, *RedisSaver, *RedisDeleter) {
	rdb := NewRedisClient(config)
	loader := NewRedisLoader(rdb)
	saver := NewRedisSaver(rdb)
	deleter := NewRedisDeleter(rdb, loader)
	return loader, saver, deleter
}

func makeCustomerId(id string) string {
	return fmt.Sprintf("customer:%s", id)
}

func makeCustomerAddressId(id string) string {
	return fmt.Sprintf("customer_shipping_address:%s", id)
}
