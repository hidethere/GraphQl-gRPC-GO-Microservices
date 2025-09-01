package catalog

import (
	"context"

	"github.com/segmentio/ksuid"
)

type Service interface {
	PostProduct(ctx context.Context, name, description string, price float64) (*Product, error)
	GetProductById(ctx context.Context, id string) (*Product, error)
	GetProducts(ctx context.Context, skip uint64, take uint64) ([]Product, error)
	GetProductWithIDs(ctx context.Context, ids []string) ([]Product, error)
	SearchProducts(ctx context.Context, query string, skip uint64, take uint64) ([]Product, error)
}

type catalogService struct {
	repository Repository
}

type Product struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

func NewService(r Repository) Service {
	return &catalogService{r}
}

func (c *catalogService) GetProductById(ctx context.Context, id string) (*Product, error) {
	res, err := c.repository.GetProductByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *catalogService) GetProductWithIDs(ctx context.Context, ids []string) ([]Product, error) {
	res, err := c.repository.ListProductWithIDs(ctx, ids)
	if err != nil {
		return nil, err
	}
	return res, err
}

func (c *catalogService) GetProducts(ctx context.Context, skip uint64, take uint64) ([]Product, error) {
	if take > 100 || (skip == 0 && take == 0) {
		take = 100
	}

	res, err := c.repository.ListProducts(ctx, skip, take)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *catalogService) PostProduct(ctx context.Context, name, description string, price float64) (*Product, error) {
	p := &Product{
		Name:        name,
		Description: description,
		Price:       price,
		ID:          ksuid.New().Next().String(),
	}
	if err := c.repository.PutProduct(ctx, *p); err != nil {
		return nil, err
	}
	return p, nil

}

func (c *catalogService) SearchProducts(ctx context.Context, query string, skip uint64, take uint64) ([]Product, error) {
	if take > 100 || (skip == 0 && take == 0) {
		take = 100
	}
	res, err := c.repository.SearchProducts(ctx, query, skip, take)
	if err != nil {
		return nil, err
	}
	return res, nil
}
