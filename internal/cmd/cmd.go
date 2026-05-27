package cmd

import (
	"context"

	"bit303_shop/internal/consts"
	"bit303_shop/internal/controller"
	"bit303_shop/internal/service"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcmd"
)

var Main = gcmd.Command{
	Name:  consts.ProjectName,
	Usage: consts.ProjectUsage,
	Brief: consts.ProjectBrief,
	Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
		s := g.Server()
		s.Group("/", func(group *ghttp.RouterGroup) {
			group.Middleware(
				service.Middleware().CORS,
				service.Middleware().ResponseHandler,
			)
			group.Bind(controller.Health)
		})
		s.Run()
		return nil
	},
}
