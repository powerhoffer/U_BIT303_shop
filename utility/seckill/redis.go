package seckill

import (
	"context"
	"errors"
	"fmt"
	"bit303_shop/internal/consts"
	"time"

	"hash/fnv"

	"github.com/gogf/gf/v2/database/gredis"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/util/gconv"
)

const (
	// Redis键前缀
	KeyPrefixSeckillResult = "seckill:result:" // 秒杀结果缓存前缀
	KeyPrefixSeckillOrder  = "seckill:order:"  // 秒杀订单前缀
	KeyPrefixSeckillStock  = "seckill:stock:"  // 秒杀库存前缀
	KeyPrefixSeckillLock   = "seckill:lock:"   // 秒杀分布式锁前缀
	KeyPrefixBloomFilter   = "seckill:bloom:"  // 布隆过滤器前缀
	KeyPrefixOrderSent     = "seckill:sent:"   // 订单消息已发送标记前缀

	// Lua 脚本 - 分布式锁
	LuaScriptLock = `
		if redis.call('setnx', KEYS[1], ARGV[1]) == 1 then
			redis.call('expire', KEYS[1], ARGV[2])
			return 1
		else
			return 0
		end
	`

	// Lua 脚本 - 解锁
	LuaScriptUnlock = `
		if redis.call('get', KEYS[1]) == ARGV[1] then
			return redis.call('del', KEYS[1])
		else
			return 0
		end
	`

	// 默认锁超时时间 (秒)
	DefaultLockExpiry = 10
)

// 定义错误类型
var (
	ErrLockAcquireFailed = errors.New("failed to acquire lock")
	ErrLockReleaseFailed = errors.New("failed to release lock")
	ErrCacheMiss         = errors.New("cache miss")
)

// RedisKeys 返回秒杀系统使用的Redis键
type RedisKeys struct {
	StockKey           string // 库存键
	GoodsInfoKey       string // 商品信息键
	UserBoughtKey      string // 用户购买记录键
	SuccessCountKey    string // 成功计数键
	BloomFilterKey     string // 布隆过滤器键
	ProcessingFlagKey  string // 处理中标记键
	ResultKey          string // 结果键
	LockKey            string // 分布式锁键
	OrderProcessingKey string // 订单处理中键
	OrderCompletedKey  string // 订单完成键
}

// NewRedisKeys 创建Redis键
func NewRedisKeys(goodsId, goodsOptionsId, userId int64) *RedisKeys {
	return &RedisKeys{
		StockKey:           fmt.Sprintf("%s%d:%d", consts.SeckillStockPrefix, goodsId, goodsOptionsId),
		GoodsInfoKey:       fmt.Sprintf("%s%d:%d", consts.SeckillGoodsPrefix, goodsId, goodsOptionsId),
		UserBoughtKey:      fmt.Sprintf("%s%d:%d:%d", consts.SeckillUserBoughtPrefix, userId, goodsId, goodsOptionsId),
		SuccessCountKey:    fmt.Sprintf("%s%d:%d", consts.SeckillSuccessPrefix, goodsId, goodsOptionsId),
		BloomFilterKey:     fmt.Sprintf("seckill:bloom:%d:%d", goodsId, goodsOptionsId),
		ProcessingFlagKey:  fmt.Sprintf("seckill:processing:%d:%d:%d", userId, goodsId, goodsOptionsId),
		ResultKey:          fmt.Sprintf("%s%d:%d:%d", consts.SeckillResultPrefix, userId, goodsId, goodsOptionsId),
		LockKey:            fmt.Sprintf("%s%d:%d", consts.SeckillLockPrefix, goodsId, goodsOptionsId),
		OrderProcessingKey: fmt.Sprintf("seckill:order:processing:%d", userId),
		OrderCompletedKey:  fmt.Sprintf("seckill:order:completed:%d", userId),
	}
}

// GetRedisClient 获取Redis客户端
func GetRedisClient() *gredis.Redis {
	// 获取默认Redis实例，如果不存在则返回nil
	redis := g.Redis()
	if redis == nil {
		g.Log().Warning(gctx.New(), "Redis实例未配置或无法获取")
	}
	return redis
}

// AcquireLock 获取分布式锁
func AcquireLock(ctx context.Context, lockKey string, lockValue string, expiry int) (bool, error) {
	// 默认10秒过期
	if expiry <= 0 {
		expiry = DefaultLockExpiry
	}

	// 使用Lua脚本确保原子性操作
	result, err := g.Redis().Do(ctx, "EVAL", LuaScriptLock, 1, lockKey, lockValue, expiry)
	if err != nil {
		return false, err
	}

	// 成功获取锁
	if gconv.Int(result) == 1 {
		return true, nil
	}

	// 获取锁失败
	return false, nil
}

// ReleaseLock 释放分布式锁
func ReleaseLock(ctx context.Context, lockKey string, lockValue string) (bool, error) {
	// 使用Lua脚本确保原子性操作
	result, err := g.Redis().Do(ctx, "EVAL", LuaScriptUnlock, 1, lockKey, lockValue)
	if err != nil {
		return false, err
	}

	// 成功释放锁
	if gconv.Int(result) == 1 {
		return true, nil
	}

	// 锁已过期或被其他客户端获取
	return false, nil
}

// TryWithLock 尝试获取锁执行函数
func TryWithLock(ctx context.Context, lockKey string, fn func() error) error {
	// 生成锁的值 (随机值)
	lockValue := fmt.Sprintf("%d", time.Now().UnixNano())

	// 尝试获取锁
	acquired, err := AcquireLock(ctx, lockKey, lockValue, DefaultLockExpiry)
	if err != nil {
		return err
	}

	if !acquired {
		return ErrLockAcquireFailed
	}

	// 确保在函数返回时释放锁
	defer func() {
		_, _ = ReleaseLock(ctx, lockKey, lockValue)
	}()

	// 执行传入的函数
	return fn()
}

// GetCache 从Redis获取缓存
func GetCache(ctx context.Context, key string, value interface{}) error {
	redis := GetRedisClient()
	if redis == nil {
		return errors.New("Redis未配置")
	}

	result, err := redis.Do(ctx, "GET", key)
	if err != nil {
		g.Log().Warningf(ctx, "获取缓存失败[%s]: %v", key, err)
		return err
	}

	if result == nil {
		return ErrCacheMiss
	}

	// 将结果转换为指定类型
	return gconv.Struct(result, value)
}

// SetCache 设置Redis缓存
func SetCache(ctx context.Context, key string, value interface{}, expiry time.Duration) error {
	redis := GetRedisClient()
	if redis == nil {
		return errors.New("Redis未配置")
	}

	if expiry > 0 {
		_, err := redis.Do(ctx, "SETEX", key, int64(expiry.Seconds()), value)
		if err != nil {
			g.Log().Warningf(ctx, "设置缓存失败[%s]: %v", key, err)
		}
		return err
	}

	_, err := redis.Do(ctx, "SET", key, value)
	if err != nil {
		g.Log().Warningf(ctx, "设置缓存失败[%s]: %v", key, err)
	}
	return err
}

// DelCache 删除Redis缓存
func DelCache(ctx context.Context, key string) error {
	_, err := g.Redis().Do(ctx, "DEL", key)
	return err
}

// Exists 检查键是否存在
func Exists(ctx context.Context, key string) (bool, error) {
	redis := GetRedisClient()
	if redis == nil {
		return false, errors.New("Redis未配置")
	}

	result, err := redis.Do(ctx, "EXISTS", key)
	if err != nil {
		g.Log().Warningf(ctx, "检查键是否存在失败[%s]: %v", key, err)
		return false, err
	}

	return gconv.Int(result) == 1, nil
}

// GetDistributedLock 获取分布式锁
// 使用SETNX实现的简单分布式锁
func GetDistributedLock(ctx context.Context, key string, expireSeconds int) (bool, error) {
	// 设置锁，使用NX确保只有一个客户端可以获取锁
	// 设置过期时间防止死锁
	result, err := g.Redis().Do(ctx, "SET", key, "1", "NX", "EX", expireSeconds)
	if err != nil {
		return false, err
	}

	// 检查是否获取到锁
	return result != nil && gconv.String(result) == "OK", nil
}

// ReleaseDistributedLock 释放分布式锁
func ReleaseDistributedLock(ctx context.Context, key string) error {
	_, err := g.Redis().Do(ctx, "DEL", key)
	return err
}

// DeductStockScript Lua脚本：扣减库存
// 原子性地检查和扣减库存，同时检查用户是否已购买
// 参数:
//
//	KEYS[1] - 库存键
//	KEYS[2] - 用户购买记录键
//	KEYS[3] - 成功计数键
//	ARGV[1] - 扣减数量
//	ARGV[2] - 订单号
//	ARGV[3] - 超时时间（秒）
//
// 返回:
//
//	1 - 扣减成功，返回订单号
//	0,msg - 扣减失败，返回失败消息
const DeductStockScript = `
local stock = redis.call('get', KEYS[1])
if not stock or tonumber(stock) < tonumber(ARGV[1]) then
    return {0, "库存不足"}
end

-- 检查用户是否已购买
local bought = redis.call('exists', KEYS[2])
if tonumber(bought) > 0 then
    return {0, "您已参与过该商品的秒杀"}
end

-- 扣减库存
redis.call('decrby', KEYS[1], ARGV[1])

-- 记录用户购买
redis.call('setex', KEYS[2], ARGV[3], ARGV[2])

-- 增加成功计数
redis.call('incr', KEYS[3])

return {1, ARGV[2]}
`

// DeductStockWithLua 使用Lua脚本扣减库存
func DeductStockWithLua(ctx context.Context, keys *RedisKeys, count int, orderNo string) (bool, string, error) {
	// 执行Lua脚本
	result, err := g.Redis().Do(ctx, "EVAL", DeductStockScript, 3,
		keys.StockKey, keys.UserBoughtKey, keys.SuccessCountKey,
		count, orderNo, 86400) // 24小时过期

	if err != nil {
		return false, "", err
	}

	// 解析结果
	resultArray := gconv.SliceAny(result)
	if len(resultArray) != 2 {
		return false, "", fmt.Errorf("无效的脚本返回结果")
	}

	statusCode := gconv.Int(resultArray[0])
	resultMessage := gconv.String(resultArray[1])

	return statusCode == 1, resultMessage, nil
}

// RollbackStockScript Lua脚本：回滚库存
// 原子性地回滚库存并删除用户购买记录
const RollbackStockScript = `
local stock = redis.call('get', KEYS[1])
if stock then
    redis.call('incrby', KEYS[1], ARGV[1])
end

redis.call('del', KEYS[2])

-- 减少成功计数
local successCount = redis.call('get', KEYS[3])
if successCount and tonumber(successCount) > 0 then
    redis.call('decr', KEYS[3])
end

return 1
`

// RollbackStock 回滚库存
func RollbackStock(ctx context.Context, keys *RedisKeys, count int) error {
	_, err := g.Redis().Do(ctx, "EVAL", RollbackStockScript, 3,
		keys.StockKey, keys.UserBoughtKey, keys.SuccessCountKey, count)
	return err
}

// InitStockScript Lua脚本：初始化库存
// 原子性地初始化或更新库存值
const InitStockScript = `
redis.call('set', KEYS[1], ARGV[1])
redis.call('del', KEYS[2]) -- 删除成功计数
return 1
`

// InitStock 初始化库存
func InitStock(ctx context.Context, goodsId, goodsOptionsId int64, stock int) error {
	keys := NewRedisKeys(goodsId, goodsOptionsId, 0)
	_, err := g.Redis().Do(ctx, "EVAL", InitStockScript, 2,
		keys.StockKey, keys.SuccessCountKey, stock)
	return err
}

// CacheGoodsInfo 缓存商品信息
func CacheGoodsInfo(ctx context.Context, goodsId, goodsOptionsId int64, goodsInfo []byte, expireSeconds int) error {
	keys := NewRedisKeys(goodsId, goodsOptionsId, 0)
	_, err := g.Redis().Do(ctx, "SETEX", keys.GoodsInfoKey, expireSeconds, goodsInfo)
	return err
}

// GetGoodsInfo 获取商品信息
func GetGoodsInfo(ctx context.Context, goodsId, goodsOptionsId int64) ([]byte, error) {
	keys := NewRedisKeys(goodsId, goodsOptionsId, 0)
	result, err := g.Redis().Do(ctx, "GET", keys.GoodsInfoKey)
	if err != nil {
		return nil, err
	}

	if result == nil {
		return nil, nil
	}

	return []byte(gconv.String(result)), nil
}

// GetStockAndSuccess 获取库存和成功数
func GetStockAndSuccess(ctx context.Context, goodsId, goodsOptionsId int64) (int, int, error) {
	keys := NewRedisKeys(goodsId, goodsOptionsId, 0)

	// 使用MGET一次获取多个值
	result, err := g.Redis().Do(ctx, "MGET", keys.StockKey, keys.SuccessCountKey)
	if err != nil {
		return 0, 0, err
	}

	resultArray := gconv.SliceAny(result)
	if len(resultArray) != 2 {
		return 0, 0, fmt.Errorf("无效的返回结果")
	}

	var stock, success int
	if resultArray[0] != nil {
		stock = gconv.Int(resultArray[0])
	}

	if resultArray[1] != nil {
		success = gconv.Int(resultArray[1])
	}

	return stock, success, nil
}

// SetProcessingFlag 设置处理中标记
func SetProcessingFlag(ctx context.Context, keys *RedisKeys, expireSeconds int) error {
	_, err := g.Redis().Do(ctx, "SETEX", keys.ProcessingFlagKey, expireSeconds, 1)
	return err
}

// IsProcessing 检查是否处理中
func IsProcessing(ctx context.Context, keys *RedisKeys) (bool, error) {
	result, err := g.Redis().Do(ctx, "EXISTS", keys.ProcessingFlagKey)
	if err != nil {
		return false, err
	}

	return gconv.Int(result) > 0, nil
}

// SetSeckillResult 设置秒杀结果
func SetSeckillResult(ctx context.Context, keys *RedisKeys, success bool, orderNo string, expireSeconds int) error {
	value := "fail"
	if success {
		value = orderNo
	}

	_, err := g.Redis().Do(ctx, "SETEX", keys.ResultKey, expireSeconds, value)
	return err
}

// GetSeckillResult 获取秒杀结果
func GetSeckillResult(ctx context.Context, keys *RedisKeys) (bool, string, error) {
	result, err := g.Redis().Do(ctx, "GET", keys.ResultKey)
	if err != nil {
		return false, "", err
	}

	if result == nil {
		return false, "", nil
	}

	resultStr := gconv.String(result)
	if resultStr == "fail" {
		return true, "", nil
	}

	return true, resultStr, nil
}

// WarmUpCache 预热秒杀系统缓存
func WarmUpCache(ctx context.Context) {
	g.Log().Info(ctx, "正在预热秒杀系统缓存...")

	// 创建后台上下文
	if ctx == nil {
		ctx = gctx.New()
	}

	// 预热过程
	go func() {
		// 加载活跃商品库存到Redis
		if err := warmUpStockCache(ctx); err != nil {
			g.Log().Error(ctx, "预热库存缓存失败:", err)
		}

		// 预热结束
		g.Log().Info(ctx, "秒杀系统缓存预热完成")
	}()
}

// 预热库存缓存
func warmUpStockCache(ctx context.Context) error {
	// 这里应该查询数据库中的秒杀商品，并将库存加载到Redis
	// 简化实现，实际应用中应该根据具体业务逻辑查询

	// 获取当前时间
	now := time.Now()

	// 查询进行中的秒杀活动
	var goods []map[string]interface{}
	err := g.DB().Model("seckill_goods").
		Where("start_time <= ? AND end_time > ? AND status = ?", now, now, 1).
		Scan(&goods)

	if err != nil {
		return err
	}

	// 遍历商品，预热缓存
	for _, item := range goods {
		goodsId := gconv.Int64(item["goods_id"])
		goodsOptionsId := gconv.Int64(item["goods_options_id"])
		stock := gconv.Int(item["seckill_stock"])

		// 初始化库存
		if err := InitStock(ctx, goodsId, goodsOptionsId, stock); err != nil {
			g.Log().Error(ctx, fmt.Sprintf("初始化商品[%d]库存失败:", goodsId), err)
			continue
		}

		// 缓存商品信息
		data, err := gjson.Encode(item)
		if err != nil {
			g.Log().Error(ctx, fmt.Sprintf("序列化商品[%d]信息失败:", goodsId), err)
			continue
		}

		if err := CacheGoodsInfo(ctx, goodsId, goodsOptionsId, data, 86400); err != nil {
			g.Log().Error(ctx, fmt.Sprintf("缓存商品[%d]信息失败:", goodsId), err)
			continue
		}

		g.Log().Info(ctx, fmt.Sprintf("商品[%d]库存预热成功, 库存: %d", goodsId, stock))
	}

	return nil
}

// CheckOrderStatusScript Lua脚本：检查订单状态
// 检查订单是否已经处理，避免重复处理
const CheckOrderStatusScript = `
local processing = redis.call('exists', KEYS[1])
local completed = redis.call('exists', KEYS[2])

if tonumber(processing) > 0 then
    return {0, "订单正在处理中"}
end

if tonumber(completed) > 0 then
    return {0, "订单已处理"}
end

-- 标记为处理中
redis.call('setex', KEYS[1], ARGV[1], 1)
return {1, "可以处理"}
`

// CheckAndMarkOrderProcessing 检查并标记订单处理中
func CheckAndMarkOrderProcessing(ctx context.Context, userId int64, expireSeconds int) (bool, error) {
	keys := NewRedisKeys(0, 0, userId)

	result, err := g.Redis().Do(ctx, "EVAL", CheckOrderStatusScript, 2,
		keys.OrderProcessingKey, keys.OrderCompletedKey, expireSeconds)

	if err != nil {
		return false, err
	}

	resultArray := gconv.SliceAny(result)
	if len(resultArray) != 2 {
		return false, fmt.Errorf("无效的脚本返回结果")
	}

	statusCode := gconv.Int(resultArray[0])
	return statusCode == 1, nil
}

// MarkOrderCompleted 标记订单已完成
func MarkOrderCompleted(ctx context.Context, userId int64, orderNo string, expireSeconds int) error {
	keys := NewRedisKeys(0, 0, userId)

	// 删除处理中标记，设置完成标记
	_, err := g.Redis().Do(ctx, "DEL", keys.OrderProcessingKey)
	if err != nil {
		return err
	}

	_, err = g.Redis().Do(ctx, "SETEX", keys.OrderCompletedKey, expireSeconds, orderNo)
	return err
}

// BloomFilter 布隆过滤器
type BloomFilter struct {
	key    string        // Redis键
	size   uint64        // 位图大小
	hashes int           // 哈希函数数量
	redis  *gredis.Redis // Redis客户端
}

// NewBloomFilter 创建布隆过滤器
func NewBloomFilter(key string, size uint64, hashes int) *BloomFilter {
	return &BloomFilter{
		key:    key,
		size:   size,
		hashes: hashes,
		redis:  GetRedisClient(),
	}
}

// Add 添加元素到布隆过滤器
func (bf *BloomFilter) Add(ctx context.Context, value string) error {
	locations := bf.getLocations(value)
	for _, loc := range locations {
		_, err := bf.redis.Do(ctx, "SETBIT", bf.key, loc, 1)
		if err != nil {
			return err
		}
	}
	return nil
}

// Exists 检查元素是否可能存在
func (bf *BloomFilter) Exists(ctx context.Context, value string) (bool, error) {
	locations := bf.getLocations(value)
	for _, loc := range locations {
		exists, err := bf.redis.Do(ctx, "GETBIT", bf.key, loc)
		if err != nil {
			return false, err
		}
		if exists.Int() == 0 {
			return false, nil
		}
	}
	return true, nil
}

// getLocations 计算元素的哈希位置
func (bf *BloomFilter) getLocations(value string) []uint64 {
	locations := make([]uint64, bf.hashes)
	for i := 0; i < bf.hashes; i++ {
		// 使用不同的种子生成不同的哈希值
		hash := fnv.New64a()
		hash.Write([]byte(value))
		hash.Write([]byte{byte(i)})
		locations[i] = hash.Sum64() % bf.size
	}
	return locations
}
