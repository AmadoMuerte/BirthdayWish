package wishlist

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/AmadoMuerte/BirthdayWish/API/internal/config"
	"github.com/AmadoMuerte/BirthdayWish/API/internal/http_server/handlers/minio"
	http_helper "github.com/AmadoMuerte/BirthdayWish/API/internal/lib"
	"github.com/AmadoMuerte/BirthdayWish/API/internal/pkg/jwt"
	"github.com/AmadoMuerte/BirthdayWish/API/internal/pkg/response"
	"github.com/AmadoMuerte/BirthdayWish/API/internal/storage"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type WishlistHandler struct {
	cfg     *config.Config
	storage *storage.Storage
	log     *slog.Logger
}

type wishItemReq struct {
	Price     float64 `json:"price"`
	Link      string  `json:"link"`
	Image     string  `json:"image_data"`
	ImageType string  `json:"image_type"`
	Name      string  `json:"name"`
}

type wishItemRes struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	Price     float64   `json:"price"`
	Link      string    `json:"link"`
	ImageUrl  string    `json:"image_url"`
	ImageName string    `json:"image_name"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type IWishlistHandler interface {
	GetWishlist(w http.ResponseWriter, r *http.Request)
	AddToWishlist(w http.ResponseWriter, r *http.Request)
	RemoveFromWishlist(w http.ResponseWriter, r *http.Request)
	GetShareList(w http.ResponseWriter, r *http.Request)
	GetWish(w http.ResponseWriter, r *http.Request)
	PartialUpdateWish(w http.ResponseWriter, r *http.Request)
}

func New(cfg *config.Config, storage *storage.Storage, log *slog.Logger) IWishlistHandler {
	return &WishlistHandler{cfg, storage, log}
}

func (h *WishlistHandler) GetShareList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	authHeader := r.Header.Get("Authorization")
	token := strings.TrimPrefix(authHeader, "Bearer ")
	if token == authHeader {
		h.log.Error("token error")
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.Error("internal error"))
		return
	}

	userID, err := h.storage.GetUserIDByToken(ctx, token)
	if err != nil {
		h.log.Error("userID", "error", err)
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.Error("internal error"))
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

func (h *WishlistHandler) GetWish(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	wishID, err := strconv.ParseInt(chi.URLParam(r, "wish_id"), 10, 64)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.Error("invalid wish id"))
		return
	}

	claims, err := jwt.GetClaims(ctx)

	path := fmt.Sprintf("%s/%d/%d", h.cfg.Services.WishListAddr, wishID, claims.UserID)

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

	var wishItemReq wishItemReq
	if err := json.NewDecoder(r.Body).Decode(&wishItemReq); err != nil {
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

	wishItemRes := wishItemRes{
		UserID: claims.UserID,
		Price:  wishItemReq.Price,
		Link:   wishItemReq.Link,
		Name:   wishItemReq.Name,
	}

	var path string
	var resp *http.Response

	if wishItemReq.Image != "" && wishItemReq.ImageType != "" {
		path := fmt.Sprintf("%s/%s", h.cfg.Services.Minio, "images")
		resp, err = http_helper.DoRequest(
			ctx,
			"POST",
			path,
			map[string]string{"data": wishItemReq.Image, "type": wishItemReq.ImageType},
			map[string]string{"1": "application/json", "2": "Content-Type"},
		)
		if resp.StatusCode != 201 || err != nil {
			h.log.Error("minio service is not create image", "error", err)
			render.Status(r, http.StatusBadGateway)
			render.JSON(w, r, response.Error("service unavailable"))
			return
		}

		var imageRecord minio.ImageRecord
		if err := json.NewDecoder(resp.Body).Decode(&imageRecord); err != nil {
			h.log.Error("failed to decode request body", "error", err)
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("invalid request body"))
			return
		}
		wishItemRes.ImageUrl = imageRecord.PublicURL
		wishItemRes.ImageName = imageRecord.OriginalName
	}

	path = fmt.Sprintf("%s", h.cfg.Services.WishListAddr)
	headers := map[string]string{
		"1": "Content-Type",
		"2": "application/json",
	}

	resp, err = http_helper.DoRequest(ctx, "POST", path, wishItemRes, headers)
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


func (h *WishlistHandler) PartialUpdateWish(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()

    bodyBytes, err := io.ReadAll(r.Body)
    if err != nil {
        h.log.Error("failed to read request body", "error", err)
        render.Status(r, http.StatusBadRequest)
        render.JSON(w, r, response.Error("invalid request body"))
        return
    }
    defer r.Body.Close()

    var wishItemReq wishItemReq
    if err := json.Unmarshal(bodyBytes, &wishItemReq); err != nil {
        h.log.Error("failed to decode request body", "error", err)
        render.Status(r, http.StatusBadRequest)
        render.JSON(w, r, response.Error("invalid request body"))
        return
    }

    var updateData map[string]any
    if err := json.Unmarshal(bodyBytes, &updateData); err != nil {
        h.log.Error("failed to decode request body", "error", err)
        render.Status(r, http.StatusBadRequest)
        render.JSON(w, r, response.Error("invalid request body"))
        return
    }

    wishID, err := strconv.ParseInt(chi.URLParam(r, "wish_id"), 10, 64)
    if err != nil {
        h.log.Error("invalid wish_id", "error", err)
        render.JSON(w, r, response.Error("invalid wish_id"))
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
        h.log.Error("user does not exist", "user_id", claims.UserID, "error", err)
        w.WriteHeader(http.StatusNotFound)
        render.JSON(w, r, response.Error("user not found"))
        return
    }

    if wishItemReq.Image != "" && wishItemReq.ImageType != "" {
        imagePath := fmt.Sprintf("%s/images", h.cfg.Services.Minio)
        imageResp, err := http_helper.DoRequest(
            ctx,
            "POST",
            imagePath,
            map[string]string{"data": wishItemReq.Image, "type": wishItemReq.ImageType},
            map[string]string{"Content-Type": "application/json"},
        )
        if err != nil || imageResp.StatusCode != http.StatusCreated {
            h.log.Error("minio service failed", "error", err)
            render.Status(r, http.StatusBadGateway)
            render.JSON(w, r, response.Error("image service unavailable"))
            return
        }
        defer imageResp.Body.Close()

        var imageRecord minio.ImageRecord
        if err := json.NewDecoder(imageResp.Body).Decode(&imageRecord); err != nil {
            h.log.Error("failed to decode image response", "error", err)
            render.Status(r, http.StatusInternalServerError)
            render.JSON(w, r, response.Error("failed to process image"))
            return
        }

        updateData["image_url"] = imageRecord.PublicURL
        updateData["image_name"] = imageRecord.OriginalName
    }

    path := fmt.Sprintf("%s/%d/%d", h.cfg.Services.WishListAddr, wishID, claims.UserID)
    
    resp, err := http_helper.DoRequest(
        ctx,
        "PATCH",
        path,
        updateData,
        map[string]string{"Content-Type": "application/json"},
    )
    if err != nil {
        h.log.Error("service call failed", "error", err)
        render.Status(r, http.StatusBadGateway)
        render.JSON(w, r, response.Error("service unavailable"))
        return
    }
    defer resp.Body.Close()

    w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
    w.WriteHeader(resp.StatusCode)
    if _, err := io.Copy(w, resp.Body); err != nil {
        h.log.Error("failed to write response", "error", err)
    }
}