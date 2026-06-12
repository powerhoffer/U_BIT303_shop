package cmd

import (
	"context"

	"bit303_shop/internal/consts"
	"bit303_shop/internal/controller"
	"bit303_shop/internal/controller/backend"
	"bit303_shop/internal/controller/frontend"
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
		s.Group("/backend", func(group *ghttp.RouterGroup) {
			group.Middleware(
				service.Middleware().CORS,
				service.Middleware().ResponseHandler,
			)
			group.Bind(
				backend.Employee.Register,
				backend.Employee.Login,
			)
			group.Group("/", func(group *ghttp.RouterGroup) {
				group.Middleware(service.Middleware().EmployeeAuth)
				group.Bind(
					backend.Employee.Logout,
					backend.Employee.Info,
					backend.Employee.UpdatePassword,
					backend.Employee.ManageCreate,
					backend.Employee.ManageList,
					backend.Employee.ManageDetail,
					backend.Employee.ManageUpdate,
					backend.Employee.ManageStatus,
					backend.Employee.ManageResetPassword,
					backend.Points.Balance,
					backend.Points.Records,
					backend.Points.ManageAdd,
					backend.Points.ManageDeduct,
					backend.Points.ManageRecords,
					backend.Category.List,
					backend.Goods.Create,
					backend.Goods.List,
					backend.Goods.Detail,
					backend.Goods.Update,
					backend.Goods.Status,
				)
			})
		})
		s.Group("/frontend", func(group *ghttp.RouterGroup) {
			group.Middleware(
				service.Middleware().CORS,
				service.Middleware().ResponseHandler,
			)
			group.Bind(
				frontend.Category.List,
			)
		})
		s.Run()
		return nil
	},
}
