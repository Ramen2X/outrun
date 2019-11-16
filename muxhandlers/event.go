package muxhandlers

import (
	"encoding/json"
	"fmt"

	"github.com/fluofoxxo/outrun/config"
	"github.com/fluofoxxo/outrun/config/eventconf"
	"github.com/fluofoxxo/outrun/db"
	"github.com/fluofoxxo/outrun/emess"
	"github.com/fluofoxxo/outrun/helper"
	"github.com/fluofoxxo/outrun/logic/conversion"
	"github.com/fluofoxxo/outrun/netobj"
	"github.com/fluofoxxo/outrun/obj"
	"github.com/fluofoxxo/outrun/requests"
	"github.com/fluofoxxo/outrun/responses"
	"github.com/fluofoxxo/outrun/status"
)

func GetEventList(helper *helper.Helper) {
	player, err := helper.GetCallingPlayer()
	if err != nil {
		helper.InternalErr("Error getting calling player", err)
		return
	}
	baseInfo := helper.BaseInfo(emess.OK, status.OK)
	// construct event list
	eventList := []obj.Event{}
	if eventconf.CFile.AllowEvents {
		if eventconf.CFile.EnforceGlobal || len(player.PersonalEvents) == 0 {
			for _, confEvent := range eventconf.CFile.CurrentEvents {
				newEvent := conversion.ConfiguredEventToEvent(confEvent)
				eventList = append(eventList, newEvent)
			}
		} else {
			for _, ce := range player.PersonalEvents {
				e := conversion.ConfiguredEventToEvent(ce)
				eventList = append(eventList, e)
			}
		}
	}
	if config.CFile.DebugPrints {
		helper.Out("Personal event list: " + fmt.Sprintln(player.PersonalEvents))
		helper.Out("Global event list: " + fmt.Sprintln(eventconf.CFile.CurrentEvents))
		helper.Out("Event list: " + fmt.Sprintln(eventList))
	}
	response := responses.EventList(baseInfo, eventList)
	err = helper.SendResponse(response)
	if err != nil {
		helper.InternalErr("Error sending response", err)
	}
}

func GetEventReward(helper *helper.Helper) {
	recv := helper.GetGameRequest()
	var request requests.GetEventRewardRequest
	err := json.Unmarshal(recv, &request)
	if err != nil {
		helper.Err("Error unmarshalling", err)
		return
	}
	baseInfo := helper.BaseInfo(emess.OK, status.OK)
	response := responses.DefaultEventRewardList(baseInfo)
	err = helper.SendResponse(response)
	if err != nil {
		helper.InternalErr("Error sending response", err)
	}
}

func GetEventState(helper *helper.Helper) {
	recv := helper.GetGameRequest()
	var request requests.GetEventStateRequest
	err := json.Unmarshal(recv, &request)
	if err != nil {
		helper.Err("Error unmarshalling", err)
		return
	}
	player, err := helper.GetCallingPlayer()
	if err != nil {
		helper.InternalErr("Error getting calling player", err)
		return
	}
	if request.EventID != player.EventState.PreviousEventID {
		player.EventState.Param = 0 //reset values
		player.EventState.RewardID = 0
		player.EventState.PreviousEventID = request.EventID
	}
	baseInfo := helper.BaseInfo(emess.OK, status.OK)
	response := responses.EventState(baseInfo, player.EventState)	
	err = db.SavePlayer(player)
	if err != nil {
		helper.InternalErr("Error saving player", err)
		return
	}
	err = helper.SendResponse(response)
	if err != nil {
		helper.InternalErr("Error sending response", err)
	}
}

// 1.1.4 raid bosses
func GetEventUserRaidbossState(helper *helper.Helper) {
	baseInfo := helper.BaseInfo(emess.OK, status.OK)
	response := responses.EventUserRaidbossState(baseInfo, netobj.DefaultUserRaidbossState())
	err := helper.SendResponse(response)
	if err != nil {
		helper.InternalErr("Error sending response", err)
	}
}
