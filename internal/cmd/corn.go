package cmd

import (
	"context"
	"bit303_shop/internal/consts"
	"bit303_shop/internal/dao"
	"bit303_shop/internal/model/entity"
	"bit303_shop/utility"
	"time"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcron"
)

func UserOrderDefaultComments(ctx context.Context) (err error) {
	//每分钟执行一次
	_, err = gcron.Add(ctx, "0 */1 * * * *", func(ctx context.Context) {
		g.Log().Info(ctx, "开始执行订单默认评价定时任务...")

		// 强制重新初始化数据库连接
		// 确保使用正确的数据库配置
		gdb.SetConfig(gdb.Config{
			"default": gdb.ConfigGroup{
				gdb.ConfigNode{
					Type:  "mysql",
					Host:  "127.0.0.1",
					Port:  "3306",
					User:  "root",
					Pass:  "111111",
					Name:  "shop",
					Debug: true,
				},
			},
		})

		// 重新获取数据库连接
		db := g.DB()

		// 尝试执行一个简单查询确认连接正常
		_, testErr := db.Query(ctx, "SELECT 1")
		if testErr != nil {
			g.Log().Error(ctx, "数据库连接测试失败:", testErr)
			return
		}
		g.Log().Info(ctx, "数据库连接测试成功")

		// 检查数据库连接
		dbConfig := db.GetConfig()
		g.Log().Info(ctx, "当前数据库配置:", dbConfig)

		condition := g.Map{
			dao.OrderInfo.Columns().Status: 4,
		}
		// 使用我们创建的数据库连接
		model := db.Model(dao.OrderInfo.Table()).Where(condition)
		minTime := utility.TimeStampToDateTime(time.Now().Unix() - consts.UserOrderDefaultCommentsTime)
		count, err := model.Where("updated_at <=?", minTime).Count()
		if err != nil {
			g.Log().Error(ctx, "查询待评价订单失败:", err)
			return
		}

		g.Log().Info(ctx, "查询到", count, "个待自动评价的订单")

		if count > 0 {
			var orderList []entity.OrderInfo
			err = db.Model(dao.OrderInfo.Table()).Where(condition).Scan(&orderList)
			if err != nil {
				g.Log().Error(ctx, "获取订单列表失败:", err)
				return
			}
			for _, order := range orderList {
				//新增评价
				orderGoods := entity.OrderGoodsInfo{}
				err := db.Model(dao.OrderGoodsInfo.Table()).Where(dao.OrderGoodsInfo.Columns().OrderId, order.Id).Scan(&orderGoods)
				if err != nil {
					g.Log().Error(ctx, "获取订单商品信息失败:", err)
					continue
				}
				data := g.Map{
					dao.OrderGoodsCommentsInfo.Columns().OrderId:        order.Id,
					dao.OrderGoodsCommentsInfo.Columns().GoodsId:        orderGoods.GoodsId,
					dao.OrderGoodsCommentsInfo.Columns().GoodsOptionsId: orderGoods.GoodsOptionsId,
					dao.OrderGoodsCommentsInfo.Columns().Content:        consts.UserOrderDefaultComments,
				}
				_, err = db.Model(dao.OrderGoodsCommentsInfo.Table()).Data(data).InsertAndGetId()
				if err != nil {
					g.Log().Error(ctx, "新增订单评价失败:", err)
					continue
				}
				//更新订单状态
				in := g.Map{
					dao.OrderInfo.Columns().Status: consts.UserOrderStatus,
				}
				_, err = model.WherePri(order.Id).Data(in).Update()
				if err != nil {
					g.Log().Error(ctx, "更新订单状态失败:", err)
					continue
				}

				g.Log().Info(ctx, "订单", order.Id, "自动评价成功")
			}
		}
		return
	}, "UserOrderDefaultComments")
	return err
}
