package service

import (
	"fmt"
	"test/internal/models"
	"testing"
)

// Create(pet *models.Pet) error
// 	Update(pet *models.Pet) error
// 	Update_put(pet *models.Pet) error
// 	Delete(id int64) error
// 	GetByID(id int64) (*models.Pet, error)
// 	GetByStatus(status string) ([]models.Pet, error)

type MockStorage struct {
	Create_mock      func(pet *models.Pet) error
	Update_mock      func(pet *models.Pet) error
	Update_put_mock  func(pet *models.Pet) error
	Delete_mock      func(id int64) error
	GetByID_mock     func(id int64) (*models.Pet, error)
	GetByStatus_mock func(status string) ([]models.Pet, error)
}

func testpetctor() *models.Pet {
	name := "test"
	return &models.Pet{
		ID:        1,
		Category:  &models.Category{ID: 1, Name: "test"},
		Name:      &name,
		PhotoUrls: make([]string, 0),
		Status:    "available",
		Tags:      make([]*models.Tag, 0),
	}
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

func TestMockService(t *testing.T) {
	mockStorage := MockStorage{}

	mockStorage.Create_mock = func(pet *models.Pet) error {
		return nil
	}

	mockStorage.GetByID_mock = func(id int64) (*models.Pet, error) {
		return testpetctor(), nil
	}
	mockStorage.GetByStatus_mock = func(status string) ([]models.Pet, error) {
		return []models.Pet{}, nil
	}
	mockStorage.Delete_mock = func(id int64) error {
		return nil
	}
	mockStorage.Update_mock = func(pet *models.Pet) error {
		return nil
	}
	mockStorage.Update_put_mock = func(pet *models.Pet) error {
		return nil
	}

	storeService := NewPetService(&mockStorage)
	t.Run("Create", func(t *testing.T) {
		resp := storeService.Create(&models.Pet{})
		fmt.Println(resp)
	})
	t.Run("Get by ID", func(t *testing.T) {
		resp, _ := storeService.GetByID(0)
		if resp == nil {
			t.Errorf("expected error got nil")
		}

	})

	t.Run("Get by Status", func(t *testing.T) {
		resp, _ := storeService.GetByStatus("")
		if resp == nil {
			t.Errorf("expected error got nil")
		}

	})

	t.Run("Update", func(t *testing.T) {
		name := "test"
		status := "test"
		ID := 0
		resp := storeService.Update(&name, &status, int64(ID))
		if resp != nil {
			t.Errorf("expected error got nil")
		}

	})

	t.Run("Update_put", func(t *testing.T) {
		resp := storeService.Update_put(testpetctor())
		if resp != nil {
			t.Errorf("expected error got nil")
		}

	})

	t.Run("Delete", func(t *testing.T) {
		resp := storeService.Delete(0)
		if resp != nil {
			t.Errorf("expected error got nil")
		}

	})

}
