package svc

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"user-api/internal/config"
	"user-api/internal/db"
)

type ServiceContext struct {
	Config config.Config
	Mysql  sqlx.SqlConn
	//Redis  *redis.Redis
	RedisClient *redis.Redis
}

func NewServiceContext(c config.Config) *ServiceContext {
	mysql := db.NewMysql(c.MysqlConfig)
	redisClient := db.NewRedis(c.RedisConfig)
	return &ServiceContext{
		Config:      c,
		Mysql:       mysql,
		RedisClient: redisClient,
	}
}
