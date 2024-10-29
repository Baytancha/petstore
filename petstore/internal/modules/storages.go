package modules

import (
	pet_storage "test/internal/modules/pet/repository"
	store_storage "test/internal/modules/store/repository"
	user_storage "test/internal/modules/user/repository"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type Storages struct {
	UserStorage  user_storage.IUserStorage
	PetStorage   pet_storage.IPetStorage
	StoreStorage store_storage.IStoreStorage
}

func NewStorages(sql *sqlx.DB, logger *zap.Logger) *Storages {
	return &Storages{
		UserStorage:  user_storage.NewUserModel(sql),
		PetStorage:   pet_storage.NewPetStorage(sql, logger),
		StoreStorage: store_storage.NewStoreStorage(sql, logger),
	}
}
