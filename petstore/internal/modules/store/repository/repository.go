package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"

	"test/internal/models"
)

type StoreStorage struct {
	logger *zap.Logger
	DB     *sqlx.DB
}

func NewStoreStorage(db *sqlx.DB, logger *zap.Logger) IStoreStorage {
	return &StoreStorage{
		logger: logger,
		DB:     db}
}

func (ps *StoreStorage) Create(order *models.Order) error {
	query := `
	INSERT INTO orders (pet_id, quantity, ship_date, status, complete) 
	VALUES ($1, $2, $3, $4, $5) 
	RETURNING id`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []any{
		order.PetID,
		order.Quantity,
		order.ShipDate,
		order.Status,
		order.Complete,
	}

	err := ps.DB.QueryRowContext(ctx, query, args...).Scan(&order.ID)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			ps.logger.Error("no rows in orders table", zap.Error(err))
			return ErrEditConflict
		default:
			ps.logger.Error(" error on inserting order", zap.Error(err))
			return err
		}
	}

	return nil
}

func (ps *StoreStorage) GetByID(id int64) (*models.Order, error) {
	query := `
	SELECT * FROM orders
	WHERE id = $1`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	order := &models.Order{}

	err := ps.DB.QueryRowContext(ctx, query, id).Scan(&order.ID, &order.PetID, &order.Quantity, &order.ShipDate, &order.Status, &order.Complete)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			ps.logger.Error("no rows in orders table", zap.Error(err))
			return nil, ErrNoOrderPlaced
		default:
			ps.logger.Error(" error on inserting order", zap.Error(err))
			return nil, err
		}
	}

	return order, nil
}

func (ps *StoreStorage) Delete(id int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := ps.DB.ExecContext(ctx, `
		DELETE FROM orders WHERE id = $1
	`, id)
	if err != nil {
		ps.logger.Error("error on deleting order", zap.Error(err))
		return err
	}

	return nil
}

func (ps *StoreStorage) GetInventory() (map[string]int, error) {
	inventory := make(map[string]int)
	query := `
SELECT status, count(*) OVER(PARTITION BY status)
FROM pets
`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	rows, err := ps.DB.QueryContext(ctx, query)
	if err != nil {
		ps.logger.Error("some error on getting rows from pet tags", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var status string
		var count int
		err := rows.Scan(&status, &count)
		if err != nil {
			ps.logger.Error("some error on getting rows from pet tags", zap.Error(err))
			return nil, err
		}
		inventory[status] = count

	}
	return inventory, nil

}
