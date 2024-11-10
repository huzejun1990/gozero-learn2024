package config

import (
	"github.com/zeromicro/go-queue/kq"
	configurator "github.com/zeromicro/go-zero/core/configcenter"
	"github.com/zeromicro/go-zero/core/configcenter/subscriber"
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

func PullConfig() Config {
	ss := subscriber.MustNewEtcdSubscriber(subscriber.EtcdConf{
		Hosts: []string{"localhost:2379"}, // etcd 地址
		Key:   "user-api-test",            // 配置key
	})

	// 创建 configurator
	cc := configurator.MustNewConfigCenter[Config](configurator.Config{
		Type: "yaml", // 配置值类型：json,yaml,toml
	}, ss)

	// 获取配置
	// 注意: 配置如果发生变更，调用的结果永远获取到最新的配置
	v, err := cc.GetConfig()
	if err != nil {
		panic(err)
	}
	cc.AddListener(func() {
		v, err := cc.GetConfig()
		if err != nil {
			panic(err)
		}
		//这个地方要写 触发配置变化后 需要处理的操作
		println("config changed:", v.Name)
	})
	// 如果想监听配置变化，可以添加 listener
	return v
}

func PullConsulConfig() Config { // /user-config/test/user-api
	//ss := NewConsulSubscriber("http://localhost:8500", "user-config/test/user-api")
	ss := NewConsulSubscriber("http://localhost:8500", "user-config/test/user-api")

	// 创建 configurator
	cc := configurator.MustNewConfigCenter[Config](configurator.Config{
		Type: "yaml", // 配置值类型：json,yaml,toml
	}, ss)

	// 获取配置
	// 注意: 配置如果发生变更，调用的结果永远获取到最新的配置
	v, err := cc.GetConfig()
	if err != nil {
		panic(err)
	}
	cc.AddListener(func() {
		v, err := cc.GetConfig()
		if err != nil {
			panic(err)
		}
		//这个地方要写 触发配置变化后 需要处理的操作
		println("config changed:", v.Name)
	})
	// 如果想监听配置变化，可以添加 listener
	return v
}
