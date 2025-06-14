package worker

import (
	"context"

	"github.com/AmadoMuerte/BirthdayWish/API/apps/wishlister/internal/storage"
)

func (w *Worker) handleUpdateWish(ctx context.Context, storage *storage.Storage, data []byte) error {
	// var task apimodels.UpdateWishTask
	// if err := json.Unmarshal(data, &task); err != nil {
	// 	return fmt.Errorf("unmarshal error: %w", err)
	// }

	// w.log.Info("Processing update wish",
	// 	"user_id", task.UserID,
	// 	"wish_id", task.WishID)

	// // if err := w.service.UpdateWish(ctx, task); err != nil {
	// //     return err
	// // }
	return nil
}
