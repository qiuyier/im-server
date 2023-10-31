package model

type BarMessage struct {
	ReceiverId int    `json:"receiver_id"`
	BarType    int    `json:"bar_type"`
	RecordId   string `json:"record_id"`
}

type BarReq struct {
	Type       string `json:"type"`
	ReceiverId uint   `json:"receiver_id" v:"required"`
}

type BarRes struct {
	RecordId string `json:"record_id"`
}

type BarTypeAInput struct {
	BarReq
	Bar string `json:"bar"`
}

type BarTypeBInput struct {
	BarReq
	Foo string `json:"foo"`
}
