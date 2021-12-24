package responses

import (
	"github.com/Ramen2X/outrun/logic/conversion"
	"github.com/Ramen2X/outrun/netobj"
	"github.com/Ramen2X/outrun/obj"
	"github.com/Ramen2X/outrun/obj/constobjs"
	"github.com/Ramen2X/outrun/responses/responseobjs"
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

type UpdateDailyBattleStatusResponseWithRewardNoRival struct {
	BaseResponse
	EndTime          int64            `json:"endTime"`
	BattleStatus     obj.BattleStatus `json:"battleStatus"`
	RewardFlag       bool             `json:"rewardFlag"` // TODO: allow not false after testing
	RewardStartTime  int64            `json:"rewardStartTime"`
	RewardEndTime    int64            `json:"rewardEndTime"`
	RewardBattleData obj.BattleData   `json:"rewardBattleData"`
}

type UpdateDailyBattleStatusResponseWithRewardNoData struct {
	BaseResponse
	EndTime         int64            `json:"endTime"`
	BattleStatus    obj.BattleStatus `json:"battleStatus"`
	RewardFlag      bool             `json:"rewardFlag"` // TODO: allow not false after testing
	RewardStartTime int64            `json:"rewardStartTime"`
	RewardEndTime   int64            `json:"rewardEndTime"`
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

func UpdateDailyBattleStatusWithReward(base responseobjs.BaseInfo, endTime int64, battleStatus obj.BattleStatus, rewardStartTime, rewardEndTime int64, rewardBattleData, rewardRivalBattleData interface{}) interface{} {
	baseResponse := NewBaseResponse(base)
	if rewardBattleData != nil {
		if rewardRivalBattleData != nil {
			battleReward := obj.NewRewardBattlePair(rewardStartTime, rewardEndTime, rewardBattleData.(obj.BattleData), rewardRivalBattleData.(obj.BattleData))
			return UpdateDailyBattleStatusResponseWithReward{
				baseResponse,
				endTime,
				battleStatus,
				true,
				battleReward,
			}
		} else {
			return UpdateDailyBattleStatusResponseWithRewardNoRival{
				baseResponse,
				endTime,
				battleStatus,
				true,
				rewardStartTime,
				rewardEndTime,
				rewardBattleData.(obj.BattleData),
			}
		}
	} else {
		return UpdateDailyBattleStatusResponseWithRewardNoData{
			baseResponse,
			endTime,
			battleStatus,
			true,
			rewardStartTime,
			rewardEndTime,
		}
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

type PostDailyBattleResultResponseWithRewardNoRival struct {
	BaseResponse
	StartTime        int64            `json:"startTime"`
	EndTime          int64            `json:"endTime"`
	BattleStatus     obj.BattleStatus `json:"battleStatus"`
	RewardFlag       bool             `json:"rewardFlag"` // TODO: allow not false after testing
	RewardStartTime  int64            `json:"rewardStartTime"`
	RewardEndTime    int64            `json:"rewardEndTime"`
	RewardBattleData obj.BattleData   `json:"rewardBattleData"`
}

type PostDailyBattleResultResponseWithRewardNoData struct {
	BaseResponse
	StartTime       int64            `json:"startTime"`
	EndTime         int64            `json:"endTime"`
	BattleStatus    obj.BattleStatus `json:"battleStatus"`
	RewardFlag      bool             `json:"rewardFlag"` // TODO: allow not false after testing
	RewardStartTime int64            `json:"rewardStartTime"`
	RewardEndTime   int64            `json:"rewardEndTime"`
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

func PostDailyBattleResultWithReward(base responseobjs.BaseInfo, startTime, endTime int64, battleStatus obj.BattleStatus, rewardStartTime, rewardEndTime int64, rewardBattleData, rewardRivalBattleData interface{}) interface{} {
	baseResponse := NewBaseResponse(base)
	if rewardBattleData != nil {
		if rewardRivalBattleData != nil {
			battleReward := obj.NewRewardBattlePair(rewardStartTime, rewardEndTime, rewardBattleData.(obj.BattleData), rewardRivalBattleData.(obj.BattleData))
			return PostDailyBattleResultResponseWithReward{
				baseResponse,
				startTime,
				endTime,
				battleStatus,
				true,
				battleReward,
			}
		} else {
			return PostDailyBattleResultResponseWithRewardNoRival{
				baseResponse,
				startTime,
				endTime,
				battleStatus,
				true,
				rewardStartTime,
				rewardEndTime,
				rewardBattleData.(obj.BattleData),
			}
		}
	} else {
		return PostDailyBattleResultResponseWithRewardNoData{
			baseResponse,
			startTime,
			endTime,
			battleStatus,
			true,
			rewardStartTime,
			rewardEndTime,
		}
	}
}

type GetPrizeDailyBattleResponse struct {
	BaseResponse
	BattlePrizeDataList []obj.OperatorScore `json:"battlePrizeDataList"`
}

func GetPrizeDailyBattle(base responseobjs.BaseInfo, battlePrizeDataList []obj.OperatorScore) GetPrizeDailyBattleResponse {
	baseResponse := NewBaseResponse(base)
	return GetPrizeDailyBattleResponse{
		baseResponse,
		battlePrizeDataList,
	}
}

func DefaultGetPrizeDailyBattle(base responseobjs.BaseInfo) GetPrizeDailyBattleResponse {
	return GetPrizeDailyBattle(
		base,
		constobjs.DefaultDailyBattlePrizeList,
	)
}
