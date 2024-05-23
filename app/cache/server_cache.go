package cache

import "vcs_server/entity"

type ServerCache interface {
	Set(key string, value *entity.Server)
	Get(key string) *entity.Server
	Del(key string) error
	ASet(key string, value []entity.Server)
	AGet(key string) []entity.Server
	Ping()
}
