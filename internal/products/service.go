package products

import (
	"context"
	"fmt"

	repo "github.com/rmcampos/ecom/internal/adapters/postgres/sqlc"
)

// constructor arguments type
type svc struct {
	repo repo.Querier
}

// constructor logic
func NewService(repo repo.Querier) Service {
	return &svc{repo: repo}
}

// methods declaration/signatures
type Service interface {
	ListProducts(ctx context.Context) ([]repo.Product, error)
	FindProductByID(ctx context.Context, id int64) (repo.Product, error)
	CreateProduct(ctx context.Context, product createProductDto) (repo.Product, error)
}

// methods implementation
func (s *svc) ListProducts(ctx context.Context) ([]repo.Product, error) {
	return s.repo.ListProducts(ctx)
}

func (s *svc) FindProductByID(ctx context.Context, id int64) (repo.Product, error) {
	return s.repo.FindProductByID(ctx, id)
}

func (s *svc) CreateProduct(ctx context.Context, product createProductDto) (repo.Product, error) {
	// validations
	// check if name is nil or empty
	if product.Name == "" {
		return repo.Product{}, fmt.Errorf("name is required")
	}
	// check if price is nil or negative
	if product.PriceInCents <= 0 {
		return repo.Product{}, fmt.Errorf("price is required")
	}
	// check if quantity is nil or negative
	if product.Quantity <= 0 {
		return repo.Product{}, fmt.Errorf("quantity is required")
	}

	// create product
	productCreated, err := s.repo.CreateProduct(ctx, repo.CreateProductParams{
		Name:         product.Name,
		PriceInCents: product.PriceInCents,
		Quantity:     product.Quantity,
	})
	if err != nil {
		return repo.Product{}, fmt.Errorf("failed to create product: %w", err)
	}
	return productCreated, nil
}
