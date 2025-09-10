// apps/gateway/internal/handlers/api_impl.go
package handlers

import (
	"log/slog"
	"net/http"

	"github.com/AmadoMuerte/BirthdayWish/API/apps/gateway/internal/client"
	api "github.com/AmadoMuerte/BirthdayWish/API/apps/gateway/internal/gen"
	"github.com/AmadoMuerte/BirthdayWish/API/apps/gateway/internal/handlers/auth_handler"
	"github.com/AmadoMuerte/BirthdayWish/API/apps/gateway/internal/handlers/wish_handler"
	"github.com/go-chi/jwtauth/v5"
)

type APIImplementation struct {
	authHandler *auth_handler.AuthHandler
	wishHandler *wish_handler.WishHandler
}

func NewAPIImplementation(authClient *client.AuthClient, wishClient *client.WishlisterClient, log *slog.Logger, tokenAuth *jwtauth.JWTAuth) *APIImplementation {
	return &APIImplementation{
		authHandler: auth_handler.NewAuthHandler(authClient, log),
		wishHandler: wish_handler.NewWishHandler(wishClient, log),
	}
}

// Auth handlers
func (a *APIImplementation) PostAuthLogin(w http.ResponseWriter, r *http.Request) {
	a.authHandler.SignIn(w, r)
}

func (a *APIImplementation) PostAuthSignup(w http.ResponseWriter, r *http.Request) {
	a.authHandler.SignUp(w, r)
}

// Wish handlers
func (a *APIImplementation) GetWishWishId(w http.ResponseWriter, r *http.Request, wishId int64) {
	a.wishHandler.GetWish(w, r, wishId)
}

func (a *APIImplementation) GetWishes(w http.ResponseWriter, r *http.Request) {
	a.wishHandler.ListWishes(w, r)
}

func (a *APIImplementation) PostWish(w http.ResponseWriter, r *http.Request) {
	a.wishHandler.CreateWish(w, r)
}

func (a *APIImplementation) PatchWishWishId(w http.ResponseWriter, r *http.Request, wishId int64) {
	a.wishHandler.UpdateWish(w, r, wishId)
}

func (a *APIImplementation) DeleteWishWishId(w http.ResponseWriter, r *http.Request, wishId int64) {
	a.wishHandler.DeleteWish(w, r, wishId)
}

// Share link handlers
func (a *APIImplementation) PostShareLinks(w http.ResponseWriter, r *http.Request) {
	a.wishHandler.CreateShareLink(w, r)
}

func (a *APIImplementation) GetShareLinks(w http.ResponseWriter, r *http.Request) {
	a.wishHandler.GetShareLinks(w, r)
}

func (a *APIImplementation) GetWishesShared(w http.ResponseWriter, r *http.Request) {
	a.wishHandler.GetSharedWishes(w, r)
}

func (a *APIImplementation) DeleteShareLinks(w http.ResponseWriter, r *http.Request) {
	a.wishHandler.DeleteShareLink(w, r)
}

var _ api.ServerInterface = (*APIImplementation)(nil)
