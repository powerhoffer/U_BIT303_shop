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
		uploadPath := g.Cfg().MustGet(ctx, "upload.path").String()
		uploadPrefix := g.Cfg().MustGet(ctx, "upload.urlPrefix").String()
		if uploadPath != "" && uploadPrefix != "" {
			s.AddStaticPath(uploadPrefix, uploadPath)
		}
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
				backend.Admin.Login,
			)
			group.Group("/", func(group *ghttp.RouterGroup) {
				group.Middleware(service.Middleware().EmployeeAuth)
				group.Bind(
					backend.Employee.Logout,
					backend.Employee.Info,
					backend.Employee.UpdatePassword,
					backend.Points.Balance,
					backend.Points.Records,
				)
			})
			group.Group("/", func(group *ghttp.RouterGroup) {
				group.Middleware(service.Middleware().AdminAuth)
				group.Bind(
					backend.Admin.Logout,
					backend.Admin.Info,
					backend.Admin.UpdatePassword,
				)
				group.Group("/", func(group *ghttp.RouterGroup) {
					group.Middleware(service.Middleware().AdminPermissionAuth)
					group.Bind(
						backend.Employee.ManageCreate,
						backend.Employee.ManageList,
						backend.Employee.ManageDetail,
						backend.Employee.ManageUpdate,
						backend.Employee.ManageStatus,
						backend.Employee.ManageResetPassword,
						backend.Employee.ManageDelete,
						backend.Points.ManageAdd,
						backend.Points.ManageBatchAdd,
						backend.Points.ManageDeduct,
						backend.Points.ManageRecords,
						backend.Category.List,
						backend.Goods.Create,
						backend.Goods.List,
						backend.Goods.Detail,
						backend.Goods.Update,
						backend.Goods.Status,
						backend.Stock.Adjust,
						backend.Stock.Records,
						backend.Order.List,
						backend.Order.Detail,
						backend.Order.Complete,
						backend.Order.Cancel,
						backend.Admin.ManageCreate,
						backend.Admin.ManageList,
						backend.Admin.ManageDetail,
						backend.Admin.ManageUpdate,
						backend.Admin.ManageStatus,
						backend.Admin.ManageResetPassword,
						backend.Admin.ManageRoles,
						backend.Role.Create,
						backend.Role.List,
						backend.Role.Detail,
						backend.Role.Update,
						backend.Role.Status,
						backend.Role.Permissions,
						backend.Permission.Create,
						backend.Permission.List,
						backend.Permission.Detail,
						backend.Permission.Update,
						backend.Permission.Status,
						backend.Upload.GoodsImage,
					)
				})
			})
		})
		s.Group("/frontend", func(group *ghttp.RouterGroup) {
			group.Middleware(
				service.Middleware().CORS,
				service.Middleware().ResponseHandler,
			)
			group.Bind(
				frontend.Category.List,
				frontend.Goods.List,
				frontend.Goods.Detail,
			)
			group.Group("/", func(group *ghttp.RouterGroup) {
				group.Middleware(service.Middleware().EmployeeAuth)
				group.Bind(
					frontend.Cart.Add,
					frontend.Cart.List,
					frontend.Cart.Update,
					frontend.Cart.Remove,
					frontend.Order.Create,
					frontend.Order.List,
					frontend.Order.Detail,
					frontend.Order.Cancel,
				)
			})
		})
		s.Run()
		return nil
	},
}
