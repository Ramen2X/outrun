package responses

import (
	"github.com/Ramen2X/outrun/netobj"
	"github.com/Ramen2X/outrun/obj"
	"github.com/Ramen2X/outrun/responses/responseobjs"
)

type MessageListResponse struct {
	BaseResponse
	MessageList           []obj.Message         `json:"messageList"`
	TotalMessages         int64                 `json:"totalMessage"`
	OperatorMessageList   []obj.OperatorMessage `json:"operatorMessageList"`
	TotalOperatorMessages int64                 `json:"totalOperatorMessage"`
}

func MessageList(base responseobjs.BaseInfo, msgl []obj.Message, opmsgl []obj.OperatorMessage) MessageListResponse {
	baseResponse := NewBaseResponse(base)
	out := MessageListResponse{
		baseResponse,
		msgl,
		int64(len(msgl)),
		opmsgl,
		int64(len(opmsgl)),
	}
	return out
}

func DefaultMessageList(base responseobjs.BaseInfo) MessageListResponse {
	return MessageList(
		base,
		[]obj.Message{},
		[]obj.OperatorMessage{
			obj.DefaultOperatorMessage(),
		},
	)
}

type GetMessageResponse struct {
	BaseResponse
	PlayerState                 netobj.PlayerState `json:"playerState"`
	CharacterState              []netobj.Character `json:"characterState"`
	ChaoState                   []netobj.Chao      `json:"chaoState"`
	PresentList                 []obj.Present      `json:"presentList"`                // obtained gifts?
	RemainingMessageIDs         []int64            `json:"notRecvMessageList"`         // IDs of messages not yet received?
	RemainingOperatorMessageIDs []int64            `json:"notRecvOperatorMessageList"` // IDs of operator messages not yet received?
}

func GetMessage(base responseobjs.BaseInfo, player netobj.Player, presentList []obj.Present, remainingMessageIDs, remainingOperatorMessageIDs []int64) GetMessageResponse {
	baseResponse := NewBaseResponse(base)
	out := GetMessageResponse{
		baseResponse,
		player.PlayerState,
		player.CharacterState,
		player.ChaoState,
		presentList,
		remainingMessageIDs,
		remainingOperatorMessageIDs,
	}
	return out
}
