package obj

type BattlePair struct { // This is just used for organization within the response
	StartTime       int64      `json:"startTime"`
	EndTime         int64      `json:"endTime"`
	BattleData      BattleData `json:"battleData"`
	RivalBattleData BattleData `json:"rivalBattleData"`
}

func NewBattlePair(startTime, endTime int64, battleData, rivalBattleData BattleData) BattlePair {
	return BattlePair{
		startTime,
		endTime,
		battleData,
		rivalBattleData,
	}
}
