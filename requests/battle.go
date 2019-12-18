package requests

type GetDailyBattleHistoryRequest struct {
	Count int64 `json:"count"`
}

type ResetDailyBattleMatchingRequest struct {
	Type int64 `json:"type"`
}
