package repository

import (
	"test/internal/models"
	"testing"
)

func NewMockStorage() *MockStorage {
	return &MockStorage{

		Create_mock: func(pet *models.Pet) error {
			return nil
		},
		Update_mock: func(pet *models.Pet) error {
			return nil
		},
		Update_put_mock: func(pet *models.Pet) error {
			return nil
		},
		Delete_mock: func(id int64) error {
			return nil
		},
		GetByID_mock: func(id int64) (*models.Pet, error) {
			return &models.Pet{}, nil
		},
		GetByStatus_mock: func(status string) ([]models.Pet, error) {
			return []models.Pet{}, nil
		},
	}

}

type MockStorage struct {
	Create_mock      func(pet *models.Pet) error
	Update_mock      func(pet *models.Pet) error
	Update_put_mock  func(pet *models.Pet) error
	Delete_mock      func(id int64) error
	GetByID_mock     func(id int64) (*models.Pet, error)
	GetByStatus_mock func(status string) ([]models.Pet, error)
}

func (m *MockStorage) Create(pet *models.Pet) error {
	return m.Create_mock(pet)
}

func (m *MockStorage) Update(pet *models.Pet) error {
	return m.Update_mock(pet)
}

func (m *MockStorage) Update_put(pet *models.Pet) error {
	return m.Update_put_mock(pet)
}

func (m *MockStorage) Delete(id int64) error {
	return m.Delete_mock(id)
}

func (m *MockStorage) GetByID(id int64) (*models.Pet, error) {
	return m.GetByID_mock(id)
}

func (m *MockStorage) GetByStatus(status string) ([]models.Pet, error) {
	return m.GetByStatus_mock(status)
}

func TestRepo(t *testing.T) {
	petRepository := NewMockStorage()

	t.Run("Create", func(t *testing.T) {
		resp := petRepository.Create(&models.Pet{})
		if resp != nil {
			t.Errorf("expected error got nil")
		}
	})

	t.Run("Update", func(t *testing.T) {
		resp := petRepository.Update(&models.Pet{})
		if resp != nil {
			t.Errorf("expected error got nil")
		}
	})

	t.Run("Delete", func(t *testing.T) {
		resp := petRepository.Delete(0)
		if resp != nil {
			t.Errorf("expected error got nil")
		}
	})

	t.Run("Get by ID", func(t *testing.T) {
		resp, _ := petRepository.GetByID(0)
		if resp == nil {
			t.Errorf("expected error got nil")
		}
	})

	t.Run("Get by status", func(t *testing.T) {
		resp, _ := petRepository.GetByStatus("")
		if resp == nil {
			t.Errorf("expected error got nil")
		}
	})
}
