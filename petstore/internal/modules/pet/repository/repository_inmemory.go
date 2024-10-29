package repository

import (
	"sync"
	"test/internal/models"

	"go.uber.org/zap"
)

type PetStorage_map struct {
	logger             *zap.Logger
	data               []*models.Pet
	primaryKeyIDx      map[int64]*models.Pet
	autoIncrementCount int
	sync.Mutex
}

func NewPetStorage_map(logger *zap.Logger) IPetStorage {
	models.Pets = make([]*models.Pet, 0, 10000)
	models.PrimaryKeyIDx = make(map[int64]*models.Pet, 10000)

	return &PetStorage_map{
		logger:             logger,
		primaryKeyIDx:      models.PrimaryKeyIDx,
		autoIncrementCount: 0,
		data:               models.Pets,
	}
}

func (ps *PetStorage_map) Create(pet *models.Pet) error {
	ps.Lock()
	defer ps.Unlock()
	pet.ID = int64(ps.autoIncrementCount)
	ps.primaryKeyIDx[pet.ID] = pet
	ps.autoIncrementCount++
	ps.data = append(ps.data, pet)
	return nil

}

func (ps *PetStorage_map) Update(pet *models.Pet) error {
	ps.Lock()
	defer ps.Unlock()

	if v, ok := ps.primaryKeyIDx[pet.ID]; ok {
		*v = *pet
		return nil
	}

	for i := 0; i < len(ps.data); i++ {
		if ps.data[i].ID == pet.ID {
			ps.data[i] = pet
			return nil
		}
	}

	return ErrEditConflict
}

func (ps *PetStorage_map) Update_put(pet *models.Pet) error {
	ps.Lock()
	defer ps.Unlock()

	if v, ok := ps.primaryKeyIDx[pet.ID]; ok {
		*v = *pet
		return nil
	}

	for i := 0; i < len(ps.data); i++ {
		if ps.data[i].ID == pet.ID {
			ps.data[i] = pet
			return nil
		}
	}

	return ErrEditConflict
}

func (ps *PetStorage_map) Delete(id int64) error {
	ps.Lock()
	defer ps.Unlock()

	if _, ok := ps.primaryKeyIDx[id]; ok {
		delete(ps.primaryKeyIDx, id)
	}

	for i := 0; i < len(ps.data); i++ {
		if ps.data[i].ID == id {
			ps.data = append(ps.data[:i], ps.data[i+1:]...)
			return nil
		}
	}

	return ErrEditConflict

}

func (ps *PetStorage_map) GetByID(id int64) (*models.Pet, error) {
	ps.Lock()
	defer ps.Unlock()

	if v, ok := ps.primaryKeyIDx[id]; ok {
		return v, nil
	}

	for i := 0; i < len(ps.data); i++ {
		if ps.data[i].ID == id {
			return ps.data[i], nil
		}
	}

	return nil, ErrPetNotFound

}

func (ps *PetStorage_map) GetByStatus(status string) ([]models.Pet, error) {
	ps.Lock()
	defer ps.Unlock()

	var pets []models.Pet
	for _, pet := range ps.data {
		if pet.Status == status {
			pets = append(pets, *pet)
		}
	}
	return pets, nil

}
