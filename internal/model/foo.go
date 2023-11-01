package model

type FooMessage struct {
	ReceiverId int    `json:"receiver_id"`
	AcceptData string `json:"accept_data"`
}
