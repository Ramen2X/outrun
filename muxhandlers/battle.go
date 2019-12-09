package muxhandlers

import (
	"encoding/json"
	"time"

	"github.com/fluofoxxo/outrun/emess"
	"github.com/fluofoxxo/outrun/helper"
	"github.com/fluofoxxo/outrun/obj"
	"github.com/fluofoxxo/outrun/requests"
	"github.com/fluofoxxo/outrun/responses"
	"github.com/fluofoxxo/outrun/status"
)

func GetDailyBattleData(helper *helper.Helper) {
	// TODO: Right now, send agnostic data. In reality, this definitely should be player based!
	/*
	   player, err := helper.GetCallingPlayer()
	   if err != nil {
	       helper.InternalErr("error getting calling player", err)
	       return
	   }
	*/
	baseInfo := helper.BaseInfo(emess.OK, status.OK)
	response := responses.DefaultDailyBattleData(baseInfo)
	err := helper.SendCompatibleResponse(response)
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
	endTime := time.Now().UTC().Unix() + 180 // three minutes from now, for testing
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
