package obj

type BattleStatus struct {
	Wins            int64 `json:"numWin"`           // battle wins
	Losses          int64 `json:"numLose"`          // battle losses
	Draws           int64 `json:"numDraw"`          // battle draws
	LossesByDefault int64 `json:"numLoseByDefault"` // battle failures
	GoOnWins        int64 `json:"goOnWin"`          // battle win streak
	GoOnLosses      int64 `json:"goOnLosses"`       // battle loss streak???
}

func DefaultBattleStatus() BattleStatus {
	return BattleStatus{
		0,
		0,
		0,
		0,
		0,
		0,
	}
}
