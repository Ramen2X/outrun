package muxhandlers

import (
	"encoding/json"
	"time"

	"github.com/fluofoxxo/outrun/db"
	"github.com/fluofoxxo/outrun/emess"
	"github.com/fluofoxxo/outrun/helper"
	"github.com/fluofoxxo/outrun/logic/battle"
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
			helper.DebugOut("No rival")
			response = responses.NoRivalDailyBattleData(baseInfo,
				player.BattleState.BattleStartsAt,
				player.BattleState.BattleEndsAt,
				conversion.DebugPlayerToBattleData(player),
			)
		}
	} else {
		helper.DebugOut("No score recorded")
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
	if player.BattleState.PendingReward {
		rewardBattleStartTime = player.BattleState.PendingRewardData.StartTime
		rewardBattleEndTime = player.BattleState.PendingRewardData.EndTime
		rewardBattlePlayerData = player.BattleState.PendingRewardData.BattleData
		rewardBattleRivalData = player.BattleState.PendingRewardData.RivalBattleData
		player.BattleState.PendingReward = false
		doReward = true
		if time.Now().UTC().Unix() > player.BattleState.BattleEndsAt {
			player.BattleState.BattleStartsAt = now.BeginningOfDay().UTC().Unix()
			player.BattleState.BattleEndsAt = now.EndOfDay().UTC().Unix()
			player.BattleState.ScoreRecordedToday = false
			player.BattleState.MatchedUpWithRival = false
		}
	} else {
		if time.Now().UTC().Unix() > player.BattleState.BattleEndsAt {
			if player.BattleState.ScoreRecordedToday {
				if player.BattleState.MatchedUpWithRival {
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
					rivalBattlePair := obj.NewBattlePair(
						rewardBattleStartTime,
						rewardBattleEndTime,
						rewardBattleRivalData,
						rewardBattlePlayerData,
					)
					rivalPlayer.BattleState.PendingReward = true
					rivalPlayer.BattleState.PendingRewardData = obj.NewRewardBattlePair(
						rewardBattleStartTime,
						rewardBattleEndTime,
						rewardBattleRivalData,
						rewardBattlePlayerData,
					)
					if player.BattleState.DailyBattleHighScore > rivalPlayer.BattleState.DailyBattleHighScore {
						player.BattleState.Wins++
						player.BattleState.WinStreak++
						player.BattleState.LossStreak = 0
						rivalPlayer.BattleState.Losses++
						rivalPlayer.BattleState.LossStreak++
						rivalPlayer.BattleState.WinStreak = 0
						// Then we'd send the appropriate things to the gift boxes, but...
						// TODO: Add the reward functionality
					} else {
						if player.BattleState.DailyBattleHighScore < rivalPlayer.BattleState.DailyBattleHighScore {
							player.BattleState.Losses++
							player.BattleState.LossStreak++
							player.BattleState.WinStreak = 0
							rivalPlayer.BattleState.Wins++
							rivalPlayer.BattleState.WinStreak++
							rivalPlayer.BattleState.LossStreak = 0
							// Then we'd send the appropriate things to the gift boxes, but...
							// TODO: Add the reward functionality
						} else {
							player.BattleState.Draws++
							player.BattleState.WinStreak = 0
							player.BattleState.LossStreak = 0
							rivalPlayer.BattleState.Draws++
							rivalPlayer.BattleState.WinStreak = 0
							rivalPlayer.BattleState.LossStreak = 0
							// Then we'd send the appropriate things to the gift boxes, but...
							// TODO: Add the reward functionality
						}
					}
					player.BattleState.BattleHistory = append(player.BattleState.BattleHistory, battlePair)
					rivalPlayer.BattleState.BattleHistory = append(rivalPlayer.BattleState.BattleHistory, rivalBattlePair)
					err = db.SavePlayer(rivalPlayer)
					if err != nil {
						helper.InternalErr("Error saving player", err)
						return
					}
					doReward = true
				} else {
					// There appears to be no reward for failures
					// TODO: Is that right?
					player.BattleState.Failures++
					player.BattleState.LossStreak++
					player.BattleState.WinStreak = 0
				}
			}
			player.BattleState.BattleStartsAt = now.BeginningOfDay().UTC().Unix()
			player.BattleState.BattleEndsAt = now.EndOfDay().UTC().Unix()
			player.BattleState.ScoreRecordedToday = false
			player.BattleState.MatchedUpWithRival = false
		}
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
	err = db.SavePlayer(player)
	if err != nil {
		helper.InternalErr("Error saving player", err)
		return
	}
}

// Reroll daily battle rival
func ResetDailyBattleMatching(helper *helper.Helper) {
	data := helper.GetGameRequest()
	var request requests.ResetDailyBattleMatchingRequest
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
	baseInfo := helper.BaseInfo(emess.OK, status.OK)
	battleData := conversion.DebugPlayerToBattleData(player)
	startTime := player.BattleState.BattleStartsAt
	endTime := player.BattleState.BattleEndsAt

	helper.DebugOut("Type: %v", request.Type)
	switch request.Type {
	case 1:
		if player.PlayerState.NumRedRings < 5 {
			baseInfo.StatusCode = status.NotEnoughRedRings
			err = helper.SendResponse(responses.NewBaseResponse(baseInfo))
			if err != nil {
				helper.InternalErr("error sending response", err)
			}
			return
		}
	case 2:
		if player.PlayerState.NumRedRings < 10 {
			baseInfo.StatusCode = status.NotEnoughRedRings
			err = helper.SendResponse(responses.NewBaseResponse(baseInfo))
			if err != nil {
				helper.InternalErr("error sending response", err)
			}
			return
		}
	}
	oldRivalID := player.BattleState.RivalID
	oldRival, err := db.GetPlayer(oldRivalID)
	if err != nil {
		helper.InternalErr("error getting rival player", err)
		return
	}
	player.BattleState.MatchedUpWithRival = false
	oldRival.BattleState.MatchedUpWithRival = false
	err = db.SavePlayer(oldRival)
	if err != nil {
		helper.InternalErr("Error saving old rival", err)
		return
	}
	if request.Type == 2 {
		helper.InvalidRequest()
		return
	} else {
		player.BattleState = battle.DrawBattleRival(player)
	}

	if player.BattleState.RivalID != oldRivalID && player.BattleState.MatchedUpWithRival {
		switch request.Type {
		case 1:
			player.PlayerState.NumRedRings -= 5
		case 2:
			player.PlayerState.NumRedRings -= 10
		}
	}

	var response interface{}
	if player.BattleState.MatchedUpWithRival {
		rivalPlayer, err := db.GetPlayer(player.BattleState.RivalID)
		if err != nil {
			helper.InternalErr("error getting rival player", err)
			return
		}
		rivalBattleData := conversion.DebugPlayerToBattleData(rivalPlayer)
		response = responses.ResetDailyBattleMatching(baseInfo, startTime, endTime, battleData, rivalBattleData, player)
	} else {
		response = responses.ResetDailyBattleMatchingNoOpponent(baseInfo, startTime, endTime, battleData, player)
	}
	err = helper.SendCompatibleResponse(response)
	if err != nil {
		helper.InternalErr("error sending response", err)
	}
	err = db.SavePlayer(player)
	if err != nil {
		helper.InternalErr("Error saving player", err)
		return
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

func PostDailyBattleResult(helper *helper.Helper) {
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
	if player.BattleState.PendingReward {
		rewardBattleStartTime = player.BattleState.PendingRewardData.StartTime
		rewardBattleEndTime = player.BattleState.PendingRewardData.EndTime
		rewardBattlePlayerData = player.BattleState.PendingRewardData.BattleData
		rewardBattleRivalData = player.BattleState.PendingRewardData.RivalBattleData
		player.BattleState.PendingReward = false
		doReward = true
		if time.Now().UTC().Unix() > player.BattleState.BattleEndsAt {
			player.BattleState.BattleStartsAt = now.BeginningOfDay().UTC().Unix()
			player.BattleState.BattleEndsAt = now.EndOfDay().UTC().Unix()
			player.BattleState.ScoreRecordedToday = false
			player.BattleState.MatchedUpWithRival = false
		}
	} else {
		if time.Now().UTC().Unix() > player.BattleState.BattleEndsAt {
			if player.BattleState.ScoreRecordedToday {
				if player.BattleState.MatchedUpWithRival {
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
					rivalBattlePair := obj.NewBattlePair(
						rewardBattleStartTime,
						rewardBattleEndTime,
						rewardBattleRivalData,
						rewardBattlePlayerData,
					)
					rivalPlayer.BattleState.PendingReward = true
					rivalPlayer.BattleState.PendingRewardData = obj.NewRewardBattlePair(
						rewardBattleStartTime,
						rewardBattleEndTime,
						rewardBattleRivalData,
						rewardBattlePlayerData,
					)
					if player.BattleState.DailyBattleHighScore > rivalPlayer.BattleState.DailyBattleHighScore {
						player.BattleState.Wins++
						player.BattleState.WinStreak++
						player.BattleState.LossStreak = 0
						rivalPlayer.BattleState.Losses++
						rivalPlayer.BattleState.LossStreak++
						rivalPlayer.BattleState.WinStreak = 0
						// Then we'd send the appropriate things to the gift boxes, but...
						// TODO: Add the reward functionality
					} else {
						if player.BattleState.DailyBattleHighScore < rivalPlayer.BattleState.DailyBattleHighScore {
							player.BattleState.Losses++
							player.BattleState.LossStreak++
							player.BattleState.WinStreak = 0
							rivalPlayer.BattleState.Wins++
							rivalPlayer.BattleState.WinStreak++
							rivalPlayer.BattleState.LossStreak = 0
							// Then we'd send the appropriate things to the gift boxes, but...
							// TODO: Add the reward functionality
						} else {
							player.BattleState.Draws++
							player.BattleState.WinStreak = 0
							player.BattleState.LossStreak = 0
							rivalPlayer.BattleState.Draws++
							rivalPlayer.BattleState.WinStreak = 0
							rivalPlayer.BattleState.LossStreak = 0
							// Then we'd send the appropriate things to the gift boxes, but...
							// TODO: Add the reward functionality
						}
					}
					player.BattleState.BattleHistory = append(player.BattleState.BattleHistory, battlePair)
					rivalPlayer.BattleState.BattleHistory = append(rivalPlayer.BattleState.BattleHistory, rivalBattlePair)
					err = db.SavePlayer(rivalPlayer)
					if err != nil {
						helper.InternalErr("Error saving player", err)
						return
					}
					doReward = true
				} else {
					// There appears to be no reward for failures
					// TODO: Is that right?
					player.BattleState.Failures++
					player.BattleState.LossStreak++
					player.BattleState.WinStreak = 0
				}
			}
			player.BattleState.BattleStartsAt = now.BeginningOfDay().UTC().Unix()
			player.BattleState.BattleEndsAt = now.EndOfDay().UTC().Unix()
			player.BattleState.ScoreRecordedToday = false
			player.BattleState.MatchedUpWithRival = false
		}
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
		response = responses.PostDailyBattleResultWithReward(baseInfo,
			player.BattleState.BattleStartsAt,
			player.BattleState.BattleEndsAt,
			battleStatus,
			rewardBattleStartTime,
			rewardBattleEndTime,
			rewardBattlePlayerData,
			rewardBattleRivalData,
		)
	} else {
		if player.BattleState.ScoreRecordedToday {
			if player.BattleState.MatchedUpWithRival {
				rivalPlayerData, err := db.GetPlayer(player.BattleState.RivalID)
				if err != nil {
					helper.InternalErr("error getting rival player", err)
					return
				}
				response = responses.PostDailyBattleResult(baseInfo,
					player.BattleState.BattleStartsAt,
					player.BattleState.BattleEndsAt,
					conversion.DebugPlayerToBattleData(player),
					conversion.DebugPlayerToBattleData(rivalPlayerData),
					battleStatus,
				)
			} else {
				helper.DebugOut("No rival")
				response = responses.PostDailyBattleResultNoRival(baseInfo,
					player.BattleState.BattleStartsAt,
					player.BattleState.BattleEndsAt,
					conversion.DebugPlayerToBattleData(player),
					battleStatus,
				)
			}
		} else {
			helper.DebugOut("No score recorded")
			response = responses.PostDailyBattleResultNoData(baseInfo,
				player.BattleState.BattleStartsAt,
				player.BattleState.BattleEndsAt,
				battleStatus,
			)
		}
	}

	err = helper.SendCompatibleResponse(response)
	if err != nil {
		helper.InternalErr("Error sending response", err)
		return
	}
	err = db.SavePlayer(player)
	if err != nil {
		helper.InternalErr("Error saving player", err)
		return
	}
}
