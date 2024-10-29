package modules

import (
	"test/internal/infrastructure/components"
	pet_service "test/internal/modules/pet/service"
	store_service "test/internal/modules/store/service"
	user_service "test/internal/modules/user/service"
)

type Services struct {
	UserService  user_service.IUserService
	PetService   pet_service.IPetstoreService
	StoreService store_service.IStoreService
}

func NewServices(cmp *components.Components, storages *Storages) *Services {
	return &Services{
		UserService:  user_service.NewUserService(storages.UserStorage),
		PetService:   pet_service.NewPetService(storages.PetStorage),
		StoreService: store_service.NewStoreService(storages.StoreStorage),
	}
}
