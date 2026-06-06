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
		s.Group("/employee", func(group *ghttp.RouterGroup) {
			group.Middleware(
				service.Middleware().CORS,
				service.Middleware().ResponseHandler,
			)
			group.Bind(
				controller.Employee.Register,
				controller.Employee.Login,
			)
			group.Group("/", func(group *ghttp.RouterGroup) {
				group.Middleware(service.Middleware().EmployeeAuth)
				group.Bind(
					controller.Employee.Logout,
					controller.Employee.Info,
					controller.Employee.UpdatePassword,
					controller.Employee.ManageCreate,
					controller.Employee.ManageList,
					controller.Employee.ManageDetail,
					controller.Employee.ManageUpdate,
					controller.Employee.ManageStatus,
					controller.Employee.ManageResetPassword,
				)
			})
		})
		s.Group("/points", func(group *ghttp.RouterGroup) {
			group.Middleware(
				service.Middleware().CORS,
				service.Middleware().ResponseHandler,
				service.Middleware().EmployeeAuth,
			)
			group.Bind(
				controller.Points.Balance,
				controller.Points.Records,
				controller.Points.ManageAdd,
				controller.Points.ManageDeduct,
				controller.Points.ManageRecords,
			)
		})
		s.Run()
		return nil
	},
}
