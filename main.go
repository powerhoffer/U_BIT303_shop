package main

import (
	"bit303_shop/internal/cmd"
	_ "bit303_shop/internal/logic"

	"github.com/gogf/gf/v2/os/gctx"

	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"
)

func main() {
	cmd.Main.Run(gctx.New())
}
