package responses

import (
	"github.com/fluofoxxo/outrun/logic/conversion"
	"github.com/fluofoxxo/outrun/netobj"
	"github.com/fluofoxxo/outrun/obj"
	"github.com/fluofoxxo/outrun/responses/responseobjs"
	"github.com/jinzhu/now"
)

type NoScoreDailyBattleDataResponse struct {
	BaseResponse
	StartTime int64 `json:"startTime"`
	EndTime   int64 `json:"endTime"`
}

type NoRivalDailyBattleDataResponse struct {
	BaseResponse
	StartTime  int64          `json:"startTime"`
	EndTime    int64          `json:"endTime"`
	BattleData obj.BattleData `json:"battleData"`
}

type DailyBattleDataResponse struct {
	BaseResponse
	obj.BattlePair
}

func NoScoreDailyBattleData(base responseobjs.BaseInfo, startTime, endTime int64) NoScoreDailyBattleDataResponse {
	baseResponse := NewBaseResponse(base)
	return NoScoreDailyBattleDataResponse{
		baseResponse,
		startTime,
		endTime,
	}
}

func NoRivalDailyBattleData(base responseobjs.BaseInfo, startTime, endTime int64, battleData obj.BattleData) NoRivalDailyBattleDataResponse {
	baseResponse := NewBaseResponse(base)
	return NoRivalDailyBattleDataResponse{
		baseResponse,
		startTime,
		endTime,
		battleData,
	}
}

func DailyBattleData(base responseobjs.BaseInfo, startTime, endTime int64, battleData, rivalBattleData obj.BattleData) DailyBattleDataResponse {
	baseResponse := NewBaseResponse(base)
	return DailyBattleDataResponse{
		baseResponse,
		obj.NewBattlePair(startTime, endTime, battleData, rivalBattleData),
	}
}

func DefaultDailyBattleData(base responseobjs.BaseInfo, player netobj.Player) NoRivalDailyBattleDataResponse {
	battleData := conversion.DebugPlayerToBattleData(player)
	//	rivalBattleData := obj.DebugRivalBattleData()
	return NoRivalDailyBattleData(
		base,
		now.BeginningOfDay().UTC().Unix(),
		now.EndOfDay().UTC().Unix(),
		battleData,
		//		rivalBattleData,
	)
}

func DefaultMatchedDailyBattleData(base responseobjs.BaseInfo, player netobj.Player) DailyBattleDataResponse {
	battleData := conversion.DebugPlayerToBattleData(player)
	rivalBattleData := obj.DebugRivalBattleData()
	return DailyBattleData(
		base,
		now.BeginningOfDay().UTC().Unix(),
		now.EndOfDay().UTC().Unix(),
		battleData,
		rivalBattleData,
	)
}

type UpdateDailyBattleStatusResponse struct {
	BaseResponse
	EndTime      int64            `json:"endTime"`
	BattleStatus obj.BattleStatus `json:"battleStatus"`
	RewardFlag   bool             `json:"rewardFlag"` // TODO: allow not false after testing
}

type UpdateDailyBattleStatusResponseWithReward struct {
	BaseResponse
	EndTime      int64            `json:"endTime"`
	BattleStatus obj.BattleStatus `json:"battleStatus"`
	RewardFlag   bool             `json:"rewardFlag"` // TODO: allow not false after testing
	obj.RewardBattlePair
}

func UpdateDailyBattleStatus(base responseobjs.BaseInfo, endTime int64, battleStatus obj.BattleStatus) UpdateDailyBattleStatusResponse {
	baseResponse := NewBaseResponse(base)
	return UpdateDailyBattleStatusResponse{
		baseResponse,
		endTime,
		battleStatus,
		false,
	}
}

func UpdateDailyBattleStatusWithReward(base responseobjs.BaseInfo, endTime int64, battleStatus obj.BattleStatus, rewardStartTime, rewardEndTime int64, rewardBattleData, rewardRivalBattleData obj.BattleData) UpdateDailyBattleStatusResponseWithReward {
	baseResponse := NewBaseResponse(base)
	battleReward := obj.NewRewardBattlePair(rewardStartTime, rewardEndTime, rewardBattleData, rewardRivalBattleData)
	return UpdateDailyBattleStatusResponseWithReward{
		baseResponse,
		endTime,
		battleStatus,
		true,
		battleReward,
	}
}

type ResetDailyBattleMatchingResponse struct {
	BaseResponse
	obj.BattlePair
	PlayerState netobj.PlayerState `json:"playerState"`
}

type ResetDailyBattleMatchingNoOpponentResponse struct {
	BaseResponse
	StartTime   int64              `json:"startTime"`
	EndTime     int64              `json:"endTime"`
	BattleData  obj.BattleData     `json:"battleData"`
	PlayerState netobj.PlayerState `json:"playerState"`
}

func ResetDailyBattleMatching(base responseobjs.BaseInfo, startTime, endTime int64, battleData, rivalBattleData obj.BattleData, player netobj.Player) ResetDailyBattleMatchingResponse {
	baseResponse := NewBaseResponse(base)
	return ResetDailyBattleMatchingResponse{
		baseResponse,
		obj.NewBattlePair(startTime, endTime, battleData, rivalBattleData),
		player.PlayerState,
	}
}

func ResetDailyBattleMatchingNoOpponent(base responseobjs.BaseInfo, startTime, endTime int64, battleData obj.BattleData, player netobj.Player) ResetDailyBattleMatchingNoOpponentResponse {
	baseResponse := NewBaseResponse(base)
	return ResetDailyBattleMatchingNoOpponentResponse{
		baseResponse,
		startTime,
		endTime,
		battleData,
		player.PlayerState,
	}
}

type GetDailyBattleHistoryResponse struct {
	BaseResponse
	BattleDataList []obj.BattlePair `json:"battleDataList"`
}

func GetDailyBattleHistory(base responseobjs.BaseInfo, battleDataList []obj.BattlePair) GetDailyBattleHistoryResponse {
	baseResponse := NewBaseResponse(base)
	return GetDailyBattleHistoryResponse{
		baseResponse,
		battleDataList,
	}
}

type GetDailyBattleStatusResponse struct {
	BaseResponse
	EndTime      int64            `json:"endTime"`
	BattleStatus obj.BattleStatus `json:"battleStatus"`
}

func GetDailyBattleStatus(base responseobjs.BaseInfo, endTime int64, battleStatus obj.BattleStatus) GetDailyBattleStatusResponse {
	baseResponse := NewBaseResponse(base)
	return GetDailyBattleStatusResponse{
		baseResponse,
		endTime,
		battleStatus,
	}
}

type PostDailyBattleResultResponse struct {
	BaseResponse
	obj.BattlePair
	BattleStatus obj.BattleStatus `json:"battleStatus"`
	RewardFlag   bool             `json:"rewardFlag"` // TODO: allow not false after testing
}

type PostDailyBattleResultResponseNoRival struct {
	BaseResponse
	StartTime    int64            `json:"startTime"`
	EndTime      int64            `json:"endTime"`
	BattleData   obj.BattleData   `json:"battleData"`
	BattleStatus obj.BattleStatus `json:"battleStatus"`
	RewardFlag   bool             `json:"rewardFlag"` // TODO: allow not false after testing
}

type PostDailyBattleResultResponseNoData struct {
	BaseResponse
	StartTime    int64            `json:"startTime"`
	EndTime      int64            `json:"endTime"`
	BattleStatus obj.BattleStatus `json:"battleStatus"`
	RewardFlag   bool             `json:"rewardFlag"` // TODO: allow not false after testing
}

type PostDailyBattleResultResponseWithReward struct {
	BaseResponse
	StartTime    int64            `json:"startTime"`
	EndTime      int64            `json:"endTime"`
	BattleStatus obj.BattleStatus `json:"battleStatus"`
	RewardFlag   bool             `json:"rewardFlag"` // TODO: allow not false after testing
	obj.RewardBattlePair
}

func PostDailyBattleResult(base responseobjs.BaseInfo, startTime, endTime int64, battleData, rivalBattleData obj.BattleData, battleStatus obj.BattleStatus) PostDailyBattleResultResponse {
	baseResponse := NewBaseResponse(base)
	return PostDailyBattleResultResponse{
		baseResponse,
		obj.NewBattlePair(startTime, endTime, battleData, rivalBattleData),
		battleStatus,
		false,
	}
}

func PostDailyBattleResultNoData(base responseobjs.BaseInfo, startTime, endTime int64, battleStatus obj.BattleStatus) PostDailyBattleResultResponseNoData {
	baseResponse := NewBaseResponse(base)
	return PostDailyBattleResultResponseNoData{
		baseResponse,
		startTime,
		endTime,
		battleStatus,
		false,
	}
}

func PostDailyBattleResultNoRival(base responseobjs.BaseInfo, startTime, endTime int64, battleData obj.BattleData, battleStatus obj.BattleStatus) PostDailyBattleResultResponseNoRival {
	baseResponse := NewBaseResponse(base)
	return PostDailyBattleResultResponseNoRival{
		baseResponse,
		startTime,
		endTime,
		battleData,
		battleStatus,
		false,
	}
}

func PostDailyBattleResultWithReward(base responseobjs.BaseInfo, startTime, endTime int64, battleStatus obj.BattleStatus, rewardStartTime, rewardEndTime int64, rewardBattleData, rewardRivalBattleData obj.BattleData) PostDailyBattleResultResponseWithReward {
	baseResponse := NewBaseResponse(base)
	battleReward := obj.NewRewardBattlePair(rewardStartTime, rewardEndTime, rewardBattleData, rewardRivalBattleData)
	return PostDailyBattleResultResponseWithReward{
		baseResponse,
		startTime,
		endTime,
		battleStatus,
		true,
		battleReward,
	}
}
