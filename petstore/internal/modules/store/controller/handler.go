package controller

import (
	"encoding/json"
	"net/http"
	"strconv"
	"test/internal/infrastructure/responder"
	"test/internal/models"
	"test/internal/modules/store/service"

	"github.com/go-chi/chi"
)

// r.Post("/pet", ctrl.petController.PetCreate)
// r.Post("/pet/{petID}", ctrl.petController.PetUpdate)
// r.Get("/pet/{petID}", ctrl.petController.PetGetByID)
// r.Get("/pet/findByStatus", ctrl.petController.PetGetByStatus)
// r.Put("/pet/{petID}", ctrl.petController.PetUpdate)
// r.Delete("/pet/{petID}", ctrl.petController.PetDelete)

type IStoreController interface {
	GetInventory(w http.ResponseWriter, r *http.Request)
	GetOrderByID(w http.ResponseWriter, r *http.Request)
	CreateOrder(w http.ResponseWriter, r *http.Request)
	DeleteOrder(w http.ResponseWriter, r *http.Request)
}

type StoreController struct {
	responder responder.Responder
	service   service.IStoreService
}

func NewStoreController(responder responder.Responder, service service.IStoreService) *StoreController {
	return &StoreController{
		responder: responder,
		service:   service,
	}
}

func (s *StoreController) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var order *models.Order
	err := json.NewDecoder(r.Body).Decode(&order)

	if err != nil {
		s.responder.ErrorBadRequest(w, err)
		return
	}

	err = s.service.Create(order)
	if err != nil {
		s.responder.ErrorInternal(w, err)
		return
	}

	s.responder.OutputJSON(w, order)
}

func (s *StoreController) GetOrderByID(w http.ResponseWriter, r *http.Request) {
	var (
		err        error
		orderIDRaw string
		orderID    int
	)

	orderIDRaw = chi.URLParam(r, "orderID")

	orderID, err = strconv.Atoi(orderIDRaw)
	if err != nil {
		s.responder.ErrorBadRequest(w, err)
		return
	}

	pet, err := s.service.GetByID(int64(orderID))
	if err != nil {
		s.responder.ErrorInternal(w, err)
		return
	}

	s.responder.OutputJSON(w, pet)
}

func (s *StoreController) DeleteOrder(w http.ResponseWriter, r *http.Request) {

	var (
		err        error
		orderIDRaw string
		orderID    int
	)

	orderIDRaw = chi.URLParam(r, "orderID")

	orderID, err = strconv.Atoi(orderIDRaw)
	if err != nil {
		s.responder.ErrorBadRequest(w, err)
		return
	}

	err = s.service.Delete(int64(orderID))
	if err != nil {
		s.responder.ErrorInternal(w, err)
		return
	}
	s.responder.OutputJSON(w, "Order deleted successfully")

}

func (s *StoreController) GetInventory(w http.ResponseWriter, r *http.Request) {

	inventory, err := s.service.GetInventory()
	if err != nil {
		s.responder.ErrorInternal(w, err)
		return
	}

	s.responder.OutputJSON(w, inventory)
}
