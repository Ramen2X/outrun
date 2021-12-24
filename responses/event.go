package responses

import (
	//	"strconv"

	//	"github.com/Ramen2X/outrun/enums"
	"github.com/Ramen2X/outrun/logic"
	"github.com/Ramen2X/outrun/netobj"
	"github.com/Ramen2X/outrun/obj"
	"github.com/Ramen2X/outrun/responses/responseobjs"
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
			/*obj.NewEventReward(
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
			),*/
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

type EventUserRaidbossStateResponse struct {
	BaseResponse
	netobj.EventUserRaidbossState `json:"eventUserRaidboss"`
}

func EventUserRaidbossState(base responseobjs.BaseInfo, userRaidbossState netobj.EventUserRaidbossState) EventUserRaidbossStateResponse {
	baseResponse := NewBaseResponse(base)
	out := EventUserRaidbossStateResponse{
		baseResponse,
		userRaidbossState,
	}
	return out
}

type EventUserRaidbossListResponse struct {
	BaseResponse
	netobj.EventUserRaidbossState `json:"eventUserRaidboss"`
	EventRaidbossStates           []netobj.EventRaidbossState `json:"eventUserRaidbossList"`
}

func DefaultEventUserRaidbossList(base responseobjs.BaseInfo, userRaidbossState netobj.EventUserRaidbossState) EventUserRaidbossListResponse {
	baseResponse := NewBaseResponse(base)
	out := EventUserRaidbossListResponse{
		baseResponse,
		userRaidbossState,
		[]netobj.EventRaidbossState{
			netobj.DefaultRaidbossState(),
		},
	}
	return out
}

type EventActStartResponse struct {
	ActStartBaseResponse
	netobj.EventUserRaidbossState `json:"eventUserRaidboss"`
}

func EventActStart(base responseobjs.BaseInfo, playerState netobj.PlayerState, campaignList []obj.Campaign, eventUserRaidbossState netobj.EventUserRaidbossState) EventActStartResponse {
	actStartBase := ActStartBase(base, playerState, campaignList)
	return EventActStartResponse{
		actStartBase,
		eventUserRaidbossState,
	}
}

func DefaultEventActStart(base responseobjs.BaseInfo, player netobj.Player) EventActStartResponse {
	campaignList := obj.DefaultCampaigns()
	playerState := player.PlayerState
	eventUserRaidbossState := player.EventUserRaidbossState
	return EventActStart(
		base,
		playerState,
		campaignList,
		eventUserRaidbossState,
	)
}

type EventPostGameResultsResponse struct {
	BaseResponse
	netobj.EventUserRaidbossState `json:"eventUserRaidboss"`
}

func EventPostGameResults(base responseobjs.BaseInfo, userRaidbossState netobj.EventUserRaidbossState) EventPostGameResultsResponse {
	baseResponse := NewBaseResponse(base)
	out := EventPostGameResultsResponse{
		baseResponse,
		userRaidbossState,
	}
	return out
}

type EventUpdateGameResultsResponse struct {
	BaseResponse
	PlayerState             netobj.PlayerState    `json:"playerState"`
	ChaoState               []netobj.Chao         `json:"chaoState"`
	DailyChallengeIncentive []obj.Incentive       `json:"dailyChallengeIncentive"` // should be obj.Item, but game doesn't care
	CharacterState          []netobj.Character    `json:"characterState"`
	MessageList             []obj.Message         `json:"messageList"`
	OperatorMessageList     []obj.OperatorMessage `json:"operatorMessageList"`
	TotalMessage            int64                 `json:"totalMessage"`
	TotalOperatorMessage    int64                 `json:"totalOperatorMessage"`
	PlayCharacterState      []netobj.Character    `json:"playCharacterState"`
	EventIncentiveList      []obj.Item            `json:"eventIncentiveList"`
	WheelOptions            netobj.WheelOptions   `json:"wheelOptions"`
	EventState              netobj.EventState     `json:"eventState,omitempty"`
}

func EventUpdateGameResults(base responseobjs.BaseInfo, player netobj.Player, dci []obj.Incentive, ml []obj.Message, oml []obj.OperatorMessage, pcs []netobj.Character, eil []obj.Item, wo netobj.WheelOptions, es netobj.EventState) EventUpdateGameResultsResponse {
	baseResponse := NewBaseResponse(base)
	playerState := player.PlayerState
	chaoState := player.ChaoState
	dailyChallengeIncentive := dci
	characterState := player.CharacterState
	messageList := []obj.Message{}
	operatorMessageList := []obj.OperatorMessage{}
	totalMessage := int64(len(messageList))
	totalOperatorMessage := int64(len(operatorMessageList))
	playCharacterState := pcs
	return EventUpdateGameResultsResponse{
		baseResponse,
		playerState,
		chaoState,
		dailyChallengeIncentive,
		characterState,
		messageList,
		operatorMessageList,
		totalMessage,
		totalOperatorMessage,
		playCharacterState,
		eil,
		wo,
		es,
	}
}

func DefaultEventUpdateGameResults(base responseobjs.BaseInfo, player netobj.Player, pcs []netobj.Character, es netobj.EventState) EventUpdateGameResultsResponse {
	baseResponse := NewBaseResponse(base)
	playerState := player.PlayerState
	chaoState := player.ChaoState
	dailyChallengeIncentive := []obj.Incentive{}
	characterState := player.CharacterState
	messageList := []obj.Message{}
	operatorMessageList := []obj.OperatorMessage{}
	totalMessage := int64(len(messageList))
	totalOperatorMessage := int64(len(operatorMessageList))
	eil := []obj.Item{}
	player.LastWheelOptions = logic.WheelRefreshLogic(player, player.LastWheelOptions)
	wo := player.LastWheelOptions
	return EventUpdateGameResultsResponse{
		baseResponse,
		playerState,
		chaoState,
		dailyChallengeIncentive,
		characterState,
		messageList,
		operatorMessageList,
		totalMessage,
		totalOperatorMessage,
		pcs,
		eil,
		wo,
		es,
	}
}
