package repository

import (
	// New import
	// New import

	//"fmt"

	"errors"

	"test/internal/models"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
	ErrNoOrderPlaced  = errors.New("order not placed")
	ErrDuplicateEmail = errors.New("duplicate email")
)

type IStoreStorage interface {
	Create(order *models.Order) error
	Delete(id int64) error
	GetByID(id int64) (*models.Order, error)
	GetInventory() (map[string]int, error)
}
