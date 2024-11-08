package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
	"user-api/internal/biz"

	"user-api/internal/config"
	"user-api/internal/handler"
	"user-api/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/user-api.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(
		c.RestConf,
		rest.WithCors("*"),
		rest.WithCorsHeaders("X-Content-Security"),
		rest.WithUnsignedCallback(func(w http.ResponseWriter, r *http.Request, next http.Handler, strict bool, code int) {
			fmt.Println("签名未通过")
		}),
	)

	//server := rest.MustNewServer(
	//	c.RestConf,
	//	rest.WithCors("*"),
	//	rest.WithCorsHeaders("X-Content-Security"),
	//	rest.WithUnsignedCallback(func(w http.ResponseWriter, r *http.Request, next http.Handler, strict bool, code int) {
	//		fmt.Println("签名未通过")
	//	}),
	/*		rest.WithCustomCors(func(header http.Header) {
			var allowOrigin = "Access-Control-Allow-Origin"
			var allOrigins = "http://localhost:5173"
			var allowMethods = "Access-Control-Allow-Methods"
			var allowHeaders = "Access-Control-Allow-Headers"
			var exposeHeaders = "Access-Control-Expose-Headers"
			var methods = "GET, HEAD, POST, PATCH, PUT, DELETE, OPTIONS"
			var allowHeadersVal = "xxxx, Content-Type, Origin, X-CSRF-Token, Authorization, AccessToken, Token, Range"
			var exposeHeadersVal = "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers"
			var maxAgeHeader = "Access-Control-Max-Age"
			var maxAgeHeaderVal = "86400"
			header.Set(allowOrigin, allOrigins)
			header.Set(allowMethods, methods)
			header.Set(allowHeaders, allowHeadersVal)
			header.Set(exposeHeaders, exposeHeadersVal)
			header.Set(maxAgeHeader, maxAgeHeaderVal)*/
	//}
	//)
	server.Use(func(next http.HandlerFunc) http.HandlerFunc {

		return func(writer http.ResponseWriter, request *http.Request) {
			writer.Header().Set("aaaaaaaaaaaaaaa", "true")
			next(writer, request)
		}
	})
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	//统一的错误管理
	httpx.SetErrorHandler(func(err error) (int, any) {
		var e *biz.Error
		switch {
		case errors.As(err, &e):
			return http.StatusOK, biz.Fail(e)
		default:
			return http.StatusInternalServerError, nil
		}
	})
	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
