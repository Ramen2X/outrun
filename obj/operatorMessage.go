package obj

import (
	"strconv"
	"time"

	"github.com/fluofoxxo/outrun/enums"
)

type OperatorMessage struct {
	ID         string      `json:"messageId"`
	Content    string      `json:"contents"`
	Item       MessageItem `json:"item"`
	ExpireTime int64       `json:"expireTime"`
}

func DefaultOperatorMessage() OperatorMessage {
	id := "2346789"
	content := "Test Gift"
	item := NewMessageItem(
		enums.ItemIDRing,
		15000,
		0,
		0,
	)
	expireTime := time.Now().Unix() + 12600 // three and a half hours from now
	return OperatorMessage{
		id,
		content,
		item,
		expireTime,
	}
}

func NewOperatorMessage(id int64, content string, item MessageItem, expiresAfter int64) OperatorMessage {
	expireTime := time.Now().Unix() + expiresAfter
	if expiresAfter < 0 {
		expireTime = -1
	}
	return OperatorMessage{
		strconv.Itoa(int(id)),
		content,
		item,
		expireTime,
	}
}
