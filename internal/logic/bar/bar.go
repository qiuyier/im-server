package bar

import (
	"context"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/guid"
	"im/internal/consts"
	"im/internal/model"
	"im/internal/service"
)

type sBar struct {
}

func init() {
	service.RegisterBar(New())
}

func New() *sBar {
	return &sBar{}
}

var adapterMap map[string]func(ctx context.Context) (*model.BarRes, error)

func (s *sBar) Bar(ctx context.Context, in model.BarReq) (out *model.BarRes, err error) {
	return s.adapter(ctx, in.Type)
}

func (s *sBar) adapter(ctx context.Context, typeValue string) (*model.BarRes, error) {
	if adapterMap == nil {
		adapterMap = make(map[string]func(ctx context.Context) (*model.BarRes, error))
		adapterMap[consts.BarTypeA] = s.onBarTypeA
		adapterMap[consts.BarTypeB] = s.onBarTypeB
	}

	if call, ok := adapterMap[typeValue]; ok {
		return call(ctx)
	}

	return nil, nil
}

func (s *sBar) onBarTypeA(ctx context.Context) (*model.BarRes, error) {
	// 获取参数
	param := &model.BarTypeAInput{}
	err := gjson.Unmarshal(g.RequestFromCtx(ctx).GetBody(), param)
	if err != nil {
		g.Log().Error(ctx, err)
	}

	// TODO business logic
	id := guid.S(gconv.Bytes(param.Bar))
	out := &model.BarRes{RecordId: id}

	s.afterHandle(ctx, &model.BarMessage{
		ReceiverId: gconv.Int(param.ReceiverId),
		BarType:    gconv.Int(param.Type),
		RecordId:   id,
	}, consts.SubEventImMessageBar)

	return out, nil
}

func (s *sBar) onBarTypeB(ctx context.Context) (*model.BarRes, error) {
	// 获取参数
	param := &model.BarTypeBInput{}
	err := gjson.Unmarshal(g.RequestFromCtx(ctx).GetBody(), param)
	if err != nil {
		g.Log().Error(ctx, err)
	}

	// TODO business logic
	id := guid.S(gconv.Bytes(param.Foo))
	out := &model.BarRes{RecordId: id}

	s.afterHandle(ctx, &model.BarMessage{
		ReceiverId: gconv.Int(param.ReceiverId),
		BarType:    gconv.Int(param.Type),
		RecordId:   id,
	}, consts.SubEventImMessageFoo)

	return out, nil
}

func (s *sBar) afterHandle(ctx context.Context, data *model.BarMessage, event string) {
	content := gjson.MustEncodeString(&model.SubscribeContent{
		Event: event,
		Data:  gjson.MustEncodeString(data),
	})

	if _, err := g.Redis().Publish(ctx, consts.ImTopicFoo, content); err != nil {
		g.Log().Error(ctx, "[ALL]消息推送失败 err:", err)
	}
}
