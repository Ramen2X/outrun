package netobj

import "github.com/jinzhu/now"

type BattleState struct {
	ScoreRecordedToday       bool   `json:"hasRecordedScoreToday"`
	DailyBattleHighScore     int64  `json:"maxScore"`
	PrevDailyBattleHighScore int64  `json:"lastMaxScore"`
	BattleEndsAt             int64  `json:"expireTime"`
	MatchedUpWithRival       bool   `json:"matchedUpWithRival"`
	RivalID                  string `json:"rivalId"`
	Wins                     int64  `json:"numWin"`
	Losses                   int64  `json:"numLose"`
	Draws                    int64  `json:"numDraw"`
	Failures                 int64  `json:"numLoseByDefault"`
	WinStreak                int64  `json:"goOnWin"`
	LossStreak               int64  `json:"goOnLosses"`
}

func NewBattleState(scoreRecordedToday bool, dailyBattleHighScore, prevDailyBattleHighScore, battleEndTime int64, matchedUpWithRival bool, rivalID string, wins, losses, draws, failures, winStreak, lossStreak int64) BattleState {
	return BattleState{
		scoreRecordedToday,
		dailyBattleHighScore,
		prevDailyBattleHighScore,
		battleEndTime,
		matchedUpWithRival,
		rivalID,
		wins,
		losses,
		draws,
		failures,
		winStreak,
		lossStreak,
	}
}

func DefaultBattleState() BattleState {
	scoreRecordedToday := false
	dailyBattleHighScore := int64(0)
	prevDailyBattleHighScore := int64(0)
	battleEndTime := now.EndOfDay().UTC().Unix()
	matchedUpWithRival := false
	rivalID := ""
	wins := int64(0)
	losses := int64(0)
	draws := int64(0)
	failures := int64(0)
	winStreak := int64(0)
	lossStreak := int64(0)
	return NewBattleState(
		scoreRecordedToday,
		dailyBattleHighScore,
		prevDailyBattleHighScore,
		battleEndTime,
		matchedUpWithRival,
		rivalID,
		wins,
		losses,
		draws,
		failures,
		winStreak,
		lossStreak,
	)
}
