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

var (
	Main = gcmd.Command{
		Name:  consts.ProjectName,
		Usage: consts.ProjectUsage,
		Brief: consts.ProjectBrief,
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			s := g.Server()
			//订单超时未评价默认好评
			//err = UserOrderDefaultComments(ctx)
			//if err != nil {
			//	panic(err)
			//}
			//管理后台路由组
			s.Group("/backend", func(group *ghttp.RouterGroup) {
				group.Middleware(
					service.Middleware().CORS,
					service.Middleware().Ctx,
					service.Middleware().ResponseHandler,
				)
				//不需要登录的路由组绑定
				group.Bind(
					controller.Admin.Create,       // 管理员
					controller.Login.Login,        // 登录
					controller.Login.RefreshToken, // 刷新Token
					controller.Data,               // 数据大屏相关
				)
				//需要登录的路由组绑定（JWT 中间件 + 权限校验）
				group.Group("/", func(group *ghttp.RouterGroup) {
					group.Middleware(
						service.Middleware().Auth,
						service.Middleware().PermissionCheck,
					)
					group.Bind(
						controller.Login.Logout,    // 登出
						controller.Rotation.Create, // 轮播图-创建
						controller.Rotation.Delete, // 轮播图-删除
						controller.Rotation.Update, // 轮播图-更新
						controller.Rotation.List,   // 轮播图-列表
						controller.Category,        // 商品分类
						controller.Role,            // 角色
						controller.Permission,      // 权限
						controller.Admin.List,      // 管理员
						controller.Admin.Update,    // 管理员
						controller.Admin.Delete,    // 管理员
						controller.Admin.Info,      // 查询当前管理员信息
						controller.Position,        // 手工位
						controller.File,            //从0到1实现文件入库
						controller.Upload,          //实现可跨项目使用的文件上云工具类
						controller.Coupon,          //商品优惠券管理
						controller.UserCoupon,      //商品优惠券管理
						controller.Goods,           //商品管理
						controller.GoodsOptions,    //商品规格管理
						controller.Address,         //城市地址管理
						//这么写是为了避免前后端重复注册相同的路由和方法
						controller.Order.List,   //订单列表
						controller.Order.Detail, //订单详情
						backend.Article,         //文章管理&CMS
						backend.Comment,         //评论管理
						backend.Consignee,       //收货地址管理
						backend.User,            //用户管理（超管专属）
					//controller.SeckillBackend, //秒杀后台管理
					)
				})
			})
			//---------------------华丽的分割线-------------------
			//---------------------华丽的分割线-------------------
			// 前台项目使用 JWT
			//前台项目路由组
			//前台项目路由组
			//前台项目路由组
			s.Group("/frontend", func(group *ghttp.RouterGroup) {
				group.Middleware(
					service.Middleware().CORS,
					service.Middleware().Ctx,
					service.Middleware().ResponseHandler,
				)
				//不需要登录的路由组绑定
				group.Bind(
					controller.User.Register,         //用户注册
					controller.Goods,                 //商品
					controller.User.Login,            //用户登录
					controller.Rotation.ListFrontend, //前台轮播图列表（不需要登录）
					frontend.Category,                //前台商品分类（不需要登录）
				)
				//需要登录鉴权的路由组
				group.Group("/", func(group *ghttp.RouterGroup) {
					// 前台需要登录鉴权的接口，挂载带 Next 的 UserAuth JWT 中间件
					group.Middleware(
						service.Middleware().UserAuth,
					)
					//需要登录鉴权的接口放到这里
					group.Bind(
						controller.User.Info,           //当前登录用户的信息
						controller.User.UpdatePassword, //当前用户修改密码
						controller.User.Logout,         //用户登出
						controller.Collection,          //收藏
						controller.Praise,              //收藏
						controller.Comment,             //评论
						controller.Cart,                //购物车
						controller.OrderGoodsComments,  //订单评价
						frontend.Order,                 //订单相关操作（列表、详情、支付、取消、确认收货等）
						frontend.Article,               //文章 @author自愚自乐
						frontend.Refund,                //售后 @author自愚自乐
						frontend.Consignee,             //收货地址管理
						//controller.SeckillFrontend,     //秒杀前台
					)
				})
			})

			// 添加公开的秒杀路由（不需要登录）
			//s.Group("/seckill", func(group *ghttp.RouterGroup) {
			//	group.Middleware(
			//		service.Middleware().CORS,
			//		service.Middleware().Ctx,
			//		service.Middleware().ResponseHandler,
			//	)
			//	group.Bind(
			//		controller.SeckillFrontend.List,   // 秒杀商品列表
			//		controller.SeckillFrontend.Detail, // 秒杀商品详情
			//	)
			//})

			s.SetPort(8199) //设置端口
			s.Run()
			return nil
		},
	}
)
