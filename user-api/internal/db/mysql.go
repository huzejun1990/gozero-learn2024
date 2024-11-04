// @Author huzejun 2024/11/4 20:38:00
package db

import (
	"context"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"time"
	"user-api/internal/config"
)

func NewMysql(mysqlConfig config.MysqlConfig) sqlx.SqlConn {
	mysql := sqlx.NewMysql(mysqlConfig.DataSource)
	db, err := mysql.RawDB()
	if err != nil {
		panic(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(mysqlConfig.ConnectTimeout))
	defer cancel()
	err = db.PingContext(ctx)
	if err != nil {
		panic(err)
	}
	db.SetMaxOpenConns(100)
	db.SetMaxIdleConns(10)
	return mysql
}
