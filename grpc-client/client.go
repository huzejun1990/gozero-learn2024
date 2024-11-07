// @Author huzejun 2024/11/7 22:47:00
package main

import (
	"context"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"grpc-server/greet"
	"log"
)

func main() {
	var clientConf zrpc.RpcClientConf
	conf.MustLoad("etc/client.yml", &clientConf)
	//conn := zrpc.MustNewClient(clientConf)
	conn, err := grpc.NewClient(clientConf.Target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	client := greet.NewGreetClient(conn)
	resp, err := client.Ping(context.Background(), &greet.Request{Ping: "ping"})
	if err != nil {
		log.Fatal(err)
		return
	}

	log.Println(resp)
}
