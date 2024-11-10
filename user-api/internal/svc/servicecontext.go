package svc

import (
	"github.com/segmentio/kafka-go"
	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"time"
	"user-api/internal/config"
	"user-api/internal/db"
)

type ServiceContext struct {
	Config       config.Config
	Mysql        sqlx.SqlConn
	RedisClient  *redis.Redis
	KafkaPushCli *kq.Pusher
}

func NewServiceContext(c config.Config) *ServiceContext {
	mysql := db.NewMysql(c.MysqlConfig)
	redisClient := db.NewRedis(c.RedisConfig)
	pusher := kq.NewPusher(
		c.KqPusherConf.Brokers,
		c.KqPusherConf.Topic,
		kq.WithChunkSize(1024),
		kq.WithFlushInterval(time.Second),
		kq.WithAllowAutoTopicCreation(),
		kq.WithBalancer(&kafka.Hash{}),
	)
	return &ServiceContext{
		Config:       c,
		Mysql:        mysql,
		RedisClient:  redisClient,
		KafkaPushCli: pusher,
	}
}
