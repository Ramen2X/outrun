package requests

type GenericEventRequest struct {
	Base
	EventID int64 `json:"eventId,string"`
}

type EventActStartRequest struct {
	Base
	Modifier     []string `json:"modifire"` // Seems to be list of item IDs.
	RaidbossID   int64    `json:"raidbossId,string"`
	EventID      int64    `json:"eventId,string"`
	EnergyExpend int64    `json:"energyExpend,string"` // the amount of raidboss energy to be used?
}

type EventPostGameResultsRequest struct {
	Base
	EventID          int64 `json:"eventId,string"`
	NumRaidbossRings int64 `json:"numRaidbossRings,string"`
}

type EventUpdateGameResultsRequest struct {
	QuickPostGameResultsRequest
	EventID         int64 `json:"eventId,string"`
	EventValue      int64 `json:"eventValue,string"`
	RaidbossID      int64 `json:"raidbossId,string"`
	RaidbossDamage  int64 `json:"raidbossDamage,string"`
	RaidbossBeatFlg int64 `json:"raidbossBeatFlg,string"`
}
