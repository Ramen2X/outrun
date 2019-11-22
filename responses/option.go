package responses

import (
	"github.com/fluofoxxo/outrun/netobj"
	"github.com/fluofoxxo/outrun/responses/responseobjs"
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
