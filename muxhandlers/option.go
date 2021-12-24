package muxhandlers

import (
	"github.com/Ramen2X/outrun/emess"
	"github.com/Ramen2X/outrun/helper"
	"github.com/Ramen2X/outrun/responses"
	"github.com/Ramen2X/outrun/status"
)

func GetOptionUserResult(helper *helper.Helper) {
	player, err := helper.GetCallingPlayer(true)
	if err != nil {
		helper.InternalErr("Error getting calling player", err)
		return
	}
	baseInfo := helper.BaseInfo(emess.OK, status.OK)
	response := responses.OptionUserResult(baseInfo, player.OptionUserResult)
	helper.SendResponse(response)
}
