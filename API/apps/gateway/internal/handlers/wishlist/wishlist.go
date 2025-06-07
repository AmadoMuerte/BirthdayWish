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

	"github.com/AmadoMuerte/BirthdayWish/API/apps/gateway/internal/config"
	"github.com/AmadoMuerte/BirthdayWish/API/apps/gateway/internal/handlers/filer"
	"github.com/AmadoMuerte/BirthdayWish/API/apps/gateway/internal/storage"
	"github.com/AmadoMuerte/BirthdayWish/API/pkg/httphelper"
	"github.com/AmadoMuerte/BirthdayWish/API/pkg/jwt"
	"github.com/AmadoMuerte/BirthdayWish/API/pkg/response"
	"github.com/go-chi/chi/v5"
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
	op := "wishlist/getShareList"
	ctx := r.Context()

	authHeader := r.Header.Get("Authorization")
	token := strings.TrimPrefix(authHeader, "Bearer ")
	if token == authHeader {
		h.log.Error(op+ ": token error")
		response.ErrorResponseJSON(w,r, http.StatusBadRequest, "internal error")
		return
	}

	userID, err := h.storage.GetUserIDByToken(ctx, token)
	if err != nil {
		h.log.Error(op + ": userID", "error", err)
		response.ErrorResponseJSON(w,r, http.StatusBadRequest, "internal error")
		return
	}

	path := fmt.Sprintf("%s/%d", h.cfg.Services.WishListAddr, userID)
	resp, err := httphelper.DoRequest(ctx, "GET", path, nil, nil)
	if err != nil {
		h.log.Error(op + ": service call failed", "error", err)
		response.ErrorResponseJSON(w,r, http.StatusBadGateway, "service unavailable")
		return
	}

	defer resp.Body.Close()

	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func (h *WishlistHandler) GetWishlist(w http.ResponseWriter, r *http.Request) {
	op := "wishlist/GetWishlist"
	ctx := r.Context()

	claims, err := jwt.GetClaims(ctx)
	if err != nil {
		h.log.Error(op + ": failed to get claims", "error", err)
		response.ErrorResponseJSON(w,r, http.StatusInternalServerError, "internal server error")
		return
	}

	path := fmt.Sprintf("%s/%d", h.cfg.Services.WishListAddr, claims.UserID)

	resp, err := httphelper.DoRequest(ctx, "GET", path, nil, nil)
	if err != nil {
		h.log.Error(op + ": service call failed", "error", err)
		response.ErrorResponseJSON(w,r, http.StatusBadGateway, "service unavailable")
		return
	}

	defer resp.Body.Close()

	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func (h *WishlistHandler) GetWish(w http.ResponseWriter, r *http.Request) {
	op := "wishlist/GetWish"
	ctx := r.Context()

	wishID, err := strconv.ParseInt(chi.URLParam(r, "wish_id"), 10, 64)
	if err != nil {
		h.log.Error(op + ": invalid wish id", "error", err)
		response.ErrorResponseJSON(w,r, http.StatusBadRequest, "invalid wish id")
		return
	}

	claims, err := jwt.GetClaims(ctx)

	path := fmt.Sprintf("%s/%d/%d", h.cfg.Services.WishListAddr, wishID, claims.UserID)

	resp, err := httphelper.DoRequest(ctx, "GET", path, nil, nil)
	if err != nil {
		h.log.Error(op + ": service call failed", "error", err)
		response.ErrorResponseJSON(w,r, http.StatusBadGateway, "service unavailable")
		return
	}

	defer resp.Body.Close()

	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}


func (h *WishlistHandler) AddToWishlist(w http.ResponseWriter, r *http.Request) {
	op := "wishlist/AddToWishlist"
	ctx := r.Context()

	var wishItemReq wishItemReq
	if err := json.NewDecoder(r.Body).Decode(&wishItemReq); err != nil {
		h.log.Error(op + ": failed to decode request body", "error", err)
		response.ErrorResponseJSON(w,r, http.StatusBadRequest, "invalid request body")
		return
	}
	defer r.Body.Close()

	claims, err := jwt.GetClaims(ctx)
	if err != nil {
		h.log.Error(op + ": failed to get claims from token", "error", err)
		response.ErrorResponseJSON(w,r, http.StatusUnauthorized, "invalid token")
		return
	}

	if exists, err := h.storage.UserExistsByID(ctx, claims.UserID); err != nil || !exists {
		h.log.Error(op + ": user does not exist", "user_id", claims.UserID, "error", err)
		response.ErrorResponseJSON(w,r, http.StatusUnauthorized, "user not found")
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
		resp, err = httphelper.DoRequest(
			ctx,
			"POST",
			path,
			map[string]string{"data": wishItemReq.Image, "type": wishItemReq.ImageType},
			map[string]string{"1": "application/json", "2": "Content-Type"},
		)
		if resp.StatusCode != 201 || err != nil {
			h.log.Error(op + ": filer service is not create image", "error", err)
			response.ErrorResponseJSON(w,r, http.StatusBadGateway, "service unavailable")
			return
		}

		var imageRecord filer.ImageRecord
		if err := json.NewDecoder(resp.Body).Decode(&imageRecord); err != nil {
			h.log.Error(op + ": failed to decode request body", "error", err)
			response.ErrorResponseJSON(w,r, http.StatusBadRequest, "invalid request body")
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

	resp, err = httphelper.DoRequest(ctx, "POST", path, wishItemRes, headers)
	if err != nil {
		h.log.Error(op + ": failed add wish item", "error", err)
		response.ErrorResponseJSON(w,r, http.StatusInternalServerError, "internal server error")
		return
	}

	defer resp.Body.Close()

	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func (h *WishlistHandler) RemoveFromWishlist(w http.ResponseWriter, r *http.Request) {
	op := "wishlist/RemoveFromWishlist"
	ctx := r.Context()

	wishID, err := strconv.ParseInt(chi.URLParam(r, "wish_id"), 10, 64)
	if err != nil {
		h.log.Error(op + ": wish_id is empty", "error", err)
		response.ErrorResponseJSON(w,r, http.StatusBadRequest, "invalid request body")
		return
	}

	claims, err := jwt.GetClaims(ctx)
	if err != nil {
		h.log.Error(op + ": failed to get claims from token", "error", err)
		response.ErrorResponseJSON(w,r, http.StatusUnauthorized, "Unauthorized")
		return
	}

	exists, err := h.storage.UserExistsByID(ctx, claims.UserID)
	if err != nil || !exists {
		h.log.Error(op + ": user does not exist", "error", err)
		response.ErrorResponseJSON(w,r, http.StatusBadRequest, "invalid request body")
		return
	}

	path := fmt.Sprintf("%s/%d/%d", h.cfg.Services.WishListAddr, wishID, claims.UserID)

	resp, err := httphelper.DoRequest(ctx, "DELETE", path, nil, nil)
	if err != nil {
		h.log.Error(op + ": service call failed", "error", err)
		response.ErrorResponseJSON(w,r, http.StatusBadGateway, "service unavailable")
		return
	}

	defer resp.Body.Close()

	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}


func (h *WishlistHandler) PartialUpdateWish(w http.ResponseWriter, r *http.Request) {
	op := "wishlist/PartialUpdateWish"
    ctx := r.Context()

    bodyBytes, err := io.ReadAll(r.Body)
    if err != nil {
		h.log.Error(op + ": failed to read request body", "error", err)
		response.ErrorResponseJSON(w,r, http.StatusBadRequest, "invalid request body")
        return
    }
    defer r.Body.Close()

    var wishItemReq wishItemReq
    if err := json.Unmarshal(bodyBytes, &wishItemReq); err != nil {
		h.log.Error(op + ": failed to decode request body", "error", err)
		response.ErrorResponseJSON(w,r, http.StatusBadRequest, "invalid request body")
        return
    }

    var updateData map[string]any
    if err := json.Unmarshal(bodyBytes, &updateData); err != nil {
		h.log.Error(op + ": failed to decode request body", "error", err)
		response.ErrorResponseJSON(w,r, http.StatusBadRequest, "invalid request body")
        return
    }

    wishID, err := strconv.ParseInt(chi.URLParam(r, "wish_id"), 10, 64)
    if err != nil {
		h.log.Error(op + ": invalid wish_id", "error", err)
		response.ErrorResponseJSON(w,r, http.StatusBadRequest, "invalid request body")
        return
    }

    claims, err := jwt.GetClaims(ctx)
    if err != nil {
		h.log.Error(op + ": failed to get claims from token", "error", err)
		response.ErrorResponseJSON(w,r, http.StatusUnauthorized, "Unauthorized")
        return
    }

    exists, err := h.storage.UserExistsByID(ctx, claims.UserID)
    if err != nil || !exists {
		h.log.Error(op + "user does not exist", "user_id", claims.UserID, "error", err)
		response.ErrorResponseJSON(w,r, http.StatusNotFound, "user not found")
        return
    }

    if wishItemReq.Image != "" && wishItemReq.ImageType != "" {
        imagePath := fmt.Sprintf("%s/images", h.cfg.Services.Minio)
        imageResp, err := httphelper.DoRequest(
            ctx,
            "POST",
            imagePath,
            map[string]string{"data": wishItemReq.Image, "type": wishItemReq.ImageType},
            map[string]string{"Content-Type": "application/json"},
        )
        if err != nil || imageResp.StatusCode != http.StatusCreated {
			h.log.Error(op + "filer service failed", "error", err)
			response.ErrorResponseJSON(w,r, http.StatusBadGateway, "service unavailable")
            return
        }
        defer imageResp.Body.Close()

        var imageRecord filer.ImageRecord
        if err := json.NewDecoder(imageResp.Body).Decode(&imageRecord); err != nil {
			h.log.Error(op + "failed to decode image response", "error", err)
			response.ErrorResponseJSON(w,r, http.StatusInternalServerError, "service unavailable")
            return
        }

        updateData["image_url"] = imageRecord.PublicURL
        updateData["image_name"] = imageRecord.OriginalName
    }

    path := fmt.Sprintf("%s/%d/%d", h.cfg.Services.WishListAddr, wishID, claims.UserID)
    
    resp, err := httphelper.DoRequest(
        ctx,
        "PATCH",
        path,
        updateData,
        map[string]string{"Content-Type": "application/json"},
    )
    if err != nil {
		h.log.Error(op + "service call failed", "error", err)
		response.ErrorResponseJSON(w,r, http.StatusBadGateway, "service unavailable")
        return
    }
    defer resp.Body.Close()

    w.WriteHeader(resp.StatusCode)
    io.Copy(w, resp.Body)
}