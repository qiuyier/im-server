package cmd

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcmd"
	"golang.org/x/sync/errgroup"
	"im/internal/controller/bar"
	"im/internal/controller/hello"
	"im/internal/service"
	"im/utility/websocket"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	Main = gcmd.Command{
		Name:  "main",
		Usage: "main",
		Brief: "start http server",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			s := g.Server()
			s.Group("/", func(group *ghttp.RouterGroup) {
				group.Middleware(ghttp.MiddlewareHandlerResponse)
				group.Bind(
					hello.NewV1(),
				)
			})

			// 启动商家前端gToken
			gfUserToken, err := StartUserGToken()
			if err != nil {
				return err
			}

			s.Group("/api/v1", func(group *ghttp.RouterGroup) {
				group.Middleware(
					ghttp.MiddlewareHandlerResponse,
					ghttp.MiddlewareCORS,
				)

				// 需要gToken
				group.Group("/auth", func(group *ghttp.RouterGroup) {
					err := gfUserToken.Middleware(ctx, group)
					if err != nil {
						panic(err)
					}
				})

				group.Group("/publish", func(group *ghttp.RouterGroup) {
					group.Bind(
						bar.NewV1(),
					)
				})
			})

			s.Group("/", func(group *ghttp.RouterGroup) {
				err := gfUserToken.Middleware(ctx, group)
				if err != nil {
					panic(err)
				}

				group.ALL("/wss/wss.io", func(r *ghttp.Request) {
					eg, groupCtx := errgroup.WithContext(ctx)

					// 初始化IM渠道配置
					websocket.Init(groupCtx, eg, func(name string) {
						g.Dump("守护进程异常", fmt.Sprintf("守护进程异常[%s]", name))
					})

					c := make(chan os.Signal, 1)
					signal.Notify(c, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT)

					// 启动消息订阅
					time.AfterFunc(3*time.Second, func() {
						service.ServerSubscribe().Start(groupCtx, eg)
					})

					// 启动websocket连接
					err = service.ServerSubscribe().Conn(r.Response.ResponseWriter, r.Request)
					if err != nil {
						panic(err)
					}
				})
			})

			s.Run()
			return nil
		},
	}
)
