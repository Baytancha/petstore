package modules

import (
	"test/internal/infrastructure/components"
	pet_controller "test/internal/modules/pet/controllers"
	store_controller "test/internal/modules/store/controller"
	user_controller "test/internal/modules/user/controller"
)

type Controllers struct {
	UserHandler  user_controller.IUserHandler
	PetHandler   pet_controller.IPetController
	StoreHandler store_controller.IStoreController
}

func NewControllers(services *Services, components *components.Components) *Controllers {
	return &Controllers{
		UserHandler:  user_controller.NewUserHandler(components.Responder, services.UserService),
		PetHandler:   pet_controller.NewPetController(components.Responder, services.PetService),
		StoreHandler: store_controller.NewStoreController(components.Responder, services.StoreService),
	}
}
