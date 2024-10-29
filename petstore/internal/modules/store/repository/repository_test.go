package repository

import (
	"test/internal/models"
	"testing"
)

func NewMockStorage() *MockStorage {
	return &MockStorage{

		Create_mock: func(order *models.Order) error {
			return nil
		},
		Delelte_mock: func(id int64) error {
			return nil
		},
		GetByID_mock: func(id int64) (*models.Order, error) {
			return &models.Order{}, nil
		},
		GetInventory_mock: func() (map[string]int, error) {
			return map[string]int{}, nil
		},
	}

}

type MockStorage struct {
	Create_mock       func(order *models.Order) error
	Delelte_mock      func(id int64) error
	GetByID_mock      func(id int64) (*models.Order, error)
	GetInventory_mock func() (map[string]int, error)
}

func (m *MockStorage) Create(order *models.Order) error {
	return m.Create_mock(order)
}
func (m *MockStorage) Delete(id int64) error {
	return m.Delelte_mock(id)
}
func (m *MockStorage) GetByID(id int64) (*models.Order, error) {
	return m.GetByID_mock(id)
}
func (m *MockStorage) GetInventory() (map[string]int, error) {
	return m.GetInventory_mock()
}

func TestRepo(t *testing.T) {
	storeRepository := NewMockStorage()

	t.Run("Create", func(t *testing.T) {
		resp := storeRepository.Create(&models.Order{})
		if resp != nil {
			t.Errorf("expected error got nil")
		}
	})

	t.Run("Delete", func(t *testing.T) {
		resp := storeRepository.Delete(0)
		if resp != nil {
			t.Errorf("expected error got nil")
		}
	})

	t.Run("Get by ID", func(t *testing.T) {
		resp, _ := storeRepository.GetByID(0)
		if resp == nil {
			t.Errorf("expected error got nil")
		}
	})

	t.Run("Get inventory", func(t *testing.T) {
		resp, _ := storeRepository.GetInventory()
		if resp == nil {
			t.Errorf("expected error got nil")
		}
	})
}
