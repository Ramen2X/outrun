package muxhandlers

import (
	"github.com/fluofoxxo/outrun/emess"
	"github.com/fluofoxxo/outrun/helper"
	"github.com/fluofoxxo/outrun/netobj"
	"github.com/fluofoxxo/outrun/responses"
	"github.com/fluofoxxo/outrun/status"
)

func GetOptionUserResult(helper *helper.Helper) {
	//TODO: Get this info from player
	/*player, err := helper.GetCallingPlayer()
	if err != nil {
		helper.InternalErr("Error getting calling player", err)
		return
	}*/
	baseInfo := helper.BaseInfo(emess.OK, status.OK)
	response := responses.OptionUserResult(baseInfo, netobj.DefaultOptionUserResult())
	helper.SendResponse(response)
}
