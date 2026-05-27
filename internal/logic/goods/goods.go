package goods

import (
	"context"
	"bit303_shop/internal/consts"
	"bit303_shop/internal/dao"
	"bit303_shop/internal/model"
	"bit303_shop/internal/service"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/golang-jwt/jwt/v5"
)

// getMapKeys 获取map的所有键
func getMapKeys(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

type sGoods struct{}

func init() {
	service.RegisterGoods(New())
}

func New() *sGoods {
	return &sGoods{}
}

func (s *sGoods) Create(ctx context.Context, in model.GoodsCreateInput) (out model.GoodsCreateOutput, err error) {
	lastInsertID, err := dao.GoodsInfo.Ctx(ctx).Data(in).InsertAndGetId()
	if err != nil {
		return out, err
	}
	return model.GoodsCreateOutput{Id: uint(lastInsertID)}, err
}

// Delete 删除
func (s *sGoods) Delete(ctx context.Context, id uint) (err error) {
	_, err = dao.GoodsInfo.Ctx(ctx).Where(g.Map{
		dao.GoodsInfo.Columns().Id: id,
	}).Delete()
	if err != nil {
		return err
	}
	return
}

// Update 修改
func (s *sGoods) Update(ctx context.Context, in model.GoodsUpdateInput) error {
	_, err := dao.GoodsInfo.
		Ctx(ctx).
		Data(in).
		FieldsEx(dao.GoodsInfo.Columns().Id).
		Where(dao.GoodsInfo.Columns().Id, in.Id).
		Update()
	return err
}

// GetList 查询分类列表
func (s *sGoods) GetList(ctx context.Context, in model.GoodsGetListInput) (out *model.GoodsGetListOutput, err error) {
	//1.获得*gdb.Model对象，方面后续调用
	m := dao.GoodsInfo.Ctx(ctx)
	//2. 实例化响应结构体
	out = &model.GoodsGetListOutput{
		Page: in.Page,
		Size: in.Size,
	}
	//3. 分页查询
	listModel := m.Page(in.Page, in.Size)
	//4. 再查询count，判断有无数据
	out.Total, err = m.Count()
	if err != nil || out.Total == 0 {
		return out, err
	}
	//5. 延迟初始化list切片 确定有数据，再按期望大小初始化切片容量
	out.List = make([]model.GoodsGetListOutputItem, 0, in.Size)
	//6. 把查询到的结果赋值到响应结构体中
	if err := listModel.Scan(&out.List); err != nil {
		return out, err
	}
	return
}

// GetLevelList 根据分类级别获取商品列表
func (s *sGoods) GetLevelList(ctx context.Context, in model.GoodsGetLevelListInput) (out *model.GoodsGetLevelListOutput, err error) {
	//1.获得*gdb.Model对象，方便后续调用
	m := dao.GoodsInfo.Ctx(ctx)
	
	//2. 根据level_id查询商品，这里需要匹配任一级分类
	m = m.Where("level1_category_id = ? OR level2_category_id = ? OR level3_category_id = ?", in.LevelId, in.LevelId, in.LevelId)
	
	//3. 实例化响应结构体
	out = &model.GoodsGetLevelListOutput{
		Page: in.Page,
		Size: in.Size,
	}
	
	//4. 分页查询
	listModel := m.Page(in.Page, in.Size)
	
	//5. 再查询count，判断有无数据
	out.Total, err = m.Count()
	if err != nil || out.Total == 0 {
		return out, err
	}
	
	//6. 延迟初始化list切片 确定有数据，再按期望大小初始化切片容量
	out.List = make([]model.GoodsGetLevelListOutputItem, 0, in.Size)
	
	//7. 把查询到的结果赋值到响应结构体中
	if err := listModel.Scan(&out.List); err != nil {
		return out, err
	}
	
	return
}

// 商品详情
func (*sGoods) Detail(ctx context.Context, in model.GoodsDetailInput) (out model.GoodsDetailOutput, err error) {
	err = dao.GoodsInfo.Ctx(ctx).WithAll().WherePri(in.Id).Scan(&out)
	if err != nil {
		return model.GoodsDetailOutput{}, err
	}

	// 检查用户是否已收藏该商品
	userIdValue := ctx.Value(consts.CtxUserId)
	
	// 如果中间件没有获取到用户ID，尝试手动解析JWT token（如果用户携带了token）
	if userIdValue == nil {
		if r := ghttp.RequestFromCtx(ctx); r != nil {
			// 尝试从多个位置获取token
			token := r.Header.Get("Authorization")
			if token == "" {
				token = r.Get("token").String()
			}
			if token == "" {
				token = r.Cookie.Get("jwt_frontend").String()
			}
			
			// 处理Bearer前缀
			if len(token) > 7 && token[:7] == "Bearer " {
				token = token[7:]
			}
			
			if token != "" {
				// 简单的min函数实现
				tokenLen := len(token)
				if tokenLen > 20 {
					tokenLen = 20
				}
				g.Log().Infof(ctx, "发现token: %s...", token[:tokenLen])
				
				// 直接使用标准JWT库解析token，不依赖中间件
				userAuthService := service.GetUserAuthService()
				if userAuthService != nil {
					// 使用标准JWT库解析token
					parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
					// 验证签名方法
					if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
						return nil, jwt.ErrSignatureInvalid
					}
						// 使用相同的密钥
						return userAuthService.Key, nil
					})
					
					if err != nil {
						g.Log().Errorf(ctx, "解析JWT token失败: %v", err)
					} else if parsedToken.Valid {
						claims := parsedToken.Claims.(jwt.MapClaims)
						g.Log().Infof(ctx, "JWT解析成功，claims: %v", claims)
						if userId, exists := claims[userAuthService.IdentityKey]; exists {
							userIdValue = userId
							g.Log().Infof(ctx, "手动解析JWT成功，用户ID: %v", userIdValue)
						} else {
							g.Log().Infof(ctx, "JWT中未找到IdentityKey: %s", userAuthService.IdentityKey)
							g.Log().Infof(ctx, "JWT中的所有键: %v", getMapKeys(claims))
						}
					} else {
						g.Log().Errorf(ctx, "JWT token无效")
					}
				}
			} else {
				g.Log().Infof(ctx, "未发现token")
			}
		}
	}
	
	// 如果获取到了用户ID，检查收藏状态
	if userIdValue != nil && userIdValue != 0 {
		isCollect, err := service.Collection().CheckIsCollect(ctx, model.CheckIsCollectInput{
			UserId:   gconv.Uint(userIdValue),
			ObjectId: gconv.Uint(out.Id),
			Type:     consts.CollectionTypeGoods,
		})
		if err != nil {
			g.Log().Errorf(ctx, "检查商品收藏状态失败: %v", err)
		} else {
			out.IsCollect = isCollect
			g.Log().Infof(ctx, "商品收藏状态: %v", isCollect)
		}
	} else {
		g.Log().Infof(ctx, "用户未登录，设置收藏状态为false")
		out.IsCollect = false
	}

	return
}
