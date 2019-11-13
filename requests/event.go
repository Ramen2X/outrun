package requests

type GetEventRewardRequest struct {
	Base
	EventID int64  `json:"eventId,string"`
}

type GetEventStateRequest struct {
	Base
	EventID int64  `json:"eventId,string"`
}
