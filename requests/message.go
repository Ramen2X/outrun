package requests

type GetMessageRequest struct {
	Base
	MessageIDs         interface{} `json:"messageId"`          // can either be a list of int64s or "0"
	OperatorMessageIDs interface{} `json:"operationMessageId"` // can either be a list of int64s or "0"
}
