// Package websocket 统一标准websocket消息返回结构体
package websocket

type Message struct {
	Event   string `json:"event"`   // 事件名称
	Content any    `json:"content"` // 消息内容
}

type SenderContent struct {
	IsAck     bool
	broadcast bool     // 是否广播消息
	exclude   []int64  // 排除的用户(预留)
	receivers []int64  // 推送的用户
	message   *Message // 消息体
}

func NewMessage(event string, content any) *Message {
	return &Message{
		Event:   event,
		Content: content,
	}
}

func NewSenderContent() *SenderContent {
	return &SenderContent{
		exclude:   make([]int64, 0),
		receivers: make([]int64, 0),
	}
}

func (s *SenderContent) SetAck(value bool) *SenderContent {
	s.IsAck = value
	return s
}

func (s *SenderContent) SetBroadcast(value bool) *SenderContent {
	s.broadcast = value
	return s
}

func (s *SenderContent) SetMessage(event string, content any) *SenderContent {
	s.message = NewMessage(event, content)
	return s
}

func (s *SenderContent) SetReceiver(cid ...int64) *SenderContent {
	s.receivers = append(s.receivers, cid...)
	return s
}

func (s *SenderContent) SetExclude(cid ...int64) *SenderContent {
	s.exclude = append(s.exclude, cid...)
	return s
}

func (s *SenderContent) IsBroadcast() bool {
	return s.broadcast
}
