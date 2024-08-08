package storage

import (
	"context"
	"shipping/internal/domain/model"

	redis "github.com/redis/go-redis/v9"
)

type RedisSaver struct {
	rdb *redis.Client
}

func NewRedisSaver(rdb *redis.Client) *RedisSaver {
	return &RedisSaver{
		rdb,
	}
}

func (rs *RedisSaver) SaveCustomer(ctx context.Context, shipping model.Customer) error {
	// TODO: add transaction "multi", consider rollback
	hsetShippingId := makeCustomerId(shipping.Id)
	hset := rs.rdb.HSet(ctx, hsetShippingId, "Id", shipping.Id)
	if hset.Err() != nil {
		return hset.Err()
	}

	hset = rs.rdb.HSet(ctx, hsetShippingId, "Name", shipping.Name)
	if hset.Err() != nil {
		return hset.Err()
	}

	hset = rs.rdb.HSet(ctx, hsetShippingId, "CreatedAt", shipping.CreatedAt)
	if hset.Err() != nil {
		return hset.Err()
	}

	if shipping.Address != nil {
		address := shipping.Address
		hsetAddressId := makeCustomerAddressId(address.Id)
		if hset := rs.rdb.HSet(ctx, hsetAddressId, "Id", address.Id); hset.Err() != nil {
			return hset.Err()
		}
		if hset := rs.rdb.HSet(ctx, hsetAddressId, "City", address.City); hset.Err() != nil {
			return hset.Err()
		}
		if hset := rs.rdb.HSet(ctx, hsetAddressId, "Street", address.Street); hset.Err() != nil {
			return hset.Err()
		}
		if hset := rs.rdb.HSet(ctx, hsetAddressId, "Longitude", address.Longitude); hset.Err() != nil {
			return hset.Err()
		}
		if hset := rs.rdb.HSet(ctx, hsetAddressId, "Latitude", address.Latitude); hset.Err() != nil {
			return hset.Err()
		}
		// Set FK reference for `customer`
		if hset := rs.rdb.HSet(ctx, hsetShippingId, "AddressId", address.Id); hset.Err() != nil {
			return hset.Err()
		}
	}

	z := redis.Z{Score: float64(shipping.CreatedAt), Member: shipping.Id}
	zadd := rs.rdb.ZAdd(ctx, CUSTOMER_LIST, z)
	if zadd.Err() != nil {
		return zadd.Err()
	}

	return nil
}
