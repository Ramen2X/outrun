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
	// TODO: is agnostic; shouldn't be!
	wheelOptions := netobj.DefaultRaidbossWheelOptions(7, 0, 0, enums.WheelRankNormal, 0)
	baseInfo := helper.BaseInfo(emess.OK, status.OK)
	response := responses.RaidbossWheelOptions(baseInfo, wheelOptions)
	err := helper.SendCompatibleResponse(response)
	if err != nil {
		helper.InternalErr("Error sending response", err)
	}
}

func GetPrizeRaidbossWheelSpin(helper *helper.Helper) {
	// agnostic
	baseInfo := helper.BaseInfo(emess.OK, status.OK)
	response := responses.DefaultPrizeChaoWheel(baseInfo)
	err := helper.SendResponse(response)
	if err != nil {
		helper.InternalErr("Error sending response", err)
	}
}
