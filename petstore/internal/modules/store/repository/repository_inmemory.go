package repository

import (
	"sync"

	"go.uber.org/zap"

	"test/internal/models"
)

type StoreStorage_map struct {
	pets               map[int64]*models.Pet
	logger             *zap.Logger
	orders             []*models.Order
	primaryKeyIDx      map[int64]*models.Order
	autoIncrementCount int
	sync.Mutex
}

func NewStoreStorage_map(logger *zap.Logger) IStoreStorage {
	return &StoreStorage_map{
		pets:               models.PrimaryKeyIDx,
		logger:             logger,
		primaryKeyIDx:      make(map[int64]*models.Order),
		autoIncrementCount: 0,
		orders:             make([]*models.Order, 0),
	}
}

func (ps *StoreStorage_map) Create(order *models.Order) error {
	ps.Lock()
	defer ps.Unlock()
	v, ok := ps.pets[order.PetID]
	if !ok {
		return ErrRecordNotFound
	}
	v.Status = "pending"
	order.ID = int64(ps.autoIncrementCount)
	ps.autoIncrementCount++
	ps.primaryKeyIDx[order.ID] = order
	ps.orders = append(ps.orders, order)
	return nil
}

func (ps *StoreStorage_map) Delete(id int64) error {
	if _, ok := ps.primaryKeyIDx[id]; ok {
		delete(ps.primaryKeyIDx, id)
	}

	for i := 0; i < len(ps.orders); i++ {
		if ps.orders[i].ID == id {
			ps.orders = append(ps.orders[:i], ps.orders[i+1:]...)
			return nil
		}
	}

	return ErrEditConflict
}

func (ps *StoreStorage_map) GetByID(id int64) (*models.Order, error) {
	ps.Lock()
	defer ps.Unlock()

	if v, ok := ps.primaryKeyIDx[id]; ok {
		return v, nil
	}

	for i := 0; i < len(ps.orders); i++ {
		if ps.orders[i].ID == id {
			return ps.orders[i], nil
		}
	}

	return nil, ErrRecordNotFound
}

func (ps *StoreStorage_map) GetInventory() (map[string]int, error) {
	inventory := make(map[string]int)
	for _, v := range ps.pets {
		inventory[v.Status]++
	}

	return inventory, nil
}
