package cache

import "context"

type IBind interface {
	Bind(ctx context.Context, channel string, cid int64, uid int) error
	UnBind(ctx context.Context, channel string, cid int64) error
}
