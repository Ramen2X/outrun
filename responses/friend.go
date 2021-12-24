package responses

import (
	"github.com/Ramen2X/outrun/netobj"
	"github.com/Ramen2X/outrun/obj"
	"github.com/Ramen2X/outrun/responses/responseobjs"
)

type FacebookIncentiveResponse struct {
	BaseResponse
	PlayerState netobj.PlayerState `json:"playerState"`
	Presents    []obj.Present      `json:"incentive"`
}

func FacebookIncentive(base responseobjs.BaseInfo, playerState netobj.PlayerState, presents []obj.Present) FacebookIncentiveResponse {
	baseResponse := NewBaseResponse(base)
	return FacebookIncentiveResponse{
		baseResponse,
		playerState,
		presents,
	}
}

func DefaultFacebookIncentive(base responseobjs.BaseInfo, player netobj.Player) FacebookIncentiveResponse {
	playerState := player.PlayerState
	presents := []obj.Present{} // Naughty this year
	return FacebookIncentive(
		base,
		playerState,
		presents,
	)
}
