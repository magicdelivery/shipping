package parcel_locker

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"log"
	"shipping/internal/infra/config"
	"shipping/internal/infra/storage"
	"time"

	"github.com/eko/gocache/lib/v4/cache"
	"github.com/eko/gocache/lib/v4/marshaler"
	"github.com/eko/gocache/lib/v4/store"
	redis_store "github.com/eko/gocache/store/redis/v4"
)

const KEY_TPL = "cache:parcel_locker:%s"

type RedisCacher struct {
	marshaler *marshaler.Marshaler
	exp       time.Duration
}

func NewRedisCacher(config *config.Config) *RedisCacher {
	rdb := storage.NewRedisClient(config)
	redis_store := redis_store.NewRedis(rdb)
	cacheManager := cache.New[any](redis_store)
	marshaler := marshaler.New(cacheManager)
	exp := time.Duration(config.ParcelLockerService.CacheTtl) * time.Second
	return &RedisCacher{
		marshaler: marshaler,
		exp:       exp,
	}
}

func (rc *RedisCacher) Set(ctx context.Context, key string, value ParcelLockersNear) error {
	internalKey := rc.makeKey(key)
	err := rc.marshaler.Set(ctx, internalKey, value, store.WithExpiration(rc.exp))
	if err != nil {
		log.Printf("Near parcel lockers cannot be saved in cache, key: %s, err: %v\n", key, err)
		return err
	}
	return nil
}

func (rc *RedisCacher) Get(ctx context.Context, key string) (ParcelLockersNear, error) {
	internalKey := rc.makeKey(key)
	value, err := rc.marshaler.Get(ctx, internalKey, new(ParcelLockersNear))
	if err != nil {
		log.Printf("Near parcel lockers for key \"%s\" are not found in cache: %v\n", key, err)
		return ParcelLockersNear{}, err
	}
	log.Printf("Near parcel lockers for key \"%s\" are read from cache\n", key)
	return *value.(*ParcelLockersNear), nil
}

func (rc *RedisCacher) makeKey(key string) string {
	hashedKey := rc.hash(key)
	return fmt.Sprintf(KEY_TPL, hashedKey)
}

func (rc *RedisCacher) hash(input string) string {
	hash := md5.New()
	hash.Write([]byte(input))
	hashBytes := hash.Sum(nil)
	output := hex.EncodeToString(hashBytes)
	return output
}
