package responses

import (
	"github.com/fluofoxxo/outrun/logic/conversion"
	"github.com/fluofoxxo/outrun/netobj"
	"github.com/fluofoxxo/outrun/obj"
	"github.com/fluofoxxo/outrun/responses/responseobjs"
	"github.com/jinzhu/now"
)

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

type UpdateDailyBattleStatusResponse struct {
	BaseResponse
	EndTime      int64            `json:"endTime"`
	BattleStatus obj.BattleStatus `json:"battleStatus"`
	RewardFlag   bool             `json:"rewardFlag"` // TODO: allow not false after testing
}

func UpdateDailyBattleStatus(base responseobjs.BaseInfo, endTime int64, battleStatus obj.BattleStatus, rewardFlag bool) UpdateDailyBattleStatusResponse {
	baseResponse := NewBaseResponse(base)
	return UpdateDailyBattleStatusResponse{
		baseResponse,
		endTime,
		battleStatus,
		rewardFlag,
	}
}
