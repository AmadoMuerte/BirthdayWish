package wishlist

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/AmadoMuerte/BirthdayWish/API/internal/config"
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

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf("http://%s/%d", h.cfg.Services.WishListAddr, userID),
		nil,
	)
	if err != nil {
		h.log.Error("failed to create request", "error", err)
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, response.Error("internal error"))
		return
	}

	resp, err := (&http.Client{Timeout: 5 * time.Second}).Do(req)
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

	serviceURL := "http://" + h.cfg.Services.WishListAddr
	bodyBytes, _ := json.Marshal(wishItem)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, serviceURL, bytes.NewBuffer(bodyBytes))
	if err != nil {
		h.log.Error("failed to create request", "error", err)
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, response.Error("internal server error"))
		return
	}

	req.Header.Set("Content-Type", "application/json")
	if auth := r.Header.Get("Authorization"); auth != "" {
		req.Header.Set("Authorization", auth)
	}

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
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

func (h *WishlistHandler) RemoveFromWishlist(w http.ResponseWriter, r *http.Request) {}
