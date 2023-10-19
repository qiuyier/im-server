package config

import (
	"fmt"
	"github.com/gogf/gf/v2/crypto/gmd5"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/os/gcfg"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/util/grand"
	"time"
)

var Cfg *Config

type Config struct {
	sid string `json:"sid"`
	Jwt Jwt    `yaml:"jwt" json:"jwt"`
}

type Jwt struct {
	Secret      string `yaml:"secret" json:"secret"`
	ExpiresTime int    `yaml:"expires_time" json:"expires_time"`
	BufferTime  int    `yaml:"buffer_time" json:"buffer_time"`
}

func init() {
	file, _ := gcfg.NewAdapterFile()
	path, _ := file.GetFilePath()

	if err := gjson.Unmarshal(gjson.MustEncode(gcfg.Instance().MustData(gctx.New())), &Cfg); err != nil {
		panic(fmt.Sprintf("Parsing configuration files %s err: %v", path, err))
	}

	// 生成服务运行ID
	Cfg.sid = gmd5.MustEncryptString(fmt.Sprintf("%d%s", time.Now().UnixNano(), grand.S(6)))
}
