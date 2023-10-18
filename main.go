package main

import (
	_ "im/internal/packed"

	"github.com/gogf/gf/v2/os/gctx"

	"im/internal/cmd"
)

func main() {
	cmd.Main.Run(gctx.GetInitCtx())
}
