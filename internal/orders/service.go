package orders

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5"
	repo "github.com/rmcampos/ecom/internal/adapters/postgres/sqlc"
)

// constructor arguments type
type svc struct {
	repo *repo.Queries
	db   *pgx.Conn
}

// constructor logic
func NewService(repo *repo.Queries, db *pgx.Conn) Service {
	return &svc{
		repo: repo,
		db:   db,
	}
}

// methods declaration/signatures
type Service interface {
	CreateOrder(ctx context.Context, orderRequest createOrderParam) (repo.Order, error)
}

// methods implementation
func (s *svc) CreateOrder(ctx context.Context, orderRequest createOrderParam) (repo.Order, error) {
	// validate payload
	if orderRequest.CustomerID <= 0 {
		return repo.Order{}, fmt.Errorf("invalid customer ID: %d", orderRequest.CustomerID)
	}
	if len(orderRequest.Items) == 0 {
		return repo.Order{}, fmt.Errorf("no items in the order. Orders must contain at least one item")
	}
	for _, item := range orderRequest.Items {
		if item.Quantity <= 0 {
			return repo.Order{}, fmt.Errorf("invalid item quantity: %d", item.Quantity)
		}
	}

	// start db transaction
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return repo.Order{}, fmt.Errorf("failed to start transaction: %w", err)
	}
	defer func() {
		if err := tx.Rollback(ctx); err != nil {
			slog.Error("Failed to rollback database transaction", "error", err, "connection", s.db.Config())
		}
	}()

	qtx := s.repo.WithTx(tx)

	// create order
	order, err := qtx.CreateOrder(ctx, orderRequest.CustomerID)
	if err != nil {
		return repo.Order{}, fmt.Errorf("failed to create order: %w", err)
	}

	// create order items
	for _, item := range orderRequest.Items {
		product, err := qtx.FindProductByID(ctx, item.ProductID)
		if err != nil {
			return repo.Order{}, fmt.Errorf("failed to find product: %w", err)
		}

		if product.Quantity < item.Quantity {
			return repo.Order{}, fmt.Errorf("insufficient stock for product ID %d", item.ProductID)
		}

		_, err = qtx.CreateOrderItem(ctx, repo.CreateOrderItemParams{
			OrderID:    order.ID,
			ProductID:  item.ProductID,
			Quantity:   item.Quantity,
			PriceCents: product.PriceInCents,
		})
		if err != nil {
			return repo.Order{}, fmt.Errorf("failed to create order item: %w", err)
		}
	}

	// TODO: update product quantity in stock after order creation (not implemented yet)

	// commit transaction
	if err := tx.Commit(ctx); err != nil {
		return repo.Order{}, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return order, nil
}
