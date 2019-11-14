package responses

import (
	"strconv"
	
	"github.com/fluofoxxo/outrun/enums"
	"github.com/fluofoxxo/outrun/netobj"
	"github.com/fluofoxxo/outrun/obj"
	"github.com/fluofoxxo/outrun/responses/responseobjs"
)

type EventListResponse struct {
	BaseResponse
	EventList []obj.Event `json:"eventList"`
}

func EventList(base responseobjs.BaseInfo, eventList []obj.Event) EventListResponse {
	baseResponse := NewBaseResponse(base)
	out := EventListResponse{
		baseResponse,
		eventList,
	}
	return out
}

func DefaultEventList(base responseobjs.BaseInfo) EventListResponse {
	return EventList(
		base,
		[]obj.Event{
			/*
			   obj.NewEvent(
			       //enums.EventIDSpecialStage+10002, // game subtracts one from number?
			       //enums.EventIDAdvert+50002, // 50002 converts to ui_event_50005_Atlas_en?
			       //enums.EventIDBGM+70002, // 70002 goes to 70007
			       enums.EventIDQuick+60002, // 60002 goes to 60006
			       0,                        // event type
			       now.BeginningOfDay().Unix(),
			       now.EndOfDay().Unix(),
			       now.EndOfDay().Unix(),
			   ),
			*/
		},
	)
}

// 1.1.4 support
type EventRewardListResponse struct {
	BaseResponse
	EventRewardList []obj.EventReward `json:"eventRewardList"`
}

func EventRewardList(base responseobjs.BaseInfo, eventRewardList []obj.EventReward) EventRewardListResponse {
	baseResponse := NewBaseResponse(base)
	out := EventRewardListResponse{
		baseResponse,
		eventRewardList,
	}
	return out
}

func DefaultEventRewardList(base responseobjs.BaseInfo) EventRewardListResponse {
	//TODO: Get this from the config, and/or on a per-event basis
	return EventRewardList(
		base,
		[]obj.EventReward{
			obj.NewEventReward(
				1,
				1500,
				strconv.Itoa(int(enums.ItemIDAsteroid)),
				10,
			),
			obj.NewEventReward(
				2,
				3000,
				strconv.Itoa(int(enums.ItemIDTrampoline)),
				10,
			),
			obj.NewEventReward(
				3,
				5000,
				strconv.Itoa(int(enums.ItemIDDrill)),
				10,
			),
			obj.NewEventReward(
				4,
				7500,
				strconv.Itoa(int(enums.ItemIDLaser)),
				15,
			),
			obj.NewEventReward(
				5,
				11000,
				strconv.Itoa(int(enums.ItemIDInvincible)),
				15,
			),
			obj.NewEventReward(
				6,
				20000,
				strconv.Itoa(int(enums.ItemIDRing)),
				50000,
			),
			obj.NewEventReward(
				7,
				30000,
				strconv.Itoa(int(enums.ItemIDRedRing)),
				500,
			),
			obj.NewEventReward(
				8,
				45000,
				strconv.Itoa(int(enums.ChaoIDKingBoomBoo)),
				1,
			),
		},
	)
}

type EventStateResponse struct {
	BaseResponse
	netobj.EventState `json:"eventState"`
}

func EventState(base responseobjs.BaseInfo, eventState netobj.EventState) EventStateResponse {
	baseResponse := NewBaseResponse(base)
	out := EventStateResponse{
		baseResponse,
		eventState,
	}
	return out
}
