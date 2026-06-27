package backend

import (
	"context"

	"bit303_shop/internal/consts"

	"github.com/gogf/gf/v2/frame/g"
)

func currentEmployeeId(ctx context.Context) uint {
	r := g.RequestFromCtx(ctx)
	if r == nil {
		return 0
	}
	return r.GetCtxVar(consts.CtxEmployeeId).Uint()
}

func currentAdminId(ctx context.Context) uint {
	r := g.RequestFromCtx(ctx)
	if r == nil {
		return 0
	}
	return r.GetCtxVar(consts.CtxAdminId).Uint()
}
