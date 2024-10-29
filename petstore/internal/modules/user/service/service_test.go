package service

import (
	"fmt"
	filter "test/internal/infrastructure/filters"
	"test/internal/models"
	"testing"
)

type MockStorage struct {
}

func testpassword() models.Password {
	return models.Password{Hash: []byte{0x01, 0x02}}

}

func (m *MockStorage) Get(id int64) (*models.User, error) {
	return &models.User{
		ID:        1,
		Name:      "test",
		Email:     "email",
		Activated: true,
		Password:  testpassword(),
	}, nil
}

func (m *MockStorage) GetByName(email string) (*models.User, error) {
	return &models.User{
		ID:        1,
		Name:      "test",
		Email:     "email",
		Activated: true,
		Password:  testpassword(),
	}, nil
}

func (m *MockStorage) GetAll(filters filter.Filters) ([]*models.User, filter.Metadata, error) {
	return []*models.User{}, filter.Metadata{}, nil
}

func (m *MockStorage) Insert(user *models.User) error {
	return nil
}

func (m *MockStorage) Update(user *models.User) error {
	return nil
}

func (m *MockStorage) Delete(id int64) error {
	return nil
}

func TestUserService(t *testing.T) {
	mockStorage := MockStorage{}
	userService := NewUserService(&mockStorage)
	t.Run("ListUsers", func(t *testing.T) {
		resp, _, _ := userService.ListUsers(filter.Filters{})
		fmt.Println(resp)
	})
	t.Run("Get users by name", func(t *testing.T) {
		resp, _ := userService.GetUserByName("")
		if resp == nil {
			t.Errorf("expected error got nil")
		}

	})

	t.Run("Get users by ID", func(t *testing.T) {
		resp, _ := userService.GetUserById(0)
		if resp == nil {
			t.Errorf("expected error got nil")
		}

	})
	t.Run("Create", func(t *testing.T) {
		password := "password"

		resp := userService.CreateUser(password, &models.User{})
		if resp != ErrValidation {
			t.Errorf("expected error got %v", resp)
		}

	})

	t.Run("Update", func(t *testing.T) {

		name := "test"
		email := "email"
		dto := struct {
			Name  *string `json:"name"`
			Email *string `json:"email"`
		}{
			Name:  &name,
			Email: &email,
		}
		resp := userService.UpdateUser(dto, *dto.Name)
		if resp == nil {
			t.Errorf("expected error got nil")
		}

	})

	t.Run("Delete", func(t *testing.T) {
		name := "test"
		resp := userService.DeleteUser(name)
		if resp != nil {
			t.Errorf("expected error got nil")
		}

	})

}
