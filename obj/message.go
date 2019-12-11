package obj

type Message struct {
	MessageID   string      `json:"messageId"`
	MessageType int64       `json:"messageType"`
	FriendID    string      `json:"friendId"`
	Name        string      `json:"name"`
	URL         string      `json:"url"`
	Item        MessageItem `json:"item"`
	ExpireTime  int64       `json:"expireTime"`
}
