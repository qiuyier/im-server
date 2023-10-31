// Package cache 绑定客户端和用户的关系，借由redis存储
package cache

import (
	"context"
	_ "embed"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
)

var (
	//go:embed scripts/bind.lua
	bindLua string
	//go:embed scripts/unbind.lua
	unbindLua string
)

type ClientCache struct {
}

func NewClientCache() *ClientCache {
	return &ClientCache{}
}

func (c *ClientCache) Bind(ctx context.Context, channel string, cid int64, uid int) error {
	_, err := g.Redis().Eval(ctx, bindLua, 2, []string{c.clientKey(channel), c.userKey(channel, gconv.String(uid))}, []any{cid, uid})
	if err != nil {
		return err
	}
	return nil
}

func (c *ClientCache) UnBind(ctx context.Context, channel string, cid int64) error {
	key := c.clientKey(channel)
	uid, err := g.Redis().HGet(ctx, key, gconv.String(cid))
	if err != nil {
		return err
	}
	_, err = g.Redis().Eval(ctx, unbindLua, 2, []string{key, c.userKey(channel, uid.String())}, []any{cid, uid})
	if err != nil {
		return err
	}
	return nil
}

func (c *ClientCache) clientKey(channel string) string {
	return fmt.Sprintf("ws:%s:client", channel)
}

func (c *ClientCache) userKey(channel, uid string) string {
	return fmt.Sprintf("ws:%s:user:%s", channel, uid)
}

func (c *ClientCache) GetClientIdsFromUid(ctx context.Context, channel, uid string) []int64 {
	clientIds := make([]int64, 0)

	items, err := g.Redis().SMembers(ctx, c.userKey(channel, uid))
	if err != nil {
		return clientIds
	}

	for _, cid := range items.Int64s() {
		clientIds = append(clientIds, cid)
	}

	return clientIds
}
