package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"test/internal/infrastructure/responder"
	"test/internal/models"
	"test/internal/modules/pet/service"

	"github.com/go-chi/chi"
)

// r.Post("/pet", ctrl.petController.PetCreate)
// r.Post("/pet/{petID}", ctrl.petController.PetUpdate)
// r.Get("/pet/{petID}", ctrl.petController.PetGetByID)
// r.Get("/pet/findByStatus", ctrl.petController.PetGetByStatus)
// r.Put("/pet/{petID}", ctrl.petController.PetUpdate)
// r.Delete("/pet/{petID}", ctrl.petController.PetDelete)

type IPetController interface {
	PetCreate(w http.ResponseWriter, r *http.Request)
	PetUpdate(w http.ResponseWriter, r *http.Request)
	PetUpdate_post(w http.ResponseWriter, r *http.Request)
	PetGetByID(w http.ResponseWriter, r *http.Request)
	PetGetByStatus(w http.ResponseWriter, r *http.Request)
	PetDelete(w http.ResponseWriter, r *http.Request)
}

type PetController struct {
	responder responder.Responder
	service   service.IPetstoreService
}

func NewPetController(responder responder.Responder, service service.IPetstoreService) *PetController {
	return &PetController{
		responder: responder,
		service:   service,
	}
}

func (p *PetController) PetCreate(w http.ResponseWriter, r *http.Request) {
	var pet *models.Pet
	err := json.NewDecoder(r.Body).Decode(&pet)

	if err != nil {
		p.responder.ErrorBadRequest(w, err)
		return
	}

	err = p.service.Create(pet)
	if err != nil {
		p.responder.ErrorInternal(w, err)
		return
	}

	p.responder.OutputJSON(w, pet)
}

func (p *PetController) PetGetByID(w http.ResponseWriter, r *http.Request) {
	var (
		err      error
		petIDRaw string
		petID    int
	)

	petIDRaw = chi.URLParam(r, "petID")

	petID, err = strconv.Atoi(petIDRaw)
	if err != nil {
		p.responder.ErrorBadRequest(w, err)
		return
	}

	pet, err := p.service.GetByID(int64(petID))
	if err != nil {
		p.responder.ErrorInternal(w, err)
		return
	}

	p.responder.OutputJSON(w, pet)
}

func (p *PetController) PetGetByStatus(w http.ResponseWriter, r *http.Request) {

	var (
		pets      []models.Pet
		petStatus string
	)

	petStatus = r.URL.Query().Get("status")

	if petStatus != "available" && petStatus != "pending" && petStatus != "sold" {
		p.responder.ErrorBadRequest(w, errors.New("Invalid request body"))
		return
	}

	pets, err := p.service.GetByStatus(petStatus)
	if err != nil {
		p.responder.ErrorInternal(w, err)
		return
	}

	p.responder.OutputJSON(w, pets)

}

func (p *PetController) PetUpdate(w http.ResponseWriter, r *http.Request) {
	var pet *models.Pet
	err := json.NewDecoder(r.Body).Decode(&pet)

	if err != nil {
		p.responder.ErrorBadRequest(w, err)
		return
	}

	err = p.service.Update_put(pet)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrEditConflict):
			p.responder.ErrorInternal(w, errors.New("Error while updating user"))
		case errors.Is(err, service.ErrRecordNotFound):
			p.responder.ErrorInternal(w, errors.New("User not found"))
		}
		return
	}
	p.responder.OutputJSON(w, pet)
}

func (p *PetController) PetUpdate_post(w http.ResponseWriter, r *http.Request) {

	petID := chi.URLParam(r, "petID")

	err := r.ParseForm()
	if err != nil {
		fmt.Println("1")
		p.responder.ErrorBadRequest(w, err)
		return
	}

	name := r.FormValue("name")
	if name == "" {
		fmt.Println("2")
		p.responder.ErrorBadRequest(w, err)
		return
	}
	status := r.FormValue("status")
	if status == "" {
		fmt.Println("3")
		p.responder.ErrorBadRequest(w, err)
		return
	}
	ID, err := strconv.Atoi(petID)
	if err != nil {
		p.responder.ErrorBadRequest(w, err)
		return
	}

	err = p.service.Update(&name, &status, int64(ID))
	if err != nil {
		switch {
		case errors.Is(err, service.ErrEditConflict):
			p.responder.ErrorInternal(w, errors.New("Error while updating user"))
		case errors.Is(err, service.ErrRecordNotFound):
			p.responder.ErrorInternal(w, errors.New("User not found"))
		}
		return
	}
	p.responder.OutputJSON(w, "Pet updated successfully")

}

func (p *PetController) PetDelete(w http.ResponseWriter, r *http.Request) {
	petID := chi.URLParam(r, "petID")
	if petID == "" {
		p.responder.ErrorBadRequest(w, errors.New("Invalid request body"))
		return
	}
	ID, err := strconv.Atoi(petID)
	if err != nil {
		p.responder.ErrorBadRequest(w, err)
		return
	}
	err = p.service.Delete(int64(ID))
	if err != nil {
		p.responder.ErrorInternal(w, err)
		return
	}
	p.responder.OutputJSON(w, "Pet deleted successfully")
}
