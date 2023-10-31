package consts

const (
	GTokenUserPrefix = "User:" //gToken登录 前台用户 前缀区分
	CtxUserName      = "CtxUserName"
	CtxUserId        = "CtxUserId"
)

// 事件名称
const (
	MsgEventPing = "ping"
	MsgEventPong = "pong"
	MsgEventAck  = "ack"
)

// topic
const (
	ImTopicFoo = "im:message:foo"
	ImTopicBar = "im:message:bar"
)

const (
	SubEventImMessageFoo = "sub.im.message.foo"
	SubEventImMessageBar = "sub.im.message.bar"
)

const (
	PushEventImMessageFoo = "im.message.foo"
	PushEventImMessageBar = "im.message.bar"
)

const (
	BarTypeA = "barTypeA"
	BarTypeB = "barTypeB"
)
