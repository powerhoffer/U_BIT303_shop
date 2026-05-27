package init

import (
	"context"
	"fmt"
	"bit303_shop/internal/dao"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcmd"
)

var (
	Init = gcmd.Command{
		Name:  "init",
		Usage: "init",
		Brief: "初始化秒杀系统数据",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			// 1. 获取所有秒杀商品
			goods, err := dao.SeckillGoods.Ctx(ctx).All()
			if err != nil {
				return err
			}

			// 2. 初始化Redis库存
			for _, item := range goods {
				goodsId := item["goods_id"].Int()
				goodsOptionsId := item["goods_options_id"].Int()
				seckillStock := item["seckill_stock"].Int()
				
				stockKey := fmt.Sprintf("seckill:stock:%d:%d", goodsId, goodsOptionsId)
				_, err = g.Redis().Do(ctx, "SET", stockKey, seckillStock)
				if err != nil {
					return err
				}
				g.Log().Info(ctx, fmt.Sprintf("初始化商品 %d:%d 库存: %d",
					goodsId, goodsOptionsId, seckillStock))
			}

			// 3. 创建Kafka主题
			// 注意：这部分需要在Kafka中手动创建，或者使用Kafka管理工具

			g.Log().Info(ctx, "初始化完成")
			return nil
		},
	}
)
