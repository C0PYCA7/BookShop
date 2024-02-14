package redis

import (
	"context"
	"time"
)

func (r *RedisDb) SaveUserPermissions(id int, permissions string) error {
	ctx := context.Background()

	err := r.client.Set(ctx, string(id), permissions, time.Second*60).Err()
	if err != nil {
		return ErrInternalServer
	}
	return nil
}
