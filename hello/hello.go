package main

import (
	"flag"
	"fmt"

	"gozero-learn2024/hello/internal/config"
	"gozero-learn2024/hello/internal/handler"
	"gozero-learn2024/hello/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/aa.toml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	fmt.Println(c.DataBase.Url, c.NoConfStr, c.NoConfStrDefault)
	server := rest.MustNewServer(rest.RestConf{
		Host: c.Host,
		Port: c.Port,
	})
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
