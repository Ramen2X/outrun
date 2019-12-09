package responses

import (
	"time"

	"github.com/fluofoxxo/outrun/obj"
	"github.com/fluofoxxo/outrun/responses/responseobjs"
)

type DailyBattleDataResponse struct {
	BaseResponse
	EndTime      int64          `json:"endTime"`
	BattleStatus obj.BattleData `json:"battleData"`
}

func DailyBattleData(base responseobjs.BaseInfo, endTime int64, battleData obj.BattleData) DailyBattleDataResponse {
	baseResponse := NewBaseResponse(base)
	return DailyBattleDataResponse{
		baseResponse,
		endTime,
		battleData,
	}
}

func DefaultDailyBattleData(base responseobjs.BaseInfo) DailyBattleDataResponse {
	battleData := obj.DebugRivalBattleData()
	return DailyBattleData(
		base,
		time.Now().Unix()+80000, // ~22 hours from now
		battleData,
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
