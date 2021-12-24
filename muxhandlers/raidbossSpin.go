package muxhandlers

import (
	"github.com/Ramen2X/outrun/emess"
	"github.com/Ramen2X/outrun/enums"
	"github.com/Ramen2X/outrun/helper"
	"github.com/Ramen2X/outrun/netobj"
	"github.com/Ramen2X/outrun/responses"
	"github.com/Ramen2X/outrun/status"
)

func GetItemStockNum(helper *helper.Helper) {
	// TODO: Flesh out properly! The game responds with
	// [IDRouletteTicketPremium, IDRouletteTicketItem, IDSpecialEgg]
	// for item IDs, along with an event ID, likely for event characters.
	baseInfo := helper.BaseInfo(emess.OK, status.OK)
	response := responses.DefaultItemStockNum(baseInfo)
	err := helper.SendResponse(response)
	if err != nil {
		helper.InternalErr("Error sending response", err)
	}
}

//1.1.4
func GetRaidbossWheelOptions(helper *helper.Helper) {
	player, err := helper.GetCallingPlayer(true)
	if err != nil {
		helper.InternalErr("Error getting calling player", err)
		return
	}
	baseInfo := helper.BaseInfo(emess.OK, status.OK)
	if player.Suspended {
		baseInfo.StatusCode = status.MissingPlayer
		err = helper.SendResponse(responses.NewBaseResponse(baseInfo))
		if err != nil {
			helper.InternalErr("Error sending response", err)
		}
		return
	}
	wheelOptions := netobj.DefaultRaidbossWheelOptions(0, player.PlayerState.ChaoEggs, 0, enums.WheelRankNormal, 0)
	response := responses.RaidbossWheelOptions(baseInfo, wheelOptions)
	err = helper.SendCompatibleResponse(response, true)
	if err != nil {
		helper.InternalErr("Error sending response", err)
	}
}

func GetPrizeRaidbossWheelSpin(helper *helper.Helper) {
	// agnostic
	baseInfo := helper.BaseInfo(emess.OK, status.OK)
	response := responses.DefaultPrizeRaidbossWheel(baseInfo)
	err := helper.SendResponse(response)
	if err != nil {
		helper.InternalErr("Error sending response", err)
	}
}
