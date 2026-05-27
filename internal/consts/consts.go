package consts

const (
	ProjectName              = "Go开源电商实战项目"
	ProjectUsage             = "演示学习使用，作者：王中阳Go，微信：wangzhongyang1993，公众号：程序员升职加薪之旅"
	ProjectBrief             = "start http server"
	Version                  = "v0.2.0"             // 当前服务版本(用于模板展示)
	CaptchaDefaultName       = "CaptchaDefaultName" // 验证码默认存储空间名称
	ContextKey               = "ContextKey"         // 上下文变量存储键名，前后端系统共享
	FileMaxUploadCountMinute = 10                   // 同一用户1分钟之内最大上传数量
	GTokenAdminPrefix        = "Admin:"             //gtoken登录 管理后台 前缀区分
	GTokenFrontendPrefix     = "User:"              //gtoken登录 前台用户 前缀区分
	//for admin
	CtxAdminId      = "CtxAdminId"
	CtxAdminName    = "CtxAdminName"
	CtxAdminIsAdmin = "CtxAdminIsAdmin"
	CtxAdminRoleIds = "CtxAdminRoleIds"
	//for user
	CtxUserId     = "CtxUserId"
	CtxUserName   = "CtxUserName"
	CtxUserAvatar = "CtxUserAvatar"
	CtxUserSex    = "CtxUserSex"
	CtxUserSign   = "CtxUserSign"
	CtxUserStatus = "CtxUserStatus"
	//for 登录相关
	TokenType          = "Bearer"
	CacheModeRedis     = 2
	BackendServerName  = "开源电商系统"
	MultiLogin         = true
	FrontendMultiLogin = true
	GTokenExpireIn     = 10 * 24 * 60 * 60
	//统一管理错误提示
	CodeMissingParameterMsg = "请检查是否缺少参数"
	ErrLoginFaulMsg         = "登录失败，账号或密码错误"
	ErrSecretAnswerMsg      = "密保问题不正确"
	ResourcePermissionFail  = "没有权限操作"
	//收藏相关
	CollectionTypeGoods   = 1
	CollectionTypeArticle = 2
	//点赞相关
	PraiseTypeGoods   = 1
	PraiseTypeArticle = 2
	//评论相关
	CommentTypeGoods   = 1
	CommentTypeArticle = 2
	//收货地址相关
	ProvincePid = 1
	//订单评论默认时间 7天 超过7天后默认好评 7 * 24 * 60 * 60
	UserOrderDefaultCommentsTime = 7 * 24 * 60 * 60
	UserOrderStatus              = 5
	UserOrderDefaultComments     = "系统默认好评"
	//文章相关
	ArticleIsAdmin = 1 //管理员发布
	ArticleIsUser  = 2 //用户发布
	//售后相关
	RefundStatusWait   = 1
	RefundStatusAgree  = 2
	RefundStatusRejuct = 3

	// 秒杀系统相关常量
	SeckillDefaultStock = 1000  // 默认秒杀商品库存
	SeckillMaxStock     = 10000 // 最大秒杀商品库存
	SeckillSlotCount    = 16384 // Redis集群槽位数量

	// 秒杀令牌桶限流器配置
	SeckillTokenBucketSize = 5000 // 令牌桶容量，原值过小，调大以提高并发处理能力
	SeckillTokenRate       = 1000 // 令牌产生速率（每秒），提高速率以支持更高QPS

	// 秒杀漏桶限流器配置
	SeckillLeakyBucketSize = 10000 // 漏桶容量，原值过小，调大以容纳更多请求
	SeckillLeakyBucketRate = 2000  // 漏桶处理速率（每秒），提高以支持更高吞吐量

	// Kafka主题
	KafkaTopicSeckill         = "seckill_order"          // 秒杀订单主题
	KafkaTopicSeckillComplete = "seckill_order_complete" // 秒杀订单完成主题
	SeckillKafkaTopic         = "seckill_order"          // 秒杀订单主题（兼容现有代码）

	// Redis键前缀
	SeckillGoodsPrefix       = "seckill:goods:"       // 秒杀商品信息前缀
	SeckillStockPrefix       = "seckill:stock:"       // 秒杀库存前缀
	SeckillUserBoughtPrefix  = "seckill:user:bought:" // 用户购买记录前缀
	SeckillSuccessPrefix     = "seckill:success:"     // 秒杀成功计数前缀
	SeckillLockPrefix        = "seckill:lock:"        // 秒杀锁前缀
	SeckillOrderSentPrefix   = "seckill:order:sent:"  // 秒杀订单已发送标记前缀
	SeckillTokenBucketPrefix = "seckill:token:"       // 令牌桶前缀
	SeckillLeakyBucketPrefix = "seckill:leaky:"       // 漏桶前缀
	SeckillResultPrefix      = "seckill:result:"      // 秒杀结果前缀
	SeckillQueuePrefix       = "seckill:queue:"       // 秒杀队列前缀

	// 秒杀脚本相关
	SeckillScriptMaxRetries = 3  // Lua脚本最大重试次数
	SeckillScriptRetryDelay = 50 // Lua脚本重试延迟(毫秒)

	// 秒杀统计和队列
	SeckillMetricsPrefix   = "seckill:metrics:" // 秒杀统计指标前缀
	SeckillQueueExpireTime = 3600               // 秒杀队列过期时间(秒)
)
