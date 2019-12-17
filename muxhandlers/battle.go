package muxhandlers

import (
	"encoding/json"
	"time"

	"github.com/fluofoxxo/outrun/db"
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
		if player.BattleState.MatchedUpWithRival {
			rivalPlayer, err := db.GetPlayer(player.BattleState.RivalID)
			if err != nil {
				helper.InternalErr("error getting rival player", err)
				return
			}
			response = responses.DailyBattleData(baseInfo,
				player.BattleState.BattleStartsAt,
				player.BattleState.BattleEndsAt,
				conversion.DebugPlayerToBattleData(player),
				conversion.DebugPlayerToBattleData(rivalPlayer),
			)
		} else {
			response = responses.NoRivalDailyBattleData(baseInfo,
				player.BattleState.BattleStartsAt,
				player.BattleState.BattleEndsAt,
				conversion.DebugPlayerToBattleData(player),
			)
		}
	} else {
		response = responses.NoScoreDailyBattleData(baseInfo,
			player.BattleState.BattleStartsAt,
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
	player, err := helper.GetCallingPlayer()
	if err != nil {
		helper.InternalErr("error getting calling player", err)
		return
	}
	var rewardBattleStartTime int64
	var rewardBattleEndTime int64
	var rewardBattlePlayerData obj.BattleData
	var rewardBattleRivalData obj.BattleData
	doReward := false
	if time.Now().UTC().Unix() > player.BattleState.BattleEndsAt {
		if player.BattleState.ScoreRecordedToday && player.BattleState.MatchedUpWithRival {
			rivalPlayer, err := db.GetPlayer(player.BattleState.RivalID)
			if err != nil {
				helper.InternalErr("error getting rival player", err)
				return
			}
			rewardBattleStartTime = player.BattleState.BattleStartsAt
			rewardBattleEndTime = player.BattleState.BattleEndsAt
			rewardBattlePlayerData = conversion.DebugPlayerToBattleData(player)
			rewardBattleRivalData = conversion.DebugPlayerToBattleData(rivalPlayer)
			battlePair := obj.NewBattlePair(
				rewardBattleStartTime,
				rewardBattleEndTime,
				rewardBattlePlayerData,
				rewardBattleRivalData,
			)
			if player.BattleState.DailyBattleHighScore > rivalPlayer.BattleState.DailyBattleHighScore {
				player.BattleState.Wins++
				player.BattleState.WinStreak++
				player.BattleState.LossStreak = 0
			} else {
				if player.BattleState.DailyBattleHighScore < rivalPlayer.BattleState.DailyBattleHighScore {
					player.BattleState.Losses++
					player.BattleState.LossStreak++
					player.BattleState.WinStreak = 0
				} else {
					player.BattleState.Draws++
					player.BattleState.WinStreak = 0
					player.BattleState.LossStreak = 0
				}
			}
			player.BattleState.BattleHistory = append(player.BattleState.BattleHistory, battlePair)
			doReward = true
		}
		player.BattleState.BattleStartsAt = now.BeginningOfDay().UTC().Unix()
		player.BattleState.BattleEndsAt = now.EndOfDay().UTC().Unix()
		player.BattleState.ScoreRecordedToday = false
		player.BattleState.MatchedUpWithRival = false
	}
	battleStatus := obj.BattleStatus{
		player.BattleState.Wins,
		player.BattleState.Losses,
		player.BattleState.Draws,
		player.BattleState.Failures,
		player.BattleState.WinStreak,
		player.BattleState.LossStreak,
	}
	baseInfo := helper.BaseInfo(emess.OK, status.OK)
	var response interface{}
	if doReward {
		response = responses.UpdateDailyBattleStatusWithReward(baseInfo, player.BattleState.BattleEndsAt, battleStatus, rewardBattleStartTime, rewardBattleEndTime, rewardBattlePlayerData, rewardBattleRivalData)
	} else {
		response = responses.UpdateDailyBattleStatus(baseInfo, player.BattleState.BattleEndsAt, battleStatus)
	}

	err = helper.SendCompatibleResponse(response)
	if err != nil {
		helper.InternalErr("Error sending response", err)
		return
	}
}

// Reroll daily battle rival
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

func GetDailyBattleHistory(helper *helper.Helper) {
	player, err := helper.GetCallingPlayer()
	if err != nil {
		helper.InternalErr("error getting calling player", err)
		return
	}
	baseInfo := helper.BaseInfo(emess.OK, status.OK)
	response := responses.GetDailyBattleHistory(baseInfo, player.BattleState.BattleHistory)
	err = helper.SendResponse(response)
	if err != nil {
		helper.InternalErr("error sending response", err)
	}
}

func GetDailyBattleStatus(helper *helper.Helper) {
	data := helper.GetGameRequest()
	var request requests.Base
	err := json.Unmarshal(data, &request)
	if err != nil {
		helper.InternalErr("Error unmarshalling", err)
		return
	}
	player, err := helper.GetCallingPlayer()
	if err != nil {
		helper.InternalErr("error getting calling player", err)
		return
	}
	battleStatus := obj.BattleStatus{
		player.BattleState.Wins,
		player.BattleState.Losses,
		player.BattleState.Draws,
		player.BattleState.Failures,
		player.BattleState.WinStreak,
		player.BattleState.LossStreak,
	}
	baseInfo := helper.BaseInfo(emess.OK, status.OK)

	response := responses.GetDailyBattleStatus(baseInfo, player.BattleState.BattleEndsAt, battleStatus)
	err = helper.SendCompatibleResponse(response)
	if err != nil {
		helper.InternalErr("Error sending response", err)
		return
	}
}
