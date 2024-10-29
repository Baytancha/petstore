package repository

import (
	// New import
	// New import

	//"fmt"

	"errors"

	"test/internal/models"
)

var (
	ErrPetNotFound       = errors.New("pet not found")
	ErrEditConflict      = errors.New("edit conflict")
	ErrDuplicateCategory = errors.New("duplicate pet category")
	ErrDuplicateTag      = errors.New("duplicate pet tag")
)

type IPetStorage interface {
	Create(pet *models.Pet) error
	Update(pet *models.Pet) error
	Update_put(pet *models.Pet) error
	Delete(id int64) error
	GetByID(id int64) (*models.Pet, error)
	GetByStatus(status string) ([]models.Pet, error)
}
