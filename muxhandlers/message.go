package muxhandlers

import (
	"encoding/json"
	"strconv"

	"github.com/Ramen2X/outrun/enums"
	"github.com/Ramen2X/outrun/logic"
	"github.com/Ramen2X/outrun/netobj"

	"github.com/Ramen2X/outrun/db"
	"github.com/Ramen2X/outrun/emess"
	"github.com/Ramen2X/outrun/helper"
	"github.com/Ramen2X/outrun/obj"
	"github.com/Ramen2X/outrun/requests"
	"github.com/Ramen2X/outrun/responses"
	"github.com/Ramen2X/outrun/status"
)

func GetMessageList(helper *helper.Helper) {
	player, err := helper.GetCallingPlayer(true)
	if err != nil {
		helper.InternalErr("error getting calling player", err)
		return
	}
	baseInfo := helper.BaseInfo(emess.OK, status.OK)
	if player.OperatorMessages == nil {
		player.OperatorMessages = []obj.OperatorMessage{}
	}
	db.SavePlayer(player)
	// response := responses.DefaultMessageList(baseInfo)
	response := responses.MessageList(baseInfo, []obj.Message{}, player.OperatorMessages)
	err = helper.SendCompatibleResponse(response, true)
	if err != nil {
		helper.InternalErr("Error sending response", err)
	}
}

func GetMessage(helper *helper.Helper) {
	data := helper.GetGameRequest()
	var request requests.GetMessageRequest
	err := json.Unmarshal(data, &request)
	if err != nil {
		helper.InternalErr("Error unmarshalling", err)
		return
	}
	player, err := helper.GetCallingPlayer(true)
	if err != nil {
		helper.InternalErr("error getting calling player", err)
		return
	}
	baseInfo := helper.BaseInfo(emess.OK, status.OK)
	if player.OperatorMessages == nil {
		player.OperatorMessages = []obj.OperatorMessage{}
	}

	presentList := []obj.Present{}
	blankParams := 0

	switch messageIds := request.MessageIDs.(type) {
	case []interface{}:
		helper.DebugOut("%v", messageIds)
		for _, msgid := range messageIds {
			helper.DebugOut("Accepting message ID %v", msgid)
		}
	case string:
		helper.DebugOut("No messages to accept")
		blankParams++
	default:
		helper.Warn("Unexpected type of request.MessageIDs")
	}

	switch operatorMessageIds := request.OperatorMessageIDs.(type) {
	case []interface{}:
		helper.DebugOut("%v", operatorMessageIds)
		player.CleanUpExpiredOperatorMessages()
		for _, omsgid := range operatorMessageIds {
			helper.DebugOut("Accepting operator message ID %v", omsgid)
			present := player.AcceptOperatorMessage(int64(omsgid.(float64))) // TODO: why does Go think this is a float64 and not an int64?
			if present != nil {
				presentList = append(presentList, present.(obj.Present))
			}
		}
	case string:
		helper.DebugOut("No operator messages to accept")
		blankParams++
	default:
		helper.Warn("Unexpected type of request.OperatorMessageIDs")
	}

	if blankParams == 2 {
		//both messageIDs and operatorMessageIDs are blank; assume we're accepting all gifts
		player.CleanUpExpiredOperatorMessages()
		for _, omsgid := range player.GetAllOperatorMessageIDs() {
			helper.DebugOut("Accepting operator message ID %v", omsgid)
			present := player.AcceptOperatorMessage(omsgid)
			if present != nil {
				presentList = append(presentList, present.(obj.Present))
			}
		}
	}

	helper.DebugOut("%v", presentList)
	for _, currentPresent := range presentList {
		itemid := strconv.Itoa(int(currentPresent.ItemID))
		helper.DebugOut("Present: %s", itemid)
		helper.DebugOut("Present amount: %v", currentPresent.NumItem)
		if itemid[:2] == "12" { // ID is an item
			// check if the item is already in the player's inventory
			for _, item := range player.PlayerState.Items {
				if item.ID == itemid { // item found, increment amount
					item.Amount += currentPresent.NumItem
					break
				}
			}
		} else if itemid == strconv.Itoa(enums.ItemIDRing) { // Rings
			player.PlayerState.NumRings += currentPresent.NumItem
		} else if itemid == strconv.Itoa(enums.ItemIDRedRing) { // Red rings
			player.PlayerState.NumRedRings += currentPresent.NumItem
		} else if itemid == strconv.Itoa(enums.ItemIDEnergy) { // Revive tokens
			player.PlayerState.Energy += currentPresent.NumItem
		} else if itemid == strconv.Itoa(enums.ItemIDRaidbossEnergy) { // Boss challenge tokens
			player.EventUserRaidbossState.RaidBossEnergy += currentPresent.NumItem
		} else if itemid == strconv.Itoa(enums.IDSpecialEgg) {
			player.PlayerState.ChaoEggs += currentPresent.NumItem
			player.ChaoRouletteGroup.ChaoWheelOptions = netobj.DefaultChaoWheelOptions(player.PlayerState)
		} else if itemid == strconv.Itoa(enums.IDRouletteTicketPremium) {
			player.PlayerState.NumChaoRouletteTicket += currentPresent.NumItem
			player.ChaoRouletteGroup.ChaoWheelOptions = netobj.DefaultChaoWheelOptions(player.PlayerState)
		} else if itemid == strconv.Itoa(enums.IDRouletteTicketItem) {
			player.PlayerState.NumRouletteTicket += currentPresent.NumItem
			player.LastWheelOptions = logic.WheelRefreshLogic(player, player.LastWheelOptions)
		} else if itemid == strconv.Itoa(enums.IDRouletteTicketRaid) {
			// TODO: Add this functionality once raid boss roulette is added!
		} else if itemid == strconv.Itoa(enums.ItemIDRaidbossRing) {
			player.EventUserRaidbossState.NumRaidbossRings += currentPresent.NumItem
		} else if itemid[:2] == "40" { // ID is a Chao
			chaoIndex := player.IndexOfChao(itemid)
			if chaoIndex == -1 { // chao index not found, should never happen
				helper.InternalErr("cannot get index of chao '"+strconv.Itoa(chaoIndex)+"'", err)
				return
			}
			if player.ChaoState[chaoIndex].Status == enums.ChaoStatusNotOwned {
				// earn the Chao
				player.ChaoState[chaoIndex].Status = enums.ChaoStatusOwned
				player.ChaoState[chaoIndex].Acquired = 1
				player.ChaoState[chaoIndex].Level = 0
			}
			player.ChaoState[chaoIndex].Level += currentPresent.NumItem
			maxChaoLevel := int64(10)
			if request.Version == "1.0.0" {
				maxChaoLevel = int64(5)
			}
			if player.ChaoState[chaoIndex].Level > maxChaoLevel { // if max chao level
				player.ChaoState[chaoIndex].Level = maxChaoLevel              // reset to maximum
				player.ChaoState[chaoIndex].Status = enums.ChaoStatusMaxLevel // set status to MaxLevel
			}
		} else if itemid[:2] == "30" { // ID is a character
			charIndex := player.IndexOfChara(itemid)
			if charIndex == -1 { // character index not found, should never happen
				helper.InternalErr("cannot get index of character '"+strconv.Itoa(charIndex)+"'", err)
				return
			}
			if player.CharacterState[charIndex].Status == enums.CharacterStatusLocked {
				// unlock the character
				player.CharacterState[charIndex].Status = enums.CharacterStatusUnlocked
			} else {
				starUpCount := currentPresent.NumItem
				for starUpCount > 0 && player.CharacterState[charIndex].Star < 10 { // 10 is max amount of stars a character can have before game breaks
					starUpCount--
					player.CharacterState[charIndex].Star++
				}
			}
		} else {
			helper.Out("Unknown present ID %s", itemid)
		}
	}
	respPlayer := player
	if request.Version == "1.0.0" { // must send fewer characters
		// only get first 21 characters
		// TODO: enforce order 300000 to 300020?
		//cState = cState[:len(cState)-(len(cState)-10)]
		cState := respPlayer.CharacterState
		cState = cState[:15]
		helper.DebugOut("cState length: %v", len(cState))
		helper.DebugOut("Sent character IDs: ")
		for _, char := range cState {
			helper.DebugOut(char.ID)
		}
		respPlayer.CharacterState = cState
	}
	var response interface{}
	if baseInfo.StatusCode == status.OK {
		response = responses.GetMessage(baseInfo, respPlayer, presentList, []int64{}, player.GetAllOperatorMessageIDs())
	} else {
		response = responses.NewBaseResponse(baseInfo)
	}
	err = helper.SendCompatibleResponse(response, true)
	if err != nil {
		helper.InternalErr("Error sending response", err)
	}
	db.SavePlayer(player)
}
