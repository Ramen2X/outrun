package muxhandlers

import (
	"encoding/json"

	"github.com/fluofoxxo/outrun/emess"
	"github.com/fluofoxxo/outrun/helper"
	"github.com/fluofoxxo/outrun/logic/conversion"
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
	var response interface{}
	if player.BattleState.ScoreRecordedToday {
		response = responses.DefaultDailyBattleData(baseInfo, player)
	} else {
		response = responses.NoScoreDailyBattleData(baseInfo,
			now.BeginningOfDay().UTC().Unix(),
			player.BattleState.BattleEndsAt,
		)
	}
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
	/*player, err := helper.GetCallingPlayer()
	if err != nil {
		helper.InternalErr("error getting calling player", err)
		return
	}*/
	endTime := now.EndOfDay().UTC().Unix()
	//rewardBattleData := conversion.DebugPlayerToBattleData(player)
	//rewardRivalBattleData := obj.DebugRivalBattleData()
	//rewardStartTime := now.BeginningOfDay().UTC().Unix()
	//rewardEndTime := now.EndOfDay().UTC().Unix()
	battleStatus := obj.DefaultBattleStatus()
	baseInfo := helper.BaseInfo(emess.OK, status.OK)

	//response := responses.UpdateDailyBattleStatusWithReward(baseInfo, endTime, battleStatus, rewardStartTime, rewardEndTime, rewardBattleData, rewardRivalBattleData)
	response := responses.UpdateDailyBattleStatus(baseInfo, endTime, battleStatus)
	err = helper.SendCompatibleResponse(response)
	if err != nil {
		helper.InternalErr("Error sending response", err)
		return
	}
}

func ResetDailyBattleMatching(helper *helper.Helper) {
	player, err := helper.GetCallingPlayer()
	if err != nil {
		helper.InternalErr("error getting calling player", err)
		return
	}
	baseInfo := helper.BaseInfo(emess.OK, status.OK)
	battleData := conversion.DebugPlayerToBattleData(player)
	//rivalBattleData := obj.DebugRivalBattleData()
	startTime := now.BeginningOfDay().UTC().Unix()
	endTime := now.EndOfDay().UTC().Unix()
	response := responses.ResetDailyBattleMatchingNoOpponent(baseInfo, startTime, endTime, battleData, player)
	err = helper.SendCompatibleResponse(response)
	if err != nil {
		helper.InternalErr("error sending response", err)
	}
}
