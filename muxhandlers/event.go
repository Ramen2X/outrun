package muxhandlers

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/fluofoxxo/outrun/config/eventconf"
	"github.com/fluofoxxo/outrun/config/gameconf"
	"github.com/fluofoxxo/outrun/db"
	"github.com/fluofoxxo/outrun/emess"
	"github.com/fluofoxxo/outrun/helper"
	"github.com/fluofoxxo/outrun/logic/conversion"
	"github.com/fluofoxxo/outrun/logic/gameplay"
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
	helper.DebugOut("Personal event list: %v", player.PersonalEvents)
	helper.DebugOut("Global event list: %v", eventconf.CFile.CurrentEvents)
	helper.DebugOut("Event list: %v", eventList)
	response := responses.EventList(baseInfo, eventList)
	err = helper.SendResponse(response)
	if err != nil {
		helper.InternalErr("Error sending response", err)
	}
}

func GetEventReward(helper *helper.Helper) {
	recv := helper.GetGameRequest()
	var request requests.GenericEventRequest
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
	var request requests.GenericEventRequest
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
	recv := helper.GetGameRequest()
	var request requests.GenericEventRequest
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
	baseInfo := helper.BaseInfo(emess.OK, status.OK)
	response := responses.EventUserRaidbossState(baseInfo, player.EventUserRaidbossState)
	response.Seq = request.Seq
	response.Version = request.Version
	err = helper.SendInsecureResponse(response)
	if err != nil {
		helper.InternalErr("Error sending response", err)
	}
}

func EventActStart(helper *helper.Helper) {
	recv := helper.GetGameRequest()
	var request requests.EventActStartRequest
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
	helper.DebugOut(fmt.Sprintf("%v", player.PlayerState.Items))
	responseStatus := status.OK
	// consume items
	modToStringSlice := func(ns []int64) []string {
		result := []string{}
		for _, n := range ns {
			result = append(result, fmt.Sprintf("%v", n))
		}
		return result
	}
	helper.DebugOut(fmt.Sprintf("%v", player.PlayerState.Items))
	for time.Now().UTC().Unix() >= player.EventUserRaidbossState.EnergyRenewsAt && player.EventUserRaidbossState.RaidBossEnergy < 20 {
		player.EventUserRaidbossState.RaidBossEnergy++
		player.EventUserRaidbossState.EnergyRenewsAt += 1200
	}
	if player.PlayerState.Energy+player.PlayerState.EnergyBuy >= request.EnergyExpend {
		if gameconf.CFile.EnableEnergyConsumption {
			if player.PlayerState.EnergyBuy > 0 {
				player.EventUserRaidbossState.RaidBossEnergyBuy -= request.EnergyExpend
				if player.EventUserRaidbossState.RaidBossEnergyBuy < 0 { //did we go negative?
					player.EventUserRaidbossState.RaidBossEnergy += player.EventUserRaidbossState.RaidBossEnergyBuy
					player.EventUserRaidbossState.RaidBossEnergyBuy = 0
				}
			} else {
				player.PlayerState.Energy -= request.EnergyExpend
				if player.EventUserRaidbossState.RaidBossEnergy < 20 {
					player.PlayerState.EnergyRenewsAt = time.Now().UTC().Unix() + 1200
				}
			}
		}
		player.PlayerState.NumPlaying++
		if !gameconf.CFile.AllItemsFree {
			consumedItems := modToStringSlice(request.Modifier)
			consumedRings := gameplay.GetRequiredItemPayment(consumedItems)
			for _, citemID := range consumedItems {
				if citemID[:2] == "11" { // boosts, not items
					continue
				}
				index := player.IndexOfItem(citemID)
				if index == -1 {
					helper.Uncatchable(fmt.Sprintf("Player sent bad item ID '%s', cannot continue", citemID))
					helper.InvalidRequest()
					return
				}
				if player.PlayerState.Items[index].Amount >= 1 { // can use item
					player.PlayerState.Items[index].Amount--
				} else {
					if player.PlayerState.NumRings < consumedRings { // not enough rings
						responseStatus = status.NotEnoughRings
						break
					}
					player.PlayerState.NumRings -= consumedRings
				}
			}
		}
	} else {
		responseStatus = status.NotEnoughEnergy
	}
	baseInfo := helper.BaseInfo(emess.OK, responseStatus)
	respPlayer := player
	if request.Version == "1.1.4" { // must send fewer characters
		// only get first 21 characters
		// TODO: enforce order 300000 to 300020?
		//cState = cState[:len(cState)-(len(cState)-10)]
		cState := respPlayer.CharacterState
		cState = cState[:16]
		helper.DebugOut("cState length: " + strconv.Itoa(len(cState)))
		helper.DebugOut("Sent character IDs: ")
		for _, char := range cState {
			helper.DebugOut(char.ID)
		}
		respPlayer.CharacterState = cState
	}
	response := responses.DefaultEventActStart(baseInfo, respPlayer)
	err = helper.SendResponse(response)
	if err != nil {
		helper.InternalErr("Error sending response", err)
		return
	}
	err = db.SavePlayer(player)
	if err != nil {
		helper.InternalErr("Error saving player", err)
		return
	}
	//TODO: Add analytics for this
}
