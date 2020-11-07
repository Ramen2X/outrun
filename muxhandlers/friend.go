package muxhandlers

import (
	"encoding/json"

	"github.com/fluofoxxo/outrun/emess"
	"github.com/fluofoxxo/outrun/helper"
	"github.com/fluofoxxo/outrun/requests"
	"github.com/fluofoxxo/outrun/responses"
	"github.com/fluofoxxo/outrun/status"
)

func GetFacebookIncentive(helper *helper.Helper) {
	data := helper.GetGameRequest()
	var request requests.FacebookIncentiveRequest
	err := json.Unmarshal(data, &request)
	if err != nil {
		helper.InternalErr("Error unmarshalling", err)
		return
	}
	switch request.Type {
	case 0:
		helper.DebugOut("Type 0 - LOGIN")
		break
	case 1:
		helper.DebugOut("User accepted the \"leave a review\" prompt!")
		break
	case 2:
		helper.DebugOut("Type 2 - FEED")
		break
	case 3:
		helper.DebugOut("Type 3 - ACHIEVEMENT")
		break
	case 4:
		helper.DebugOut("Type 4 - PUSH_NOLOGIN")
		break
	default:
		helper.DebugOut("Unknown incentive type %v", request.Type)
		break
	}
	// We respond with no presents for now.
	player, err := helper.GetCallingPlayer(true)
	if err != nil {
		helper.InternalErr("Error getting calling player", err)
		return
	}
	baseInfo := helper.BaseInfo(emess.OK, status.OK)
	response := responses.DefaultFacebookIncentive(baseInfo, player)
	err = helper.SendResponse(response)
	if err != nil {
		helper.InternalErr("Error sending response", err)
	}
}
