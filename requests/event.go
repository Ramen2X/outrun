package requests

type GenericEventRequest struct {
	Base
	EventID int64 `json:"eventId,string"`
}

type EventActStartRequest struct {
	Base
	Modifier     []int64 `json:"modifire"` // Seems to be list of item IDs.
	RaidbossID   int64   `json:"raidbossId"`
	EventID      int64   `json:"eventId"`
	EnergyExpend int64   `json:"energyExpend"` // the amount of raidboss energy to be used?
}
