package handlers

import (
	"net/http"

	api "github.com/AmadoMuerte/BirthdayWish/API/apps/gateway/internal/gen"
	"github.com/AmadoMuerte/BirthdayWish/API/apps/gateway/internal/handlers/test_handler"
)

type APIImplementation struct {
	testHandler *test_handler.TestHandler
}

func NewAPIImplementation() *APIImplementation {
	return &APIImplementation{
		testHandler: test_handler.NewTestHandler(),
	}
}

func (a *APIImplementation) GetHealth(w http.ResponseWriter, r *http.Request) {
	a.testHandler.GetHealth(w, r)
}

var _ api.ServerInterface = (*APIImplementation)(nil)
