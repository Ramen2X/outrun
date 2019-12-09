package muxhandlers

import (
	"encoding/json"

	"github.com/fluofoxxo/outrun/emess"
	"github.com/fluofoxxo/outrun/helper"
	"github.com/fluofoxxo/outrun/obj"
	"github.com/fluofoxxo/outrun/requests"
	"github.com/fluofoxxo/outrun/responses"
	"github.com/fluofoxxo/outrun/status"
	"github.com/jinzhu/now"
)

func GetDailyBattleData(helper *helper.Helper) {
	player, err := helper.GetCallingPlayer()
	if err != nil {
		helper.InternalErr("error getting calling player", err)
		return
	}
	baseInfo := helper.BaseInfo(emess.OK, status.OK)
	response := responses.DefaultDailyBattleData(baseInfo, player)
	err = helper.SendCompatibleResponse(response)
	if err != nil {
		helper.InternalErr("error sending response", err)
	}
}

func UpdateDailyBattleStatus(helper *helper.Helper) {
	data := helper.GetGameRequest()
	var request requests.Base
	err := json.Unmarshal(data, &request)
	if err != nil {
		helper.InternalErr("Error unmarshalling", err)
		return
	}
	endTime := now.EndOfDay().UTC().Unix()
	rewardFlag := false
	battleStatus := obj.DefaultBattleStatus()
	baseInfo := helper.BaseInfo(emess.OK, status.OK)
	response := responses.UpdateDailyBattleStatus(baseInfo, endTime, battleStatus, rewardFlag)
	err = helper.SendCompatibleResponse(response)
	if err != nil {
		helper.InternalErr("Error sending response", err)
		return
	}
}
