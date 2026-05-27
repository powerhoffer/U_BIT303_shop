package test

import (
	"context"
	"os"
	"testing"
	"time"

	_ "github.com/gogf/gf/contrib/nosql/redis/v2"
	"github.com/gogf/gf/v2/database/gredis"
)

func skipRedisIntegration(t *testing.T) {
	t.Helper()
	if os.Getenv("BIT303_RUN_REDIS_TESTS") != "1" {
		t.Skip("set BIT303_RUN_REDIS_TESTS=1 to run Redis integration tests")
	}
}

// TestRedisConnection 测试无密码连接Redis
func TestRedisConnection(t *testing.T) {
	skipRedisIntegration(t)

	config := gredis.Config{
		Address:     "127.0.0.1:6379",
		Db:          1,
		Pass:        "", // 无密码
		IdleTimeout: 60 * time.Second,
	}

	redis, err := gredis.New(&config)
	if err != nil {
		t.Fatalf("创建无密码Redis客户端失败: %v", err)
	}

	ctx := context.Background()
	result, err := redis.Do(ctx, "PING")
	if err != nil {
		t.Fatalf("无密码Redis PING测试失败: %v", err)
	}

	t.Logf("无密码Redis PING成功: %v", result)
}

// TestRedisWithPasswordConnection 测试带密码连接Redis
func TestRedisWithPasswordConnection(t *testing.T) {
	skipRedisIntegration(t)

	config := gredis.Config{
		Address:     "127.0.0.1:6379",
		Db:          1,
		Pass:        "123456", // 使用配置文件中的密码
		IdleTimeout: 60 * time.Second,
	}

	redis, err := gredis.New(&config)
	if err != nil {
		t.Fatalf("创建带密码Redis客户端失败: %v", err)
	}

	ctx := context.Background()
	result, err := redis.Do(ctx, "PING")
	if err != nil {
		t.Fatalf("带密码Redis PING测试失败: %v", err)
	}

	t.Logf("带密码Redis PING成功: %v", result)
}

// TestRedisNoPasswordACL 测试使用默认用户无密码连接
func TestRedisNoPasswordACL(t *testing.T) {
	skipRedisIntegration(t)

	config := gredis.Config{
		Address:     "127.0.0.1:6379",
		Db:          1,
		User:        "default", // 使用默认用户
		Pass:        "",        // 无密码
		IdleTimeout: 60 * time.Second,
	}

	redis, err := gredis.New(&config)
	if err != nil {
		t.Fatalf("创建默认用户Redis客户端失败: %v", err)
	}

	ctx := context.Background()
	result, err := redis.Do(ctx, "PING")
	if err != nil {
		t.Fatalf("默认用户Redis PING测试失败: %v", err)
	}

	t.Logf("默认用户Redis PING成功: %v", result)
}
