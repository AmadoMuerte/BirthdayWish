package apimodels

type DeleteWishTask struct {
	UserID int64 `json:"user_id"`
	WishID int64 `json:"wish_id"`
}
