package main

import (
	_ "github.com/gogf/gf/contrib/nosql/redis/v2"

	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"

	_ "im/internal/packed"

	"github.com/gogf/gf/v2/os/gctx"

	"im/internal/cmd"
)

func main() {
	cmd.Main.Run(gctx.GetInitCtx())
}
