package wishlist

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/AmadoMuerte/BirthdayWish/API/internal/config"
	http_helper "github.com/AmadoMuerte/BirthdayWish/API/internal/lib"
	"github.com/AmadoMuerte/BirthdayWish/API/internal/storage"
	"github.com/AmadoMuerte/BirthdayWish/API/pkg/jwt"
	"github.com/AmadoMuerte/BirthdayWish/API/pkg/response"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type WishlistHandler struct {
	cfg     *config.Config
	storage *storage.Storage
	log     *slog.Logger
}

type wishItem struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	Price     float64   `json:"price"`
	Link      string    `json:"link"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type IWishlistHandler interface {
	GetWishlist(w http.ResponseWriter, r *http.Request)
	AddToWishlist(w http.ResponseWriter, r *http.Request)
	RemoveFromWishlist(w http.ResponseWriter, r *http.Request)
}

func New(cfg *config.Config, storage *storage.Storage, log *slog.Logger) IWishlistHandler {
	return &WishlistHandler{cfg, storage, log}
}

func (h *WishlistHandler) GetWishlist(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, err := strconv.ParseInt(chi.URLParam(r, "user_id"), 10, 64)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.Error("invalid user id"))
		return
	}

	claims, err := jwt.GetClaims(ctx)
	if err != nil || claims.UserID != userID {
		render.Status(r, http.StatusForbidden)
		render.JSON(w, r, response.Error("access denied"))
		return
	}

	path := fmt.Sprintf("%s/%d", h.cfg.Services.WishListAddr, userID)
	
	resp, err := http_helper.DoRequest(ctx, "GET", path, nil, nil)
	if err != nil {
		h.log.Error("service call failed", "error", err)
		render.Status(r, http.StatusBadGateway)
		render.JSON(w, r, response.Error("service unavailable"))
		return
	}
	
	defer resp.Body.Close()

	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func (h *WishlistHandler) AddToWishlist(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var wishItem wishItem
	if err := json.NewDecoder(r.Body).Decode(&wishItem); err != nil {
		h.log.Error("failed to decode request body", "error", err)
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.Error("invalid request body"))
		return
	}
	defer r.Body.Close()

	claims, err := jwt.GetClaims(ctx)
	if err != nil {
		h.log.Error("failed to get claims from token", "error", err)
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, response.Error("invalid token"))
		return
	}

	if exists, err := h.storage.UserExistsByID(ctx, claims.UserID); err != nil || !exists {
		h.log.Error("user does not exist", "user_id", claims.UserID, "error", err)
		render.Status(r, http.StatusNotFound)
		render.JSON(w, r, response.Error("user not found"))
		return
	}
	wishItem.UserID = claims.UserID

	path := fmt.Sprintf("%s", h.cfg.Services.WishListAddr)
	headers := map[string]string{
		"1":"Content-Type", 
		"2":"application/json",
	}
	
	resp, err := http_helper.DoRequest(ctx, "GET", path, wishItem, headers)
	if err != nil {
		h.log.Error("service call failed", "error", err)
		render.Status(r, http.StatusBadGateway)
		render.JSON(w, r, response.Error("service unavailable"))
		return
	}
	
	defer resp.Body.Close()

	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func (h *WishlistHandler) RemoveFromWishlist(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	wishID, err := strconv.ParseInt(chi.URLParam(r, "wish_id"), 10, 64)
	if err != nil {
		h.log.Error("wish_id is empty")
		render.JSON(w, r, response.Error("wish_id is empty"))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	claims, err := jwt.GetClaims(ctx)
	if err != nil {
		h.log.Error("failed to get claims from token", "error", err)
		w.WriteHeader(http.StatusUnauthorized)
		render.JSON(w, r, response.Error("invalid token"))
		return
	}

	exists, err := h.storage.UserExistsByID(ctx, claims.UserID)
	if err != nil || !exists {
		fmt.Printf("user id is %d", claims.UserID)
		h.log.Error("user does not exist", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		render.JSON(w, r, response.Error("user does not exist"))
		return
	}

	path := fmt.Sprintf("%s/%d/%d", h.cfg.Services.WishListAddr, wishID, claims.UserID)

	resp, err := http_helper.DoRequest(ctx, "DELETE", path, nil, nil)
	if err != nil {
		h.log.Error("service call failed", "error", err)
		render.Status(r, http.StatusBadGateway)
		render.JSON(w, r, response.Error("service unavailable"))
		return
	}

	defer resp.Body.Close()

	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}
