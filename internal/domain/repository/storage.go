package repository

import (
	"context"
	"shipping/internal/domain/model"
)

//go:generate go run github.com/vektra/mockery/v2@v2.43.2 --name=CustomerLoader
type CustomerLoader interface {
	LoadCustomerById(ctx context.Context, id string) (*model.Customer, error)
	LoadAllCustomers(ctx context.Context) ([]*model.Customer, error)
}

//go:generate go run github.com/vektra/mockery/v2@v2.43.2 --name=CustomerSaver
type CustomerSaver interface {
	SaveCustomer(ctx context.Context, shipping model.Customer) error
}

//go:generate go run github.com/vektra/mockery/v2@v2.43.2 --name=CustomerDeleter
type CustomerDeleter interface {
	DeleteCustomer(ctx context.Context, id string) error
}
