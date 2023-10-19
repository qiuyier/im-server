package cmd

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcmd"

	"im/internal/controller/hello"
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

			s.Group("/api", func(group *ghttp.RouterGroup) {
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
			})

			s.Run()
			return nil
		},
	}
)
