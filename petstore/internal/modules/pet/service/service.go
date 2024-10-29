package service

import (
	"errors"
	"test/internal/modules/pet/repository"

	"test/internal/models"
)

var (
	ErrRecordNotFound  = errors.New("record not found")
	ErrEditConflict    = errors.New("edit conflict")
	ErrDuplicateRecord = errors.New("duplicate record")
)

// r.Post("/pet", ctrl.petController.PetCreate)
// r.Post("/pet/{petID}", ctrl.petController.PetUpdate)
// r.Get("/pet/{petID}", ctrl.petController.PetGetByID)
// r.Get("/pet/findByStatus", ctrl.petController.PetGetByStatus)
// r.Put("/pet/{petID}", ctrl.petController.PetUpdate)
// r.Delete("/pet/{petID}", ctrl.petController.PetDelete)

type IPetstoreService interface {
	Create(pet *models.Pet) error
	Update(name, status *string, ID int64) error
	Update_put(pet *models.Pet) error
	Delete(id int64) error
	GetByID(id int64) (*models.Pet, error)
	GetByStatus(status string) ([]models.Pet, error)
}

type PetService struct {
	storage repository.IPetStorage
}

func NewPetService(repo repository.IPetStorage) *PetService {
	return &PetService{storage: repo}
}

func (s *PetService) Create(pet *models.Pet) error {
	err := s.storage.Create(pet)
	if err != nil {
		return ErrDuplicateRecord
	}
	return nil
}

func (s *PetService) Update(name, status *string, ID int64) error {

	updated, err := s.storage.GetByID(ID)
	if err != nil {
		return ErrRecordNotFound
	}
	updated.Name = name
	updated.Status = *status
	err = s.storage.Update(updated)
	if err != nil {
		return ErrEditConflict
	}
	return nil
}

func (s *PetService) Update_put(pet *models.Pet) error {

	updated, err := s.storage.GetByID(int64(pet.ID))
	if err != nil {
		return ErrRecordNotFound
	}
	updated.Name = pet.Name
	updated.Category.ID = pet.Category.ID
	updated.Category.Name = pet.Category.Name
	updated.Status = pet.Status
	updated.PhotoUrls = pet.PhotoUrls
	updated.Tags = pet.Tags

	err = s.storage.Update_put(updated)
	if err != nil {
		return ErrEditConflict
	}
	return nil
}
func (s *PetService) Delete(id int64) error {
	err := s.storage.Delete(id)
	if err != nil {
		return ErrRecordNotFound
	}
	return nil
}

func (s *PetService) GetByID(id int64) (*models.Pet, error) {
	pet, err := s.storage.GetByID(id)
	if err != nil {
		return nil, ErrRecordNotFound
	}
	return pet, nil
}

func (s *PetService) GetByStatus(status string) ([]models.Pet, error) {
	pet, err := s.storage.GetByStatus(status)
	if err != nil {
		return nil, ErrRecordNotFound
	}
	return pet, nil
}
