package controller

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"test/internal/infrastructure/responder"
	"test/internal/models"
	"test/internal/modules/pet/service"
	"testing"

	"github.com/go-chi/chi"
	jsoniter "github.com/json-iterator/go"
	"github.com/ptflp/godecoder"
	"go.uber.org/zap"
)

// r.Post("/pet", ctrl.petController.PetCreate)
// r.Post("/pet/{petID}", ctrl.petController.PetUpdate)
// r.Get("/pet/{petID}", ctrl.petController.PetGetByID)
// r.Get("/pet/findByStatus", ctrl.petController.PetGetByStatus)
// r.Put("/pet/{petID}", ctrl.petController.PetUpdate)
// r.Delete("/pet/{petID}", ctrl.petController.PetDelete)

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

func TestCreatePetHandler(t *testing.T) {

	t.Run("happy path", func(t *testing.T) {

		req := httptest.NewRequest("POST", "/pet", bytes.NewReader([]byte(`{  "id": 123,
  "category": {
    "id": 0,
    "name": "STORMY"
  },
  "name": "CAT",
  "photoUrls": [
    "string",
    "PIC1",
    "PIC2"
  ],
  "tags": [
    {
      "id": 0,
      "name": "PURRING"
    }
  ],
  "status": "available"
}`)))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		mock := &MockStorage{
			Create_mock: func(pet *models.Pet) error { return nil },
		}

		logger, err := zap.NewProduction()
		if err != nil {
			panic(err)
		}
		decoder := godecoder.NewDecoder(jsoniter.Config{
			EscapeHTML:             true,
			SortMapKeys:            true,
			ValidateJsonRawMessage: true,
			DisallowUnknownFields:  true,
		})

		service := service.NewPetService(mock)

		controller := NewPetController(responder.NewResponder(decoder, logger), service)
		controller.PetCreate(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status code %d but got %d", http.StatusOK, w.Code)
		}

	})

	t.Run("invalid json", func(t *testing.T) {

		req := httptest.NewRequest("POST", "/pet", bytes.NewReader([]byte(`{  "id": 123,
  "c
      "nme": "PURRING"
    }
  ],
  "status": "available"
}`)))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		mock := &MockStorage{
			Create_mock: func(pet *models.Pet) error { return nil },
		}

		logger, err := zap.NewProduction()
		if err != nil {
			panic(err)
		}
		decoder := godecoder.NewDecoder(jsoniter.Config{
			EscapeHTML:             true,
			SortMapKeys:            true,
			ValidateJsonRawMessage: true,
			DisallowUnknownFields:  true,
		})

		service := service.NewPetService(mock)

		controller := NewPetController(responder.NewResponder(decoder, logger), service)
		controller.PetCreate(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status code %d but got %d", http.StatusBadRequest, w.Code)
		}

	})

	t.Run("internal server error", func(t *testing.T) {

		req := httptest.NewRequest("POST", "/pet", bytes.NewReader([]byte(`{  "id": 123,
  "category": {
    "id": 0,
    "name": "STORMY"
  },
  "name": "CAT",
  "photoUrls": [
    "string",
    "PIC1",
    "PIC2"
  ],
  "tags": [
    {
      "id": 0,
      "name": "PURRING"
    }
  ],
  "status": "available"
}`)))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		mock := &MockStorage{
			Create_mock: func(pet *models.Pet) error { return errors.New("some error") },
		}

		logger, err := zap.NewProduction()
		if err != nil {
			panic(err)
		}
		decoder := godecoder.NewDecoder(jsoniter.Config{
			EscapeHTML:             true,
			SortMapKeys:            true,
			ValidateJsonRawMessage: true,
			DisallowUnknownFields:  true,
		})

		service := service.NewPetService(mock)

		controller := NewPetController(responder.NewResponder(decoder, logger), service)
		controller.PetCreate(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("expected status code %d but got %d", http.StatusInternalServerError, w.Code)
		}

	})

}

func TestGetByIDPetHandler(t *testing.T) {

	t.Run("happy path", func(t *testing.T) {

		req := httptest.NewRequest("GET", "/pet/1", nil)
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		mock := &MockStorage{
			GetByID_mock: func(id int64) (*models.Pet, error) { return &models.Pet{}, nil },
		}

		logger, err := zap.NewProduction()
		if err != nil {
			panic(err)
		}
		decoder := godecoder.NewDecoder(jsoniter.Config{
			EscapeHTML:             true,
			SortMapKeys:            true,
			ValidateJsonRawMessage: true,
			DisallowUnknownFields:  true,
		})

		service := service.NewPetService(mock)
		chiCtx := chi.NewRouteContext()

		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))

		chiCtx.URLParams.Add("petID", "1")

		controller := NewPetController(responder.NewResponder(decoder, logger), service)
		controller.PetGetByID(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status code %d but got %d", http.StatusOK, w.Code)
		}

	})

	t.Run("invalid  request", func(t *testing.T) {

		req := httptest.NewRequest("GET", "/pet/", nil)
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		mock := &MockStorage{
			GetByID_mock: func(id int64) (*models.Pet, error) { return &models.Pet{}, nil },
		}

		logger, err := zap.NewProduction()
		if err != nil {
			panic(err)
		}
		decoder := godecoder.NewDecoder(jsoniter.Config{
			EscapeHTML:             true,
			SortMapKeys:            true,
			ValidateJsonRawMessage: true,
			DisallowUnknownFields:  true,
		})

		service := service.NewPetService(mock)
		chiCtx := chi.NewRouteContext()

		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))

		chiCtx.URLParams.Add("id", "OIUUOU")

		controller := NewPetController(responder.NewResponder(decoder, logger), service)
		controller.PetGetByID(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status code %d but got %d", http.StatusBadRequest, w.Code)
		}

	})

	t.Run("internal server error", func(t *testing.T) {

		req := httptest.NewRequest("GET", "/pet/", nil)
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		mock := &MockStorage{
			GetByID_mock: func(id int64) (*models.Pet, error) { return nil, errors.New("some error") },
		}

		logger, err := zap.NewProduction()
		if err != nil {
			panic(err)
		}
		decoder := godecoder.NewDecoder(jsoniter.Config{
			EscapeHTML:             true,
			SortMapKeys:            true,
			ValidateJsonRawMessage: true,
			DisallowUnknownFields:  true,
		})

		service := service.NewPetService(mock)
		chiCtx := chi.NewRouteContext()

		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))

		chiCtx.URLParams.Add("petID", "1")

		controller := NewPetController(responder.NewResponder(decoder, logger), service)
		controller.PetGetByID(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("expected status code %d but got %d", http.StatusInternalServerError, w.Code)
		}

	})

}

func TestGetByStatusPetHandler(t *testing.T) {

	t.Run("happy path", func(t *testing.T) {

		req := httptest.NewRequest("GET", "/pet/", nil)
		req.Header.Set("Content-Type", "application/json")
		req.URL.RawQuery = "status=available"

		w := httptest.NewRecorder()

		mock := &MockStorage{
			GetByStatus_mock: func(status string) ([]models.Pet, error) { return []models.Pet{}, nil },
		}

		logger, err := zap.NewProduction()
		if err != nil {
			panic(err)
		}
		decoder := godecoder.NewDecoder(jsoniter.Config{
			EscapeHTML:             true,
			SortMapKeys:            true,
			ValidateJsonRawMessage: true,
			DisallowUnknownFields:  true,
		})

		service := service.NewPetService(mock)

		controller := NewPetController(responder.NewResponder(decoder, logger), service)
		controller.PetGetByStatus(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status code %d but got %d", http.StatusOK, w.Code)
		}

	})

	t.Run("invalid request", func(t *testing.T) {

		req := httptest.NewRequest("GET", "/pet/findByStatus", nil)
		req.Header.Set("Content-Type", "application/json")
		req.URL.RawQuery = "satus=avable"

		w := httptest.NewRecorder()

		mock := &MockStorage{
			GetByStatus_mock: func(status string) ([]models.Pet, error) { return []models.Pet{}, nil },
		}

		logger, err := zap.NewProduction()
		if err != nil {
			panic(err)
		}
		decoder := godecoder.NewDecoder(jsoniter.Config{
			EscapeHTML:             true,
			SortMapKeys:            true,
			ValidateJsonRawMessage: true,
			DisallowUnknownFields:  true,
		})

		service := service.NewPetService(mock)

		controller := NewPetController(responder.NewResponder(decoder, logger), service)
		controller.PetGetByStatus(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status code %d but got %d", http.StatusBadRequest, w.Code)
		}

	})

	t.Run("internal server error", func(t *testing.T) {

		req := httptest.NewRequest("GET", "/pet/findByStatus", nil)
		req.Header.Set("Content-Type", "application/json")
		req.URL.RawQuery = "status=available"

		w := httptest.NewRecorder()

		mock := &MockStorage{
			GetByStatus_mock: func(status string) ([]models.Pet, error) { return nil, errors.New("some error") },
		}

		logger, err := zap.NewProduction()
		if err != nil {
			panic(err)
		}
		decoder := godecoder.NewDecoder(jsoniter.Config{
			EscapeHTML:             true,
			SortMapKeys:            true,
			ValidateJsonRawMessage: true,
			DisallowUnknownFields:  true,
		})

		service := service.NewPetService(mock)

		controller := NewPetController(responder.NewResponder(decoder, logger), service)
		controller.PetGetByStatus(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("expected status code %d but got %d", http.StatusInternalServerError, w.Code)
		}

	})

}

func TestPetDeleteHandler(t *testing.T) {

	t.Run("happy path", func(t *testing.T) {

		req := httptest.NewRequest("DELETE", "/pet/", nil)
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		mock := &MockStorage{
			Delete_mock: func(id int64) error { return nil },
		}

		logger, err := zap.NewProduction()
		if err != nil {
			panic(err)
		}
		decoder := godecoder.NewDecoder(jsoniter.Config{
			EscapeHTML:             true,
			SortMapKeys:            true,
			ValidateJsonRawMessage: true,
			DisallowUnknownFields:  true,
		})

		service := service.NewPetService(mock)
		chiCtx := chi.NewRouteContext()

		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))

		chiCtx.URLParams.Add("petID", "1")
		controller := NewPetController(responder.NewResponder(decoder, logger), service)
		controller.PetDelete(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status code %d but got %d", http.StatusOK, w.Code)
		}

	})

	t.Run("bad request", func(t *testing.T) {

		req := httptest.NewRequest("DELETE", "/pet/", nil)
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		mock := &MockStorage{
			Delete_mock: func(id int64) error { return nil },
		}

		logger, err := zap.NewProduction()
		if err != nil {
			panic(err)
		}
		decoder := godecoder.NewDecoder(jsoniter.Config{
			EscapeHTML:             true,
			SortMapKeys:            true,
			ValidateJsonRawMessage: true,
			DisallowUnknownFields:  true,
		})

		service := service.NewPetService(mock)
		chiCtx := chi.NewRouteContext()

		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))

		chiCtx.URLParams.Add("id", "GDGD")
		controller := NewPetController(responder.NewResponder(decoder, logger), service)
		controller.PetDelete(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status code %d but got %d", http.StatusBadRequest, w.Code)
		}

	})

	t.Run("INTERNAL SERVER ERROR", func(t *testing.T) {

		req := httptest.NewRequest("DELETE", "/pet/", nil)
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		mock := &MockStorage{
			Delete_mock: func(id int64) error { return errors.New("some error") },
		}

		logger, err := zap.NewProduction()
		if err != nil {
			panic(err)
		}
		decoder := godecoder.NewDecoder(jsoniter.Config{
			EscapeHTML:             true,
			SortMapKeys:            true,
			ValidateJsonRawMessage: true,
			DisallowUnknownFields:  true,
		})

		service := service.NewPetService(mock)
		chiCtx := chi.NewRouteContext()

		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))

		chiCtx.URLParams.Add("petID", "1")
		controller := NewPetController(responder.NewResponder(decoder, logger), service)
		controller.PetDelete(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("expected status code %d but got %d", http.StatusInternalServerError, w.Code)
		}

	})

}
