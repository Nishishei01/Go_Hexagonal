package ports

import "time"

type CacheRepository interface {
	Set(key string, value any, expiration time.Duration) error
	Get(key string, dest any) error
	Delete(key string) error
}
