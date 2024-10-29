package service

import (
	"errors"
	"test/internal/modules/store/repository"

	"test/internal/models"
)

var (
	ErrRecordNotFound  = errors.New("record not found")
	ErrEditConflict    = errors.New("edit conflict")
	ErrDuplicateRecord = errors.New("duplicate record")
	ErrNoInventory     = errors.New("inventory error")
)

type IStoreService interface {
	Create(pet *models.Order) error
	Delete(id int64) error
	GetByID(id int64) (*models.Order, error)
	GetInventory() (map[string]int, error)
}

type StoreService struct {
	storage repository.IStoreStorage
}

func NewStoreService(repo repository.IStoreStorage) *StoreService {
	return &StoreService{storage: repo}
}

func (s *StoreService) Create(pet *models.Order) error {
	err := s.storage.Create(pet)
	if err != nil {
		return ErrDuplicateRecord
	}
	return nil
}

func (s *StoreService) Delete(id int64) error {
	err := s.storage.Delete(id)
	if err != nil {
		return ErrRecordNotFound
	}
	return nil
}

func (s *StoreService) GetByID(id int64) (*models.Order, error) {
	order, err := s.storage.GetByID(id)
	if err != nil {
		return nil, ErrRecordNotFound
	}
	return order, nil
}

func (s *StoreService) GetInventory() (map[string]int, error) {
	inventory, err := s.storage.GetInventory()
	if err != nil {
		return nil, ErrNoInventory
	}
	return inventory, nil
}
