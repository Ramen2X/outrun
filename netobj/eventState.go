package netobj

type EventState struct {
	Param           int64 `json:"param"`
	RewardID        int64 `json:"rewardId"`
	PreviousEventID int64 `json:"ORN_prevEventId"`
}

func DefaultEventState() EventState {
	param := int64(0)
	rewardId := int64(0)        // ???
	previousEventId := int64(0) // no previous event
	return NewEventState(param, rewardId, previousEventId)
}

func NewEventState(param, rewardId, previousEventId int64) EventState {
	return EventState{
		param,
		rewardId,
		previousEventId,
	}
}
