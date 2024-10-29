package controller

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"test/internal/infrastructure/responder"
	"test/internal/models"
	"test/internal/modules/store/service"
	"testing"

	"github.com/go-chi/chi"
	jsoniter "github.com/json-iterator/go"
	"github.com/ptflp/godecoder"
	"go.uber.org/zap"
)

// r.Get("/store/inventory", ctrl.StoreHandler.GetInventory)

// r.Get("/store/order/{orderID}", ctrl.StoreHandler.GetOrderByID)
// r.Post("/store/order", ctrl.StoreHandler.CreateOrder)
// r.Delete("/store/order/{orderID}", ctrl.StoreHandler.DeleteOrder)

type MockStorage struct {
	Create_mock       func(order *models.Order) error
	Delelte_mock      func(id int64) error
	GetByID_mock      func(id int64) (*models.Order, error)
	GetInventory_mock func() (map[string]int, error)
}

func (m *MockStorage) Create(order *models.Order) error {
	return m.Create_mock(order)
}
func (m *MockStorage) Delete(id int64) error {
	return m.Delelte_mock(id)
}
func (m *MockStorage) GetByID(id int64) (*models.Order, error) {
	return m.GetByID_mock(id)
}
func (m *MockStorage) GetInventory() (map[string]int, error) {
	return m.GetInventory_mock()
}

func TestCreateHandler(t *testing.T) {

	t.Run("happy path", func(t *testing.T) {

		req := httptest.NewRequest("POST", "/store/order", bytes.NewReader([]byte(`{
  "id": 0,
  "petId": 0,
  "quantity": 0,
  "shipDate": "2024-10-12T18:41:14.730Z",
  "status": "placed",
  "complete": true
}`)))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		mock := &MockStorage{
			Create_mock: func(order *models.Order) error { return nil },
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

		service := service.NewStoreService(mock)

		controller := NewStoreController(responder.NewResponder(decoder, logger), service)
		controller.CreateOrder(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status code %d but got %d", http.StatusOK, w.Code)
		}

	})

	t.Run("dublicate", func(t *testing.T) {

		req := httptest.NewRequest("POST", "/store/order", bytes.NewReader([]byte(`{
  "id": 0,
  "'petkhjkhId": 0,
}`)))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		mock := &MockStorage{
			Create_mock: func(order *models.Order) error { return errors.New("some error") },
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

		service := service.NewStoreService(mock)

		controller := NewStoreController(responder.NewResponder(decoder, logger), service)
		controller.CreateOrder(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status code %d but got %d", http.StatusBadRequest, w.Code)
		}

	})

}

func TestGetOrderByID(t *testing.T) {

	t.Run("happy path", func(t *testing.T) {

		req := httptest.NewRequest("GET", "/store/order/1", nil)
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		mock := &MockStorage{
			GetByID_mock: func(id int64) (*models.Order, error) { return &models.Order{}, nil },
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

		service := service.NewStoreService(mock)

		controller := NewStoreController(responder.NewResponder(decoder, logger), service)
		chiCtx := chi.NewRouteContext()

		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))

		chiCtx.URLParams.Add("orderID", "1")
		controller.GetOrderByID(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status code %d but got %d", http.StatusOK, w.Code)
		}

	})
}

func TestDeleteOrder(t *testing.T) {

	t.Run("happy path", func(t *testing.T) {

		req := httptest.NewRequest("DELETE", "/store/order/1", nil)
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		mock := &MockStorage{
			Delelte_mock: func(id int64) error { return nil },
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

		service := service.NewStoreService(mock)

		controller := NewStoreController(responder.NewResponder(decoder, logger), service)
		chiCtx := chi.NewRouteContext()

		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))

		chiCtx.URLParams.Add("orderID", "1")
		controller.DeleteOrder(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status code %d but got %d", http.StatusOK, w.Code)
		}

	})
}

func TestGetAll(t *testing.T) {

	t.Run("happy path", func(t *testing.T) {

		req := httptest.NewRequest("GET", "/store/inventory", nil)
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		mock := &MockStorage{
			GetInventory_mock: func() (map[string]int, error) {
				return map[string]int{}, nil
			},
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

		service := service.NewStoreService(mock)

		controller := NewStoreController(responder.NewResponder(decoder, logger), service)
		chiCtx := chi.NewRouteContext()

		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))

		chiCtx.URLParams.Add("id", "1")
		controller.GetInventory(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status code %d but got %d", http.StatusOK, w.Code)
		}

	})

	t.Run("internal server error", func(t *testing.T) {

		req := httptest.NewRequest("GET", "/store/inventory", nil)
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		mock := &MockStorage{
			GetInventory_mock: func() (map[string]int, error) {
				return nil, errors.New("failed to get inventory")
			},
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

		service := service.NewStoreService(mock)

		controller := NewStoreController(responder.NewResponder(decoder, logger), service)
		chiCtx := chi.NewRouteContext()

		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))

		chiCtx.URLParams.Add("id", "1")
		controller.GetInventory(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("expected status code %d but got %d", http.StatusInternalServerError, w.Code)
		}

	})
}
