package requests

import "github.com/Ramen2X/outrun/netobj"

type QuickPostGameResultsRequest struct {
	Base
	Score                  int64  `json:"score,string"`
	Rings                  int64  `json:"numRings,string"`
	FailureRings           int64  `json:"numFailureRings,string"`
	RedRings               int64  `json:"numRedStarRings,string"`
	Distance               int64  `json:"distance,string"`
	DailyChallengeValue    int64  `json:"dailyChallengeValue,string"`
	DailyChallengeComplete int64  `json:"dailyChallengeComplete,string"`
	Animals                int64  `json:"numAnimals,string"`
	MaxCombo               int64  `json:"maxCombo,string"`
	Closed                 int64  `json:"closed,string"`
	CheatResult            string `json:"cheatResult"`
}

type PostGameResultsRequest struct {
	QuickPostGameResultsRequest
	BossDestroyed int64 `json:"bossDestroyed"`
	ChapterClear  int64 `json:"chapterClear,string"`
	GetChaoEgg    int64 `json:"getChaoEgg,string"`
	NumBossAttack int64 `json:"numBossAttack,string"`
	ReachPoint    int64 `json:"reachPoint,string"`
	EventId       int64 `json:"eventId,string"`
	EventValue    int64 `json:"eventValue,string"`
}

type QuickActStartRequest struct {
	Base
	Modifier []int64 `json:"modifire"`           // Seems to be list of item IDs.
	Tutorial int64   `json:"tutorial,string"` // will omit the field if not found (breaks 1.0.0 for some reason)
}

type ActStartRequest struct {
	QuickActStartRequest
	DistanceFriendList []netobj.MileageFriend `json:"distanceFriendList"` // TODO: Discover correct type... This might be list of strings
}

type MileageRewardRequest struct {
	Base
	Episode int64 `json:"episode,string"`
	Chapter int64 `json:"chapter,string"`
}

type DrawRaidBossRequest struct {
	Base
	EventID int64 `json:"eventId,string"`
	Score   int64 `json:"score,string"`
}
