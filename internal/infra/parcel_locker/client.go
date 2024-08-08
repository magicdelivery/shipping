package parcel_locker

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"shipping/internal/domain/model"
	"shipping/internal/infra/config"
	"shipping/internal/infra/operation"
	"time"
)

const (
	parcelLockersDistanceSearchUrlTpl = "%s/parcel-locker-distance-search?longitude=%f&latitude=%f&distance=%f"
)

//go:generate go run github.com/vektra/mockery/v2@v2.43.2 --name=ParcelLockerClient
type ParcelLockerClient interface {
	FindParcelLockersNear(ctx context.Context, shipping *model.Customer, distance float64) (ParcelLockersNear, error)
}

type ParcelLockerHttpClient struct {
	locationServiceEndpoint string
	cacher                  *RedisCacher
}

type ParcelLockersDistanceSearchResponse struct {
	Name      string  `json:"name"`
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
	Distance  float64 `json:"distance"`
}

type ParcelLockersNear struct {
	ParcelLockers []ParcelLockersDistanceSearchResponse
}

func NewParcelLockerClient(config *config.Config) *ParcelLockerHttpClient {
	return &ParcelLockerHttpClient{
		locationServiceEndpoint: config.ParcelLockerService.EndpointUrl,
		cacher:                  NewRedisCacher(config),
	}
}

func (cl *ParcelLockerHttpClient) FindParcelLockersNear(ctx context.Context, shipping *model.Customer, distance float64) (ParcelLockersNear, error) {
	if shipping.Address == nil {
		return ParcelLockersNear{}, nil
	}

	endpoint := fmt.Sprintf(
		parcelLockersDistanceSearchUrlTpl,
		cl.locationServiceEndpoint,
		shipping.Address.Longitude,
		shipping.Address.Latitude,
		distance,
	)

	if parcels, err := cl.cacher.Get(ctx, endpoint); err == nil {
		return parcels, nil
	}

	if parcels, err := cl.fetchParcelLockersNear(endpoint); err == nil {
		cl.cacher.Set(ctx, endpoint, parcels)
		return parcels, nil
	} else {
		return ParcelLockersNear{}, err
	}
}

func (cl *ParcelLockerHttpClient) fetchParcelLockersNear(endpoint string) (ParcelLockersNear, error) {
	fetchFn := func() (result interface{}, err error) {
		httpClient := &http.Client{
			Timeout: 10 * time.Second,
		}
		return httpClient.Get(endpoint)
	}
	runner := &operation.Runner{}
	result, err := runner.RunWithRetries(fetchFn)
	if err != nil {
		return ParcelLockersNear{}, err
	}

	response := result.(*http.Response)
	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		return ParcelLockersNear{}, err
	}

	parcels := ParcelLockersNear{}
	if err := json.Unmarshal(responseData, &parcels.ParcelLockers); err != nil {
		return ParcelLockersNear{}, err
	}

	return parcels, nil
}
