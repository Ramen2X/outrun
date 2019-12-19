package requests

type GetDailyBattleHistoryRequest struct {
	Count int64 `json:"count,string"`
}

type ResetDailyBattleMatchingRequest struct {
	Type int64 `json:"type,string"`
}

// ResetDailyBattleMatching has 3 possible types:
// 0: Initial search
// 1: Normal search (costs 5 red rings)
// 2: Closest match search (costs 10 red rings)
