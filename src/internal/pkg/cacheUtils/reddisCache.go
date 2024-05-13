package cache_utils

import (
	"annotater/internal/models"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
)

type ICache interface {
	Set(key string, value interface{}) error
	Get(key string, value interface{}) error
	Del(key string) error
}

type ReddisCache struct {
	ttlDur   time.Duration
	ctx      context.Context
	maxBytes uint
	redis    *redis.Client
}

func NewReddisCache(redSrc *redis.Client, ctxSrc context.Context, maxBytesSrc uint, ttlDurSrc time.Duration) ICache {
	return ReddisCache{
		redis:    redSrc,
		ctx:      ctxSrc,
		maxBytes: maxBytesSrc,
		ttlDur:   ttlDurSrc,
	}

}

func (r ReddisCache) Set(key string, value interface{}) error {
	marhsalledData, err := json.Marshal(value)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("error in reddis chache marshalling value %v", value))
	}
	err = r.redis.Set(r.ctx, key, marhsalledData, r.ttlDur).Err()
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("error in reddis setting value %v", value))
	}
	return nil
}

func (r ReddisCache) Get(key string, value interface{}) error {
	marhsalledData, err := r.redis.Get(r.ctx, key).Result()
	if err == redis.Nil {
		return models.ErrNotFound
	}
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("error in reddis cache  getting marshalled value %v", value))
	}
	err = json.Unmarshal([]byte(marhsalledData), value)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("error in reddis unmarshalling value %v", value))
	}
	fmt.Print("cache used successfully")
	return nil
}

func (r ReddisCache) Del(key string) error {
	err := r.redis.Del(r.ctx, key).Err()
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("error in reddis deleting value by key %v", key))
	}
	return nil
}
