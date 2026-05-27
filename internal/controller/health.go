package controller

import (
	"context"

	"bit303_shop/api/base"
)

var Health = cHealth{}

type cHealth struct{}

func (c cHealth) Check(ctx context.Context, req *base.HealthReq) (res *base.HealthRes, err error) {
	return &base.HealthRes{Status: "ok"}, nil
}
