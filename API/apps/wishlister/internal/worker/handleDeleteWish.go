package worker

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/AmadoMuerte/BirthdayWish/API/apps/wishlister/internal/storage"
	"github.com/AmadoMuerte/BirthdayWish/API/pkg/apimodels"
)

func (w *Worker) handleDeleteWish(ctx context.Context, storage *storage.Storage, data []byte) error {
	op := "worker/handleDeleteWish"

	var task apimodels.DeleteWishTask
	if err := json.Unmarshal(data, &task); err != nil {
		return fmt.Errorf("unmarshal error: %s", err)
	}

	err := w.storage.RemoveFromWishlist(ctx, task.WishID, task.UserID)
	if err != nil {
		w.log.Error(op+": internal server error", "error", err)
		return err
	}

	w.log.Info("Wish was deleted")

	return nil
}
