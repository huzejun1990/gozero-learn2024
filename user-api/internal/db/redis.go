// @Author huzejun 2024/11/7 17:34:00
package db

import "github.com/zeromicro/go-zero/core/stores/redis"

func NewRedis(con redis.RedisConf) *redis.Redis {
	return redis.MustNewRedis(con)
}
