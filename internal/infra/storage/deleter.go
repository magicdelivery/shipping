package storage

import (
	"context"
	"shipping/internal/domain/repository"

	redis "github.com/redis/go-redis/v9"
)

type RedisDeleter struct {
	rdb    *redis.Client
	loader repository.CustomerLoader
}

func NewRedisDeleter(rdb *redis.Client, loader repository.CustomerLoader) *RedisDeleter {
	return &RedisDeleter{
		rdb,
		loader,
	}
}

func (rd *RedisDeleter) DeleteCustomer(ctx context.Context, id string) error {
	customer, err := rd.loader.LoadCustomerById(ctx, id)
	if err != nil {
		return err
	}

	if customer.Address != nil {
		if err := rd.rdb.Unlink(ctx, makeCustomerAddressId(customer.Address.Id)).Err(); err != nil {
			return err
		}
	}

	if err := rd.rdb.Unlink(ctx, makeCustomerId(id)).Err(); err != nil {
		return err
	}

	if err := rd.rdb.ZRem(ctx, CUSTOMER_LIST, id).Err(); err != nil {
		return err
	}

	return nil
}
