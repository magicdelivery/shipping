package storage

import (
	"context"
	"strconv"

	"shipping/internal/domain/model"

	redis "github.com/redis/go-redis/v9"
)

type RedisLoader struct {
	rdb *redis.Client
}

func NewRedisLoader(rdb *redis.Client) *RedisLoader {
	return &RedisLoader{
		rdb,
	}
}

func (rl *RedisLoader) LoadCustomerById(ctx context.Context, id string) (*model.Customer, error) {
	shippingHGetAll := rl.rdb.HGetAll(ctx, makeCustomerId(id))
	if err := shippingHGetAll.Err(); err != nil {
		return nil, err
	}
	shippingRes, err := shippingHGetAll.Result()
	if err != nil {
		return nil, err
	}
	if len(shippingRes) == 0 {
		return nil, nil
	}

	createdAt, _ := strconv.ParseInt(shippingRes["CreatedAt"], 10, 64)
	shipping := model.Customer{
		Id:        shippingRes["Id"],
		Name:      shippingRes["Name"],
		Address:   nil,
		CreatedAt: createdAt,
	}

	if addressId, is := shippingRes["AddressId"]; is {
		addressHGetAll := rl.rdb.HGetAll(ctx, makeCustomerAddressId(addressId))
		if err := addressHGetAll.Err(); err != nil {
			return nil, err
		}
		addressRes, err := addressHGetAll.Result()
		if err != nil {
			return nil, err
		}
		if len(addressRes) > 0 {
			longitude, err := strconv.ParseFloat(addressRes["Longitude"], 64)
			if err != nil {
				return nil, err
			}
			latitude, err := strconv.ParseFloat(addressRes["Latitude"], 64)
			if err != nil {
				return nil, err
			}
			address := model.ShippingAddress{
				Id:        addressRes["Id"],
				City:      addressRes["City"],
				Street:    addressRes["Street"],
				Longitude: longitude,
				Latitude:  latitude,
			}
			shipping.Address = &address
		}
	}

	return &shipping, nil
}

func (rl *RedisLoader) LoadAllCustomers(ctx context.Context) ([]*model.Customer, error) {
	var shippings []*model.Customer = make([]*model.Customer, 0)
	zRange := rl.rdb.ZRange(ctx, CUSTOMER_LIST, 0, -1)
	if err := zRange.Err(); err != nil {
		return nil, err
	}

	ids, err := zRange.Result()
	if err != nil {
		return nil, err
	}

	for _, id := range ids {
		if shipping, err := rl.LoadCustomerById(ctx, id); err != nil {
			return nil, err
		} else {
			shippings = append(shippings, shipping)
		}
	}
	return shippings, nil
}
