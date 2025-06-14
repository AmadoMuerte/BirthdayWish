package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

func (r *RDB) PushToQueue(ctx context.Context, queueName string, data any) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return r.Client.RPush(ctx, queueName, jsonData).Err()
}

func (r *RDB) PopFromQueue(ctx context.Context, queueName string) ([]byte, error) {
	res, err := r.Client.BLPop(ctx, 1*time.Second, queueName).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, fmt.Errorf("BLPop failed: %w", err)
	}

	if len(res) < 2 {
		return nil, errors.New("invalid queue response: expected [key, value]")
	}

	return []byte(res[1]), nil
}
