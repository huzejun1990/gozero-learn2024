package config

import (
	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf
	MysqlConfig    MysqlConfig
	Auth           Auth
	RedisConfig    redis.RedisConf
	KqPusherConf   kq.KqConf
	KqConsumerConf kq.KqConf
}

type Auth struct {
	AccessSecret string
	Expire       int64
}

type MysqlConfig struct {
	DataSource     string
	ConnectTimeout int
}
