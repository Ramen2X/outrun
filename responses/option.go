package responses

import (
	"github.com/Ramen2X/outrun/netobj"
	"github.com/Ramen2X/outrun/responses/responseobjs"
)

type OptionUserResultResponse struct {
	BaseResponse
	OptionUserResult netobj.OptionUserResult `json:"optionUserResult"`
}

func OptionUserResult(base responseobjs.BaseInfo, optionUserResult netobj.OptionUserResult) OptionUserResultResponse {
	baseResponse := NewBaseResponse(base)
	out := OptionUserResultResponse{
		baseResponse,
		optionUserResult,
	}
	return out
}
