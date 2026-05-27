package seckill

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
)

// StockSegments 定义分段数（按用户ID取模分段）
const StockSegments = 16

// StockManager 库存管理器
type StockManager struct {
	// 先放置需要原子操作的64位字段，确保8字节对齐
	stats struct {
		hits      int64 // 缓存命中次数
		misses    int64 // 缓存未命中次数
		updates   int64 // 更新次数
		conflicts int64 // 冲突次数
	}

	// 然后是指针和映射类型
	localCache sync.Map          // 本地库存缓存 (商品ID -> 库存量)
	scripts    map[string]string // Redis Lua 脚本相关常量

	// 最后是32位字段和其他数据
	defaultStock int32                     // 在Redis中存储的库存默认值
	segmentLocks [StockSegments]sync.Mutex // 分段锁 (降低锁竞争)
}

// NewStockManager 创建库存管理器
func NewStockManager(defaultStock int32) *StockManager {
	return &StockManager{
		defaultStock: defaultStock,
		scripts: map[string]string{
			"deduct": `
				-- 扣减库存Lua脚本
				local stockKey = KEYS[1]
				local quantity = tonumber(ARGV[1])
				local current = tonumber(redis.call('GET', stockKey) or '0')
				
				-- 检查库存是否足够
				if current < quantity then
					return -1
				end
				
				-- 扣减库存并返回剩余量
				local remain = current - quantity
				redis.call('SET', stockKey, remain)
				return remain
			`,
			"add": `
				-- 增加库存Lua脚本
				local stockKey = KEYS[1]
				local quantity = tonumber(ARGV[1])
				local current = tonumber(redis.call('GET', stockKey) or '0')
				
				-- 增加库存并返回当前量
				local remain = current + quantity
				redis.call('SET', stockKey, remain)
				return remain
			`,
			"check": `
				-- 检查库存Lua脚本
				local stockKey = KEYS[1]
				local minAmount = tonumber(ARGV[1] or '1')
				local current = tonumber(redis.call('GET', stockKey) or '0')
				
				-- 检查库存是否足够
				if current >= minAmount then
					return current
				else
					return -1
				end
			`,
		},
	}
}

// GetStockKey 获取商品库存在Redis中的键
func (sm *StockManager) GetStockKey(goodsId, optionId int32) string {
	if optionId > 0 {
		return fmt.Sprintf("seckill:stock:%d:%d", goodsId, optionId)
	}
	return fmt.Sprintf("seckill:stock:%d", goodsId)
}

// InitStock 初始化商品库存
func (sm *StockManager) InitStock(ctx context.Context, goodsId, optionId int32, stock int32) error {
	stockKey := sm.GetStockKey(goodsId, optionId)

	// 检查是否已经存在库存
	exists, err := g.Redis().Do(ctx, "EXISTS", stockKey)
	if err != nil {
		return err
	}

	// 如果不存在，设置初始库存
	if gconv.Int(exists) == 0 {
		_, err = g.Redis().Do(ctx, "SET", stockKey, stock)
		if err != nil {
			return err
		}

		// 更新本地缓存
		cacheKey := sm.getCacheKey(goodsId, optionId)
		sm.localCache.Store(cacheKey, stock)
	}

	return nil
}

// getCacheKey 生成本地缓存键
func (sm *StockManager) getCacheKey(goodsId, optionId int32) string {
	if optionId > 0 {
		return fmt.Sprintf("%d:%d", goodsId, optionId)
	}
	return fmt.Sprintf("%d", goodsId)
}

// getSegmentLock 获取分段锁
func (sm *StockManager) getSegmentLock(userId int64) *sync.Mutex {
	// 用户ID取模获取分段索引
	segment := userId % StockSegments
	if segment < 0 {
		segment = -segment
	}
	return &sm.segmentLocks[segment]
}

// CheckStock 检查库存（先查本地缓存，再查Redis）
func (sm *StockManager) CheckStock(ctx context.Context, goodsId, optionId int32) (int32, error) {
	cacheKey := sm.getCacheKey(goodsId, optionId)

	// 1. 先查本地缓存
	if value, ok := sm.localCache.Load(cacheKey); ok {
		stock := value.(int32)
		atomic.AddInt64(&sm.stats.hits, 1)
		// 如果本地缓存中库存大于0，直接返回
		if stock > 0 {
			return stock, nil
		}
	} else {
		atomic.AddInt64(&sm.stats.misses, 1)
	}

	// 2. 查询Redis
	stockKey := sm.GetStockKey(goodsId, optionId)
	script := sm.scripts["check"]

	// 执行Lua脚本检查库存
	result, err := g.Redis().Do(ctx, "EVAL", script, 1, stockKey, 1)
	if err != nil {
		return 0, err
	}

	stock := gconv.Int32(result)
	if stock >= 0 {
		// 更新本地缓存
		sm.localCache.Store(cacheKey, stock)
		return stock, nil
	}

	return 0, nil
}

// DeductStock 扣减库存（使用Redis Lua脚本保证原子性）
func (sm *StockManager) DeductStock(ctx context.Context, userId int64, goodsId, optionId, quantity int32) (int32, error) {
	// 获取分段锁（按用户ID分段）
	lock := sm.getSegmentLock(userId)
	lock.Lock()
	defer lock.Unlock()

	// 执行Redis库存扣减
	stockKey := sm.GetStockKey(goodsId, optionId)
	script := sm.scripts["deduct"]

	// 执行Lua脚本原子扣减库存
	result, err := g.Redis().Do(ctx, "EVAL", script, 1, stockKey, quantity)
	if err != nil {
		atomic.AddInt64(&sm.stats.conflicts, 1)
		return 0, err
	}

	remain := gconv.Int32(result)
	if remain >= 0 {
		// 更新本地缓存
		cacheKey := sm.getCacheKey(goodsId, optionId)
		sm.localCache.Store(cacheKey, remain)
		atomic.AddInt64(&sm.stats.updates, 1)
		return remain, nil
	}

	return 0, fmt.Errorf("insufficient stock")
}

// AddStock 增加库存
func (sm *StockManager) AddStock(ctx context.Context, goodsId, optionId, quantity int32) (int32, error) {
	// 执行Redis库存增加
	stockKey := sm.GetStockKey(goodsId, optionId)
	script := sm.scripts["add"]

	// 执行Lua脚本原子增加库存
	result, err := g.Redis().Do(ctx, "EVAL", script, 1, stockKey, quantity)
	if err != nil {
		return 0, err
	}

	remain := gconv.Int32(result)

	// 更新本地缓存
	cacheKey := sm.getCacheKey(goodsId, optionId)
	sm.localCache.Store(cacheKey, remain)
	atomic.AddInt64(&sm.stats.updates, 1)

	return remain, nil
}

// GetCachedStock 从本地缓存获取库存
func (sm *StockManager) GetCachedStock(goodsId, optionId int32) (int32, bool) {
	cacheKey := sm.getCacheKey(goodsId, optionId)
	if value, ok := sm.localCache.Load(cacheKey); ok {
		atomic.AddInt64(&sm.stats.hits, 1)
		return value.(int32), true
	}
	atomic.AddInt64(&sm.stats.misses, 1)
	return 0, false
}

// SyncStockToDatabase 将Redis中的库存同步到数据库
func (sm *StockManager) SyncStockToDatabase(ctx context.Context, goodsId, optionId int32) error {
	stockKey := sm.GetStockKey(goodsId, optionId)

	// 从Redis获取当前库存
	result, err := g.Redis().Do(ctx, "GET", stockKey)
	if err != nil {
		g.Log().Errorf(ctx, "获取Redis库存失败[%d:%d]: %v", goodsId, optionId, err)
		return fmt.Errorf("获取Redis库存失败: %v", err)
	}

	// 检查Redis中是否存在该库存键
	if result == nil {
		g.Log().Warningf(ctx, "Redis中不存在库存键[%d:%d]，跳过同步", goodsId, optionId)
		return fmt.Errorf("Redis库存不存在")
	}

	// 明确使用Int32转换，保持与数据库字段类型一致
	currentStock := gconv.Int32(result)
	g.Log().Infof(ctx, "准备同步库存到数据库[%d:%d]: Redis库存=%d", goodsId, optionId, currentStock)

	// 1. 更新秒杀商品表中的库存，显式使用事务确保数据一致性
	tx, err := g.DB().Ctx(ctx).Begin(ctx)
	if err != nil {
		g.Log().Errorf(ctx, "开始事务失败: %v", err)
		return fmt.Errorf("开始事务失败: %v", err)
	}

	// 2. 更新秒杀商品表
	res, err := tx.Model("seckill_goods").
		Data(g.Map{
			"seckill_stock": currentStock,
			"updated_at":    gtime.Now(),
		}).
		Where("goods_id", goodsId).
		Where("goods_options_id", optionId).
		Update()

	if err != nil {
		tx.Rollback()
		g.Log().Errorf(ctx, "更新秒杀商品表库存失败[%d:%d]: %v", goodsId, optionId, err)
		return fmt.Errorf("更新秒杀商品表库存失败: %v", err)
	}

	// 检查是否有记录被更新
	affected, err := res.RowsAffected()
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("获取影响行数失败: %v", err)
	}

	if affected == 0 {
		g.Log().Warningf(ctx, "未找到秒杀商品记录[%d:%d]，同步库存跳过", goodsId, optionId)
	} else {
		g.Log().Infof(ctx, "秒杀商品表库存已更新[%d:%d]: %d (影响行数: %d)", goodsId, optionId, currentStock, affected)
	}

	// 2. 如果是商品规格，更新商品选项表
	if optionId > 0 {
		_, err = tx.Model("goods_options_info").
			Data(g.Map{
				"stock":      currentStock,
				"updated_at": gtime.Now(),
			}).
			Where("id", optionId).
			Update()

		if err != nil {
			tx.Rollback()
			g.Log().Warningf(ctx, "更新商品规格[%d]库存失败: %v", optionId, err)
			return fmt.Errorf("更新商品规格库存失败: %v", err)
		}

		g.Log().Infof(ctx, "商品规格[%d]库存已更新为: %d", optionId, currentStock)
	}

	// 3. 始终更新商品主表库存
	_, err = tx.Model("goods_info").
		Data(g.Map{
			"stock":      currentStock,
			"updated_at": gtime.Now(),
		}).
		Where("id", goodsId).
		Update()

	if err != nil {
		tx.Rollback()
		g.Log().Warningf(ctx, "更新商品[%d]库存失败: %v", goodsId, err)
		return fmt.Errorf("更新商品主表库存失败: %v", err)
	}

	g.Log().Infof(ctx, "商品[%d]库存已更新为: %d", goodsId, currentStock)

	// 提交事务
	if err = tx.Commit(); err != nil {
		g.Log().Errorf(ctx, "提交事务失败: %v", err)
		return fmt.Errorf("提交事务失败: %v", err)
	}

	g.Log().Infof(ctx, "商品[%d:%d]库存已成功同步到数据库, 当前库存: %d", goodsId, optionId, currentStock)
	return nil
}

// BatchSyncStockToDatabase 批量同步库存到数据库
func (sm *StockManager) BatchSyncStockToDatabase(ctx context.Context, items []struct {
	GoodsId  int32
	OptionId int32
}) error {
	if len(items) == 0 {
		return nil
	}

	g.Log().Info(ctx, "开始批量同步库存，商品数量:", len(items))

	// 获取所有商品的Redis库存
	var goodsStocks = make(map[string]int32) // 修改为int32确保类型一致
	var stockKeys = make([]interface{}, 0, len(items))
	var keyToGoodsMap = make(map[string]struct {
		GoodsId  int32
		OptionId int32
	})

	// 检查Redis是否可用
	redis := g.Redis()
	if redis == nil {
		g.Log().Warning(ctx, "Redis未配置，直接从数据库同步商品库存")
		return sm.syncStockFromDatabase(ctx, items)
	}

	// 测试Redis连接
	_, pingErr := redis.Do(ctx, "PING")
	if pingErr != nil {
		g.Log().Warningf(ctx, "Redis连接测试失败: %v, 将从数据库同步", pingErr)
		return sm.syncStockFromDatabase(ctx, items)
	}

	// 准备所有需要的Key
	for _, item := range items {
		key := sm.GetStockKey(item.GoodsId, item.OptionId)
		stockKeys = append(stockKeys, key)
		keyToGoodsMap[key] = item
	}

	// 尝试使用MGET批量获取所有库存
	results, err := redis.Do(ctx, "MGET", stockKeys...)
	if err != nil {
		g.Log().Errorf(ctx, "批量获取Redis库存失败: %v, 将从数据库同步", err)
		return sm.syncStockFromDatabase(ctx, items)
	}

	// 解析结果
	resultSlice := results.Slice()
	if resultSlice == nil || len(resultSlice) == 0 {
		g.Log().Warning(ctx, "Redis返回结果为空，将从数据库同步")
		return sm.syncStockFromDatabase(ctx, items)
	}

	for i, result := range resultSlice {
		if i >= len(stockKeys) {
			continue
		}

		// 检查结果是否为空
		if result == nil {
			g.Log().Warningf(ctx, "Redis中不存在库存键: %v", stockKeys[i])
			continue
		}

		key := stockKeys[i].(string)
		goodsStocks[key] = gconv.Int32(result) // 使用Int32类型确保类型一致
	}

	if len(goodsStocks) == 0 {
		g.Log().Warning(ctx, "没有获取到任何有效的库存数据，将从数据库同步")
		return sm.syncStockFromDatabase(ctx, items)
	}

	g.Log().Infof(ctx, "成功从Redis获取了%d个商品的库存数据", len(goodsStocks))

	// 检查数据库连接
	if _, err := g.DB().GetOne(ctx, "SELECT 1"); err != nil {
		g.Log().Errorf(ctx, "数据库连接测试失败: %v, 无法同步库存", err)
		return fmt.Errorf("数据库连接失败: %v", err)
	}

	// 每个商品单独同步，使用事务确保每个商品的更新是原子的
	// 这种方式虽然比一次性批量更新慢，但更安全可靠
	successCount := 0
	failCount := 0

	for key, stock := range goodsStocks {
		item := keyToGoodsMap[key]
		goodsId := item.GoodsId
		optionId := item.OptionId

		// 开始事务
		tx, err := g.DB().Ctx(ctx).Begin(ctx)
		if err != nil {
			g.Log().Errorf(ctx, "开始事务失败[%d:%d]: %v", goodsId, optionId, err)
			failCount++
			continue // 继续处理其他商品
		}

		// 更新秒杀商品表
		res, err := tx.Model("seckill_goods").
			Data(g.Map{
				"seckill_stock": stock,
				"updated_at":    gtime.Now(),
			}).
			Where("goods_id", goodsId).
			Where("goods_options_id", optionId).
			Update()

		if err != nil {
			tx.Rollback()
			g.Log().Errorf(ctx, "更新秒杀商品表库存失败[%d:%d]: %v", goodsId, optionId, err)
			failCount++
			continue // 继续处理其他商品
		}

		// 检查是否有记录被更新
		affected, err := res.RowsAffected()
		if err != nil {
			tx.Rollback()
			g.Log().Errorf(ctx, "获取影响行数失败[%d:%d]: %v", goodsId, optionId, err)
			failCount++
			continue
		}

		if affected == 0 {
			g.Log().Warningf(ctx, "未找到秒杀商品记录[%d:%d]，同步跳过", goodsId, optionId)
			tx.Rollback() // 没有找到记录，回滚事务
			failCount++
			continue
		}

		// 如果是商品规格，更新商品选项表
		if optionId > 0 {
			_, err = tx.Model("goods_options_info").
				Data(g.Map{
					"stock":      stock,
					"updated_at": gtime.Now(),
				}).
				Where("id", optionId).
				Update()

			if err != nil {
				tx.Rollback()
				g.Log().Warningf(ctx, "更新商品规格[%d]库存失败: %v", optionId, err)
				failCount++
				continue
			}
		}

		// 更新商品主表
		_, err = tx.Model("goods_info").
			Data(g.Map{
				"stock":      stock,
				"updated_at": gtime.Now(),
			}).
			Where("id", goodsId).
			Update()

		if err != nil {
			tx.Rollback()
			g.Log().Warningf(ctx, "更新商品[%d]库存失败: %v", goodsId, err)
			failCount++
			continue
		}

		// 提交事务
		if err = tx.Commit(); err != nil {
			g.Log().Errorf(ctx, "提交事务失败[%d:%d]: %v", goodsId, optionId, err)
			failCount++
			continue
		}

		g.Log().Infof(ctx, "商品[%d:%d]库存同步成功，当前库存: %d", goodsId, optionId, stock)
		successCount++
	}

	g.Log().Infof(ctx, "批量同步库存完成，成功: %d, 失败: %d", successCount, failCount)

	if failCount > 0 {
		return fmt.Errorf("部分商品库存同步失败(%d/%d)", failCount, len(goodsStocks))
	}

	return nil
}

// syncStockFromDatabase 从数据库直接同步商品库存(Redis降级方案)
func (sm *StockManager) syncStockFromDatabase(ctx context.Context, items []struct {
	GoodsId  int32
	OptionId int32
}) error {
	// 这是Redis不可用时的降级方案
	g.Log().Info(ctx, "开始从数据库直接同步库存数据")

	// 检查数据库连接
	if _, err := g.DB().GetOne(ctx, "SELECT 1"); err != nil {
		g.Log().Errorf(ctx, "数据库连接测试失败: %v, 无法同步库存", err)
		return fmt.Errorf("数据库连接失败: %v", err)
	}

	successCount := 0
	failCount := 0

	for _, item := range items {
		// 直接从数据库读取秒杀商品当前库存
		var stockValue int32 = 0
		record, err := g.DB().Model("seckill_goods").
			Fields("seckill_stock").
			Where("goods_id", item.GoodsId).
			Where("goods_options_id", item.OptionId).
			One()

		if err != nil {
			g.Log().Errorf(ctx, "从数据库读取商品[%d:%d]库存失败: %v",
				item.GoodsId, item.OptionId, err)
			failCount++
			continue
		}

		if record == nil {
			g.Log().Warningf(ctx, "商品[%d:%d]不存在于秒杀商品表中",
				item.GoodsId, item.OptionId)
			failCount++
			continue
		}

		stockValue = gconv.Int32(record["seckill_stock"])

		// 尝试更新本地缓存
		cacheKey := sm.getCacheKey(item.GoodsId, item.OptionId)
		sm.localCache.Store(cacheKey, stockValue)

		// 尝试更新Redis缓存(如果Redis已恢复)
		redis := g.Redis()
		if redis != nil {
			// 测试Redis连接是否已恢复
			_, pingErr := redis.Do(ctx, "PING")
			if pingErr == nil {
				stockKey := sm.GetStockKey(item.GoodsId, item.OptionId)
				_, redisErr := redis.Do(ctx, "SET", stockKey, stockValue)
				if redisErr != nil {
					g.Log().Warningf(ctx, "更新Redis缓存失败[%d:%d]: %v",
						item.GoodsId, item.OptionId, redisErr)
				} else {
					g.Log().Infof(ctx, "成功将商品[%d:%d]库存更新到Redis: %d",
						item.GoodsId, item.OptionId, stockValue)
				}
			}
		}

		g.Log().Infof(ctx, "商品[%d:%d]库存同步成功，当前库存: %d",
			item.GoodsId, item.OptionId, stockValue)
		successCount++
	}

	g.Log().Infof(ctx, "数据库同步库存完成，成功: %d, 失败: %d", successCount, failCount)

	if failCount > 0 {
		return fmt.Errorf("部分商品库存同步失败(%d/%d)", failCount, len(items))
	}

	return nil
}

// BatchDeductStock 批量扣减库存（适用于购物车结算等场景）
func (sm *StockManager) BatchDeductStock(ctx context.Context, userId int64, items []struct {
	GoodsId  int32
	OptionId int32
	Quantity int32
}) (map[string]int32, error) {
	if len(items) == 0 {
		return nil, nil
	}

	// 获取分段锁（按用户ID分段）
	lock := sm.getSegmentLock(userId)
	lock.Lock()
	defer lock.Unlock()

	// 准备批量扣减
	result := make(map[string]int32)

	// 1. 逐个处理每个商品的库存扣减
	for _, item := range items {
		stockKey := sm.GetStockKey(item.GoodsId, item.OptionId)

		// 执行Lua脚本原子扣减库存
		script := sm.scripts["deduct"]
		deductResult, err := g.Redis().Do(ctx, "EVAL", script, 1, stockKey, item.Quantity)
		if err != nil {
			atomic.AddInt64(&sm.stats.conflicts, 1)

			// 扣减失败，回滚已处理的商品
			for key, _ := range result {
				// 解析库存键
				parts := strings.Split(key, ":")
				if len(parts) >= 3 {
					goodsId := gconv.Int32(parts[2])
					var optionId int32 = 0
					if len(parts) >= 4 {
						optionId = gconv.Int32(parts[3])
					}

					// 根据已经扣减的数量计算需要增加的数量
					for _, it := range items {
						if it.GoodsId == goodsId && it.OptionId == optionId {
							// 增加回对应的库存数量
							sm.AddStock(ctx, goodsId, optionId, it.Quantity)
							break
						}
					}
				}
			}

			return nil, err
		}

		remain := gconv.Int32(deductResult)
		if remain < 0 {
			// 扣减失败，回滚已处理的商品
			for key := range result {
				// 解析库存键
				parts := strings.Split(key, ":")
				if len(parts) >= 3 {
					goodsId := gconv.Int32(parts[2])
					var optionId int32 = 0
					if len(parts) >= 4 {
						optionId = gconv.Int32(parts[3])
					}

					// 根据已经扣减的数量计算需要增加的数量
					for _, it := range items {
						if it.GoodsId == goodsId && it.OptionId == optionId {
							// 增加回对应的库存数量
							sm.AddStock(ctx, goodsId, optionId, it.Quantity)
							break
						}
					}
				}
			}

			return nil, fmt.Errorf("insufficient stock for item: %d-%d", item.GoodsId, item.OptionId)
		}

		// 记录结果
		result[stockKey] = remain

		// 更新本地缓存
		cacheKey := sm.getCacheKey(item.GoodsId, item.OptionId)
		sm.localCache.Store(cacheKey, remain)
	}

	atomic.AddInt64(&sm.stats.updates, int64(len(items)))
	return result, nil
}

// BatchCheckStock 批量检查多个商品的库存
func (sm *StockManager) BatchCheckStock(ctx context.Context, items []struct {
	GoodsId  int32
	OptionId int32
	Quantity int32
}) (map[string]bool, error) {
	if len(items) == 0 {
		return nil, nil
	}

	// 结果映射
	result := make(map[string]bool)

	// 本地缓存检查结果
	cachedResults := make(map[string]int32)
	missedItems := make([]struct {
		GoodsId  int32
		OptionId int32
		Quantity int32
	}, 0)

	// 1. 先查本地缓存
	for _, item := range items {
		cacheKey := sm.getCacheKey(item.GoodsId, item.OptionId)
		if value, ok := sm.localCache.Load(cacheKey); ok {
			stock := value.(int32)
			atomic.AddInt64(&sm.stats.hits, 1)
			cachedResults[cacheKey] = stock
			// 检查库存是否足够
			result[cacheKey] = stock >= item.Quantity
		} else {
			atomic.AddInt64(&sm.stats.misses, 1)
			missedItems = append(missedItems, item)
		}
	}

	// 如果全部命中缓存，直接返回
	if len(missedItems) == 0 {
		return result, nil
	}

	// 2. 使用Lua脚本批量检查库存
	script := sm.scripts["check"]

	for _, item := range missedItems {
		stockKey := sm.GetStockKey(item.GoodsId, item.OptionId)
		checkResult, err := g.Redis().Do(ctx, "EVAL", script, 1, stockKey, item.Quantity)

		if err != nil {
			// 查询错误，认为库存不足
			result[stockKey] = false
			continue
		}

		stock := gconv.Int(checkResult)
		hasStock := stock >= int(item.Quantity)
		result[stockKey] = hasStock

		// 更新本地缓存
		if stock > 0 {
			cacheKey := sm.getCacheKey(item.GoodsId, item.OptionId)
			sm.localCache.Store(cacheKey, int32(stock))
		}
	}

	return result, nil
}

// GetStats 获取统计数据
func (sm *StockManager) GetStats() map[string]int64 {
	return map[string]int64{
		"cache_hits":   atomic.LoadInt64(&sm.stats.hits),
		"cache_misses": atomic.LoadInt64(&sm.stats.misses),
		"updates":      atomic.LoadInt64(&sm.stats.updates),
		"conflicts":    atomic.LoadInt64(&sm.stats.conflicts),
	}
}

// Reset 重置统计数据
func (sm *StockManager) Reset() {
	atomic.StoreInt64(&sm.stats.hits, 0)
	atomic.StoreInt64(&sm.stats.misses, 0)
	atomic.StoreInt64(&sm.stats.updates, 0)
	atomic.StoreInt64(&sm.stats.conflicts, 0)
}
