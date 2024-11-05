package config

import "github.com/zeromicro/go-zero/rest"

type Config struct {
	rest.RestConf
	MysqlConfig MysqlConfig
	Auth        Auth
}

type Auth struct {
	AccessSecret string
	Expire       int64
}

type MysqlConfig struct {
	DataSource     string
	ConnectTimeout int
}
