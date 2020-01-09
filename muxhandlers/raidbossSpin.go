package muxhandlers

import (
	"github.com/fluofoxxo/outrun/emess"
	"github.com/fluofoxxo/outrun/enums"
	"github.com/fluofoxxo/outrun/helper"
	"github.com/fluofoxxo/outrun/netobj"
	"github.com/fluofoxxo/outrun/responses"
	"github.com/fluofoxxo/outrun/status"
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
	player, err := helper.GetCallingPlayer()
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
	err = helper.SendCompatibleResponse(response)
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
