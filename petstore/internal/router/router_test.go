package router

import (
	"net/http"
	"net/http/httptest"
	"test/internal/infrastructure/components"
	"test/internal/infrastructure/responder"
	"test/internal/modules"
	"testing"

	jsoniter "github.com/json-iterator/go"
	"github.com/ptflp/godecoder"
	"go.uber.org/zap"
)

func TestRouter(t *testing.T) {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	decoder := godecoder.NewDecoder(jsoniter.Config{})
	responseManager := responder.NewResponder(decoder, logger)

	components := components.NewComponents(responseManager, decoder, logger, nil)
	storages := modules.NewStorages(nil, nil)
	services := modules.NewServices(components, storages)
	ctrl := modules.NewControllers(services, components)
	r := Routes(ctrl, components)
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/user/list", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("expected status code %d, got %d", http.StatusForbidden, w.Code)
	}
}
