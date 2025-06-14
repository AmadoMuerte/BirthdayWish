package worker

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/AmadoMuerte/BirthdayWish/API/apps/wishlister/internal/storage"
	redisapi "github.com/AmadoMuerte/BirthdayWish/API/pkg/redis"
)

type Worker struct {
	rdb     *redisapi.RDB
	log     *slog.Logger
	storage *storage.Storage
}

type taskHandler func(ctx context.Context, storage *storage.Storage, data []byte) error

func New(rdb *redisapi.RDB, log *slog.Logger, storage *storage.Storage) *Worker {
	return &Worker{
		rdb:     rdb,
		log:     log.With("component", "worker"),
		storage: storage,
	}
}

func (w *Worker) Start(ctx context.Context) {
	w.log.Info("Starting Redis worker")
	defer w.log.Info("Worker stopped")

	var wg sync.WaitGroup
	wg.Add(2)

	go w.runTaskProcessor(ctx, &wg, redisapi.TaskDeleteWish, w.handleDeleteWish)
	// go w.runTaskProcessor(ctx, &wg, redisapi.TaskUpdateWish, w.handleUpdateWish)

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			wg.Wait()
			return
		case <-ticker.C:
			w.log.Debug("Queue status",
				"delete_queue", w.rdb.Client.LLen(ctx, redisapi.TaskDeleteWish).Val(),
				"update_queue", w.rdb.Client.LLen(ctx, redisapi.TaskUpdateWish).Val())
		}
	}
}

func (w *Worker) runTaskProcessor(
	ctx context.Context,
	wg *sync.WaitGroup,
	queueName string,
	handler taskHandler,
) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		default:
			res, err := w.rdb.PopFromQueue(ctx, queueName)
			if err != nil {
				w.log.Error("Queue error",
					"queue", queueName,
					"error", err)
				time.Sleep(2 * time.Second)
				continue
			}

			if res == nil {
				time.Sleep(1 * time.Second)
				continue
			}

			if err := handler(ctx, w.storage, res); err != nil {
				w.log.Error("Task processing failed",
					"queue", queueName,
					"error", err)
			}
		}
	}
}
