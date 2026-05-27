package utility

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"bit303_shop/internal/dao"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/util/gconv"
)

// SeckillCacheKeys 秒杀系统使用的缓存键
const (
	// 商品库存缓存键前缀
	SeckillGoodsStockPrefix = "seckill:goods:stock:"
	// 商品信息缓存键前缀
	SeckillGoodsInfoPrefix = "seckill:goods:info:"
	// 用户参与秒杀记录前缀
	SeckillUserRecordPrefix = "seckill:user:record:"
	// 秒杀订单处理状态前缀
	SeckillOrderStatusPrefix = "seckill:order:status:"
	// 秒杀限流计数器前缀
	SeckillRateLimitPrefix = "seckill:ratelimit:"
)

// PreloadSeckillCache 预加载秒杀缓存
// 在系统启动时调用，将商品库存预加载到Redis中
func PreloadSeckillCache(ctx context.Context) error {
	g.Log().Info(ctx, "开始预加载秒杀缓存...")

	// 查询所有商品
	goodsList, err := dao.GoodsInfo.Ctx(ctx).
		Where("stock > 0").
		WhereIn("status", g.Slice{1, 2}). // 上架和推荐状态
		All()

	if err != nil {
		g.Log().Error(ctx, "预加载秒杀缓存失败:", err)
		return err
	}

	// 预加载到Redis
	for _, goods := range goodsList {
		// 设置商品库存缓存
		stockKey := SeckillGoodsStockPrefix + gconv.String(goods["id"])
		_, err = g.Redis().Do(ctx, "SETEX", stockKey, 86400, gconv.Int(goods["stock"])) // 24小时过期
		if err != nil {
			g.Log().Error(ctx, "设置商品库存缓存失败:", err)
			continue
		}

		// 设置商品信息缓存
		infoKey := SeckillGoodsInfoPrefix + gconv.String(goods["id"])
		// 将商品数据序列化为JSON
		goodsJson, err := json.Marshal(goods)
		if err != nil {
			g.Log().Error(ctx, "序列化商品信息失败:", err)
			continue
		}
		_, err = g.Redis().Do(ctx, "SETEX", infoKey, 86400, string(goodsJson)) // 24小时过期
		if err != nil {
			g.Log().Error(ctx, "设置商品信息缓存失败:", err)
			continue
		}
	}

	g.Log().Info(ctx, fmt.Sprintf("秒杀缓存预加载完成，共加载%d个商品", len(goodsList)))
	return nil
}

// DeductCacheStock 扣减缓存中的库存
// 原子性地扣减缓存中的库存，如果库存不足则返回错误
func DeductCacheStock(ctx context.Context, goodsId uint, count int) error {
	stockKey := SeckillGoodsStockPrefix + gconv.String(goodsId)

	// 使用Redis的DECRBY原子操作扣减库存
	result, err := g.Redis().Do(ctx, "EVAL", `
		local stock = redis.call('get', KEYS[1])
		if not stock or tonumber(stock) < tonumber(ARGV[1]) then
			return 0
		end
		return redis.call('decrby', KEYS[1], ARGV[1])
	`, 1, stockKey, count)

	if err != nil {
		return err
	}

	if gconv.Int(result) == 0 {
		// 库存不足，恢复缓存
		return errors.New("库存不足")
	}

	return nil
}

// RevertCacheStock 归还缓存中的库存
// 当订单处理失败时，需要将扣减的库存归还
func RevertCacheStock(ctx context.Context, goodsId uint, count int) error {
	stockKey := SeckillGoodsStockPrefix + gconv.String(goodsId)

	// 使用Redis的INCRBY原子操作增加库存
	_, err := g.Redis().Do(ctx, "INCRBY", stockKey, count)
	return err
}

// CheckUserSeckillEligibility 检查用户是否有资格参与秒杀
// 防止同一用户重复购买，或者参与次数超限
func CheckUserSeckillEligibility(ctx context.Context, userId, goodsId uint) (bool, error) {
	// 检查用户是否已经参与过此商品的秒杀
	recordKey := SeckillUserRecordPrefix + gconv.String(userId) + ":" + gconv.String(goodsId)

	// 使用Redis直接检查
	exists, err := g.Redis().Do(ctx, "EXISTS", recordKey)
	if err != nil {
		return false, err
	}

	if gconv.Int(exists) > 0 {
		return false, errors.New("您已经参与过此商品的秒杀")
	}

	// 检查用户今日秒杀次数是否超限
	limitKey := SeckillRateLimitPrefix + "user:" + gconv.String(userId) + ":" + time.Now().Format("20060102")
	count, err := g.Redis().Do(ctx, "GET", limitKey)
	if err != nil {
		return false, err
	}

	// 默认每用户每日限购5次
	if count != nil && gconv.Int(count) >= 5 {
		return false, errors.New("您今日的秒杀次数已达上限")
	}

	return true, nil
}

// RecordUserSeckill 记录用户参与秒杀
// 记录用户参与秒杀的信息，用于防止重复购买
func RecordUserSeckill(ctx context.Context, userId, goodsId uint) error {
	recordKey := SeckillUserRecordPrefix + gconv.String(userId) + ":" + gconv.String(goodsId)

	// 记录用户参与秒杀，过期时间24小时
	_, err := g.Redis().Do(ctx, "SETEX", recordKey, 86400, 1) // 24小时过期
	if err != nil {
		return err
	}

	// 增加用户今日秒杀次数
	limitKey := SeckillRateLimitPrefix + "user:" + gconv.String(userId) + ":" + time.Now().Format("20060102")
	_, err = g.Redis().Do(ctx, "INCR", limitKey)
	if err != nil {
		return err
	}

	// 设置次数缓存过期时间为当天结束
	tomorrow := time.Now().Add(24 * time.Hour)
	midnight := time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), 0, 0, 0, 0, tomorrow.Location())
	ttl := midnight.Sub(time.Now())

	_, err = g.Redis().Do(ctx, "EXPIRE", limitKey, int(ttl.Seconds()))
	return err
}

// InitSeckillSystem 初始化秒杀系统
// 在系统启动时调用，进行秒杀系统的初始化
func InitSeckillSystem() {
	ctx := gctx.New()
	g.Log().Info(ctx, "初始化秒杀系统...")

	// 预加载商品缓存
	if err := PreloadSeckillCache(ctx); err != nil {
		g.Log().Error(ctx, "预加载秒杀缓存失败:", err)
	}

	// 启动秒杀定时任务
	go scheduledTasks(ctx)

	g.Log().Info(ctx, "秒杀系统初始化完成")
}

// 定时任务
func scheduledTasks(ctx context.Context) {
	// 每小时同步一次库存数据
	syncStockTicker := time.NewTicker(1 * time.Hour)
	defer syncStockTicker.Stop()

	// 每天零点清理过期数据
	midnight := time.Now().Add(24 * time.Hour)
	midnight = time.Date(midnight.Year(), midnight.Month(), midnight.Day(), 0, 0, 0, 0, midnight.Location())
	cleanupTimer := time.NewTimer(midnight.Sub(time.Now()))
	defer cleanupTimer.Stop()

	for {
		select {
		case <-syncStockTicker.C:
			syncCacheStockToDatabase(ctx)
		case <-cleanupTimer.C:
			cleanupExpiredData(ctx)
			// 重置定时器到下一个零点
			nextMidnight := time.Now().Add(24 * time.Hour)
			nextMidnight = time.Date(nextMidnight.Year(), nextMidnight.Month(), nextMidnight.Day(), 0, 0, 0, 0, nextMidnight.Location())
			cleanupTimer.Reset(nextMidnight.Sub(time.Now()))
		}
	}
}

// 同步缓存库存到数据库
func syncCacheStockToDatabase(ctx context.Context) {
	g.Log().Info(ctx, "开始同步缓存库存到数据库...")

	// 查询所有商品
	goodsList, err := dao.GoodsInfo.Ctx(ctx).
		Where("stock > 0").
		WhereIn("status", g.Slice{1, 2}). // 上架和推荐状态
		All()

	if err != nil {
		g.Log().Error(ctx, "查询商品失败:", err)
		return
	}

	for _, goods := range goodsList {
		// 获取缓存中的库存
		stockKey := SeckillGoodsStockPrefix + gconv.String(goods["id"])
		cacheStock, err := g.Redis().Do(ctx, "GET", stockKey)
		if err != nil || cacheStock == nil {
			g.Log().Warning(ctx, fmt.Sprintf("获取商品[%d]缓存库存失败:", gconv.Int(goods["id"])), err)
			continue
		}

		// 如果缓存库存与数据库库存不一致，更新数据库
		if gconv.Int(cacheStock) != gconv.Int(goods["stock"]) {
			// 1. 更新商品表库存
			_, err = dao.GoodsInfo.Ctx(ctx).
				Data(g.Map{"stock": cacheStock}).
				WherePri(gconv.Int(goods["id"])).
				Update()

			if err != nil {
				g.Log().Error(ctx, fmt.Sprintf("更新商品[%d]库存失败:", gconv.Int(goods["id"])), err)
			} else {
				g.Log().Info(ctx, fmt.Sprintf("商品[%d]库存已从[%d]更新为[%d]",
					gconv.Int(goods["id"]),
					gconv.Int(goods["stock"]),
					gconv.Int(cacheStock)))
			}

			// 2. 查询并更新商品规格表
			goodsOptions, err := g.DB().Ctx(ctx).Model("goods_options_info").
				Where("goods_id", goods["id"]).
				All()

			if err != nil {
				g.Log().Error(ctx, fmt.Sprintf("查询商品[%d]规格失败:", gconv.Int(goods["id"])), err)
			} else if len(goodsOptions) > 0 {
				for _, option := range goodsOptions {
					optionStockKey := SeckillGoodsStockPrefix + gconv.String(goods["id"]) + ":" + gconv.String(option["id"])
					optionCacheStock, err := g.Redis().Do(ctx, "GET", optionStockKey)

					// 如果规格库存存在于Redis中
					if err == nil && optionCacheStock != nil {
						// 更新商品规格库存
						_, err = g.DB().Ctx(ctx).Model("goods_options_info").
							Data(g.Map{"stock": optionCacheStock}).
							WherePri(option["id"]).
							Update()

						if err != nil {
							g.Log().Error(ctx, fmt.Sprintf("更新商品规格[%d]库存失败:", gconv.Int(option["id"])), err)
						} else {
							g.Log().Info(ctx, fmt.Sprintf("商品规格[商品ID:%d, 规格ID:%d]库存已从[%d]更新为[%d]",
								gconv.Int(goods["id"]),
								gconv.Int(option["id"]),
								gconv.Int(option["stock"]),
								gconv.Int(optionCacheStock)))
						}
					} else {
						// 如果规格库存不在Redis中，使用商品主表的库存
						_, err = g.DB().Ctx(ctx).Model("goods_options_info").
							Data(g.Map{"stock": cacheStock}).
							WherePri(option["id"]).
							Update()

						if err != nil {
							g.Log().Error(ctx, fmt.Sprintf("更新商品规格[%d]库存失败:", gconv.Int(option["id"])), err)
						}
					}
				}
			}

			// 3. 查询并更新秒杀商品表中对应的商品
			seckillGoods, err := g.DB().Ctx(ctx).Model("seckill_goods").
				Where("goods_id", goods["id"]).
				All()

			if err != nil {
				g.Log().Error(ctx, fmt.Sprintf("查询秒杀商品[%d]失败:", gconv.Int(goods["id"])), err)
				continue
			}

			if len(seckillGoods) > 0 {
				for _, seckill := range seckillGoods {
					// 如果是带规格的秒杀商品
					if gconv.Int(seckill["goods_options_id"]) > 0 {
						optionId := gconv.Int(seckill["goods_options_id"])
						optionStockKey := SeckillGoodsStockPrefix + gconv.String(goods["id"]) + ":" + gconv.String(optionId)
						optionCacheStock, err := g.Redis().Do(ctx, "GET", optionStockKey)

						// 如果规格库存存在于Redis中
						if err == nil && optionCacheStock != nil {
							// 更新秒杀商品库存为规格库存
							_, err = g.DB().Ctx(ctx).Model("seckill_goods").
								Data(g.Map{"seckill_stock": optionCacheStock}).
								WherePri(seckill["id"]).
								Update()

							if err != nil {
								g.Log().Error(ctx, fmt.Sprintf("更新秒杀商品[%d]规格库存失败:", gconv.Int(seckill["id"])), err)
							} else {
								g.Log().Info(ctx, fmt.Sprintf("秒杀商品[ID:%d, 商品ID:%d, 规格ID:%d]库存已从[%d]更新为[%d]",
									gconv.Int(seckill["id"]),
									gconv.Int(goods["id"]),
									optionId,
									gconv.Int(seckill["seckill_stock"]),
									gconv.Int(optionCacheStock)))
							}
						} else {
							// 如果规格库存不在Redis中，使用商品主表的库存
							_, err = g.DB().Ctx(ctx).Model("seckill_goods").
								Data(g.Map{"seckill_stock": cacheStock}).
								WherePri(seckill["id"]).
								Update()

							if err != nil {
								g.Log().Error(ctx, fmt.Sprintf("更新秒杀商品[%d]库存失败:", gconv.Int(seckill["id"])), err)
							}
						}
					} else {
						// 更新秒杀商品库存
						_, err = g.DB().Ctx(ctx).Model("seckill_goods").
							Data(g.Map{"seckill_stock": cacheStock}).
							WherePri(seckill["id"]).
							Update()

						if err != nil {
							g.Log().Error(ctx, fmt.Sprintf("更新秒杀商品[%d]库存失败:", gconv.Int(seckill["id"])), err)
						} else {
							g.Log().Info(ctx, fmt.Sprintf("秒杀商品[ID:%d, 商品ID:%d]库存已从[%d]更新为[%d]",
								gconv.Int(seckill["id"]),
								gconv.Int(goods["id"]),
								gconv.Int(seckill["seckill_stock"]),
								gconv.Int(cacheStock)))
						}
					}
				}
			}
		}
	}

	g.Log().Info(ctx, "缓存库存同步到数据库完成")
}

// 清理过期数据
func cleanupExpiredData(ctx context.Context) {
	g.Log().Info(ctx, "开始清理过期数据...")

	// 清理过期的秒杀记录
	pattern := SeckillUserRecordPrefix + "*"
	_, err := g.Redis().Do(ctx, "SCAN", 0, "MATCH", pattern, "COUNT", 100)
	if err != nil {
		g.Log().Error(ctx, "清理过期数据失败:", err)
	}

	g.Log().Info(ctx, "过期数据清理完成")
}
