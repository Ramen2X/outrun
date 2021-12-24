package muxhandlers

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/Ramen2X/outrun/analytics"
	"github.com/Ramen2X/outrun/analytics/factors"
	"github.com/Ramen2X/outrun/config"
	"github.com/Ramen2X/outrun/consts"
	"github.com/Ramen2X/outrun/db"
	"github.com/Ramen2X/outrun/emess"
	"github.com/Ramen2X/outrun/enums"
	"github.com/Ramen2X/outrun/helper"
	"github.com/Ramen2X/outrun/logic/roulette"
	"github.com/Ramen2X/outrun/netobj"
	"github.com/Ramen2X/outrun/obj"
	"github.com/Ramen2X/outrun/requests"
	"github.com/Ramen2X/outrun/responses"
	"github.com/Ramen2X/outrun/status"
)

func GetChaoWheelOptions(helper *helper.Helper) {
	player, err := helper.GetCallingPlayer(true)
	if err != nil {
		helper.InternalErr("Error getting calling player", err)
		return
	}
	baseInfo := helper.BaseInfo(emess.OK, status.OK)
	if player.Suspended {
		baseInfo.StatusCode = status.MissingPlayer
		err = helper.SendResponse(responses.NewBaseResponse(baseInfo))
		if err != nil {
			helper.InternalErr("Error sending response", err)
		}
		return
	}
	response := responses.DefaultChaoWheelOptions(baseInfo, player)
	err = helper.SendResponse(response)
	if err != nil {
		helper.InternalErr("Error sending response", err)
	}
}

func GetPrizeChaoWheelSpin(helper *helper.Helper) {
	// agnostic
	baseInfo := helper.BaseInfo(emess.OK, status.OK)
	response := responses.DefaultPrizeChaoWheel(baseInfo)
	err := helper.SendResponse(response)
	if err != nil {
		helper.InternalErr("Error sending response", err)
	}
}

func EquipChao(helper *helper.Helper) {
	recv := helper.GetGameRequest()
	var request requests.EquipChaoRequest
	err := json.Unmarshal(recv, &request)
	if err != nil {
		helper.Err("Error unmarshalling", err)
		return
	}

	player, err := helper.GetCallingPlayer(true)
	if err != nil {
		helper.InternalErr("Error getting calling player", err)
		return
	}

	baseInfo := helper.BaseInfo(emess.OK, status.OK)
	if player.Suspended {
		baseInfo.StatusCode = status.MissingPlayer
		err = helper.SendResponse(responses.NewBaseResponse(baseInfo))
		if err != nil {
			helper.InternalErr("Error sending response", err)
			return
		}
		return
	}

	mainChaoID := request.MainChaoID
	subChaoID := request.SubChaoID

	// check if the user has one chao active and is just switching
	if mainChaoID == "-1" && subChaoID == player.PlayerState.MainChaoID {
		// switching from main to sub
		player.PlayerState.MainChaoID = player.PlayerState.SubChaoID
		player.PlayerState.SubChaoID = subChaoID
		goto completed
	}
	if mainChaoID == player.PlayerState.SubChaoID && subChaoID == "-1" {
		// switching from sub to main
		player.PlayerState.SubChaoID = player.PlayerState.MainChaoID
		player.PlayerState.MainChaoID = mainChaoID
		goto completed
	}

	if mainChaoID != "-1" {
		// check if the player actually has the Chao
		chaoIndex := player.IndexOfChao(mainChaoID)
		if chaoIndex != -1 {
			chao := player.ChaoState[chaoIndex]
			if chao.Acquired != 0 && chao.Status != enums.ChaoStatusNotOwned {
				player.PlayerState.MainChaoID = mainChaoID
			} else {
				helper.Warn("Bad Chao state: chao.Acquired = %v, should = 0; chao.Status = %v, should NOT equal enums.ChaoStatusNotOwned (%v)", chao.Acquired, chao.Status, enums.ChaoStatusNotOwned)
			}
			_, err = analytics.Store(player.ID, factors.AnalyticTypeChangeMainChao)
			if err != nil {
				helper.WarnErr("Error storing analytics (AnalyticTypeChangeMainChao)", err)
			}
		} else {
			helper.Warn("Unable to find chao ID '%s'", mainChaoID)
			_, err = analytics.Store(player.ID, factors.AnalyticTypeBadRequests)
			if err != nil {
				helper.WarnErr("Error storing analytics (AnalyticTypeBadRequests)", err)
			}
		}
	}
	if subChaoID != "-1" {
		// check if the player actually has the Chao
		chaoIndex := player.IndexOfChao(subChaoID)
		if chaoIndex != -1 {
			chao := player.ChaoState[chaoIndex]
			if chao.Acquired != 0 && chao.Status != enums.ChaoStatusNotOwned {
				player.PlayerState.SubChaoID = subChaoID
			} else {
				helper.Warn("Bad Chao state: chao.Acquired = %v, should = 0; chao.Status = %v, should NOT equal enums.ChaoStatusNotOwned (%v)", chao.Acquired, chao.Status, enums.ChaoStatusNotOwned)
			}
			_, err = analytics.Store(player.ID, factors.AnalyticTypeChangeSubChao)
			if err != nil {
				helper.WarnErr("Error storing analytics (AnalyticTypeChangeSubChao)", err)
			}
		} else {
			helper.Warn("Unable to find chao ID '%s'", subChaoID)
			_, err = analytics.Store(player.ID, factors.AnalyticTypeBadRequests)
			if err != nil {
				helper.WarnErr("Error storing analytics (AnalyticTypeBadRequests)", err)
			}
		}
	}
completed:
	helper.DebugOut("Main Chao: %s", mainChaoID)
	helper.DebugOut("Sub Chao: %s", subChaoID)
	if config.CFile.Debug {
		// TODO: remove
		player.PlayerState.NumRedRings += 150
	}
	db.SavePlayer(player)

	response := responses.EquipChao(baseInfo, player.PlayerState)
	err = helper.SendResponse(response)
	if err != nil {
		helper.InternalErr("Error sending response", err)
	}
}

func CommitChaoWheelSpin(helper *helper.Helper) {
	player, err := helper.GetCallingPlayer(true)
	if err != nil {
		helper.InternalErr("Error getting calling player", err)
		return
	}

	data := helper.GetGameRequest()
	var request requests.CommitChaoWheelSpinRequest
	err = json.Unmarshal(data, &request)
	if err != nil {
		helper.InternalErr("Error unmarshalling", err)
	}

	baseInfo := helper.BaseInfo(emess.OK, status.OK)
	if player.Suspended {
		baseInfo.StatusCode = status.MissingPlayer
		err = helper.SendResponse(responses.NewBaseResponse(baseInfo))
		if err != nil {
			helper.InternalErr("Error sending response", err)
			return
		}
		return
	}

	items := player.ChaoRouletteGroup.WheelChao
	weights := player.ChaoRouletteGroup.ChaoWheelOptions.ItemWeight
	// set initial prize
	prize := netobj.CharacterIDToChaoSpinPrize("0") // This will almost certainly give the game errors if improperly counting payment!
	spinResults := []netobj.ChaoSpinResult{}

	helper.DebugOut("PRE")
	helper.DebugOut("Items: %s", items)
	helper.DebugOut("Weights: %s", items)
	helper.DebugOut("Chao Eggs (Player): %v", player.PlayerState.ChaoEggs)
	helper.DebugOut("Chao Eggs (ChaoWheelOptions): %v", player.ChaoRouletteGroup.ChaoWheelOptions.NumSpecialEgg)
	helper.DebugOut("Chao Roulette tickets (Player): %v", player.PlayerState.NumChaoRouletteTicket)
	helper.DebugOut("Chao Roulette tickets (ChaoWheelOptions): %v", player.ChaoRouletteGroup.ChaoWheelOptions.NumChaoRouletteToken)
	helper.DebugOut("Chao Roulette spin cost: %v", player.ChaoRouletteGroup.ChaoWheelOptions.SpinCost)
	helper.DebugOut("Red Rings: %v", player.PlayerState.NumRedRings)
	helper.DebugOut("Bought red rings: %v", player.PlayerState.NumBuyRedRings)
	helper.DebugOut("Spin count: %v", request.Count)

	// reset ChaoRouletteInfo if needed
	rightNow := time.Now().Unix()
	if rightNow > player.ChaoRouletteGroup.ChaoRouletteInfo.RoulettePeriodEnd { // if past period
		player.ChaoRouletteGroup.ChaoRouletteInfo = netobj.DefaultRouletteInfo() // reset all values
	}

	// spin logic
	primaryLogic := func(usingTickets bool) {
		actions := request.Count
		if request.Version == "1.1.4" {
			// 1.1.4 does not specify a value for request.Count; so we'll make it always one spin
			actions = 1
		}
		if actions < 1 {
			helper.InvalidRequest()
			return
		}
		if player.PlayerState.ChaoEggs < 10 {
			if usingTickets { // paying with ticket(s)
				player.PlayerState.NumChaoRouletteTicket -= consts.ChaoRouletteTicketCost * actions // spend ticket(s)
			} else { // paying with red ring(s)
				player.PlayerState.NumRedRings -= consts.ChaoRouletteRedRingCost * actions // spend red ring(s)
			}
		} else { //paying with chao eggs
			player.PlayerState.ChaoEggs -= 10
		}
		player.ChaoRouletteGroup.ChaoRouletteInfo.RouletteCountInPeriod++ // increment times spun in timer; TODO: Should we count request.Count?
		for actions > 0 {
			player.OptionUserResult.NumChaoRoulette++
			actions--
			gottenItemIndex, err := roulette.ChooseChaoRouletteItemIndex(items, weights) // pick a potential item index (used for later)
			if err != nil {
				helper.Err("Error choosing Chao roulette item", err)
				return
			}
			gottenItem := items[gottenItemIndex]                       // ID of prize
			gottenPrize := netobj.GenericIDToChaoSpinPrize(gottenItem) // convert ID to prize
			prize = gottenPrize
			spinResult := netobj.ChaoSpinResult{
				prize,
				[]obj.Item{},
				int64(gottenItemIndex),
			}
			//amtWon := 1
			if prize.Rarity == 100 { // Character
				// increase character level by (amount)
				charIndex := player.IndexOfChara(prize.ID)
				if charIndex == -1 { // character index not found, should never happen
					helper.InternalErr("cannot get index of character '"+strconv.Itoa(charIndex)+"'", err)
					return
				}
				if player.CharacterState[charIndex].Status == enums.CharacterStatusLocked {
					// unlock the character
					player.CharacterState[charIndex].Status = enums.CharacterStatusUnlocked
				} else {
					starUpCount := consts.ChaoRouletteCharacterStarIncrease
					for starUpCount > 0 && player.CharacterState[charIndex].Star < 10 { // 10 is max amount of stars a character can have before game breaks
						starUpCount--
						player.CharacterState[charIndex].Star++
					}
					spinResult.WonPrize.Level = player.CharacterState[charIndex].Level // set level of prize to character level
				}
			} else if prize.Rarity == 2 || prize.Rarity == 1 || prize.Rarity == 0 { // Chao
				chaoIndex := player.IndexOfChao(prize.ID)
				if chaoIndex == -1 { // chao index not found, should never happen
					helper.InternalErr("cannot get index of chao '"+strconv.Itoa(chaoIndex)+"'", err)
					return
				}
				if player.ChaoState[chaoIndex].Status == enums.ChaoStatusNotOwned {
					// earn the Chao
					player.ChaoState[chaoIndex].Status = enums.ChaoStatusOwned
					player.ChaoState[chaoIndex].Acquired = 1
					player.ChaoState[chaoIndex].Level = 0 // starting level
				} else {
					highRange := int(consts.ChaoRouletteChaoLevelIncreaseHigh)
					lowRange := int(consts.ChaoRouletteChaoLevelIncreaseLow)
					prizeChaoLevel := int64(rand.Intn(highRange-lowRange+1) + lowRange) // This level is added to the current Chao level
					//amtWon = int(prizeChaoLevel)
					maxChaoLevel := int64(10)
					if request.Version == "1.1.4" {
						maxChaoLevel = int64(5)
					}
					if player.ChaoState[chaoIndex].Level < maxChaoLevel {
						player.ChaoState[chaoIndex].Level += prizeChaoLevel
						if player.ChaoState[chaoIndex].Level > maxChaoLevel { // if max chao level (https://www.deviantart.com/vocaloidbrsfreak97/journal/So-Sonic-Runners-just-recently-updated-574789098)
							excess := player.ChaoState[chaoIndex].Level - maxChaoLevel    // get amount gone over
							prizeChaoLevel -= excess                                      // shave it from prize level
							player.ChaoState[chaoIndex].Level = maxChaoLevel              // reset to maximum
							player.ChaoState[chaoIndex].Status = enums.ChaoStatusMaxLevel // set status to MaxLevel
						}
					} else {
						player.PlayerState.ChaoEggs += 3 // maxed out; give 3 special eggs as compensation
						spinResult.ItemList = append(spinResult.ItemList, obj.NewItem(strconv.Itoa(enums.IDSpecialEgg), 3))
					}
					spinResult.WonPrize.Level = player.ChaoState[chaoIndex].Level
				}
			} else { // Should never happen!
				helper.InternalErr("unknown prize rarity '"+strconv.Itoa(int(prize.Rarity))+"'", fmt.Errorf("")) // TODO: Probably shouldn't use a blank error?
			}
			/*index := 0
			for index < len(items) {
				spinResult.ItemList = append(spinResult.ItemList, obj.NewItem(items[index], int64(amtWon)))
				index++
			}*/
			spinResults = append(spinResults, spinResult) // add spin result to results list (See spinResults declaration)
		}
		// create a new wheel; must be done after ALL player operations are done
		chaoCanBeLevelled := !player.AllChaoMaxLevel(request.Version == "1.1.4")
		charactersCanBeLevelled := !player.AllCharactersMaxLevel()
		helper.DebugOut("Chao can be levelled: %v", chaoCanBeLevelled)
		helper.DebugOut("Characters can be levelled: %v", charactersCanBeLevelled)
		fixRarities := func(rarities []int64) ([]int64, bool) {
			newRarities := []int64{}
			if !chaoCanBeLevelled && !charactersCanBeLevelled {
				// Wow, they can't upgrade _anything!_
				return newRarities, false
			}
			if config.CFile.Debug {
				player.PlayerState.NumRedRings += 150
				//return []int64{100, 100, 100, 100, 100, 100, 100, 100}, true
				return []int64{0, 0, 0, 0, 0, 0, 0, 0}, true
			}
			for _, r := range rarities {
				if r == 0 || r == 1 || r == 2 { // Chao
					if chaoCanBeLevelled {
						newRarities = append(newRarities, r)
					} else {
						newRarities = append(newRarities, 100) // append a character
					}
				} else if r == 100 { // character
					if charactersCanBeLevelled {
						newRarities = append(newRarities, r)
					} else {
						newRarities = append(newRarities, int64(rand.Intn(3))) // append random rarity Chao
					}
				} else { // should never happen
					panic(fmt.Errorf("invalid rarity '" + strconv.Itoa(int(r)) + "'")) // TODO: use better way to handle
				}
			}
			return newRarities, true
		}
		player.ChaoRouletteGroup.ChaoWheelOptions = netobj.DefaultChaoWheelOptions(player.PlayerState) // create a new wheel
		newRarities, ok := fixRarities(player.ChaoRouletteGroup.ChaoWheelOptions.Rarity)
		if !ok { // if player is entirely unable to upgrade anything
			// TODO: this is probably not the right way to do this!
			player.ChaoRouletteGroup.ChaoWheelOptions.SpinCost = player.PlayerState.NumChaoRouletteTicket + player.PlayerState.NumRedRings // make it impossible for player to use roulette
		} else { // if player can upgrade
			player.ChaoRouletteGroup.ChaoWheelOptions.Rarity = newRarities
		}
		//newItems, err := roulette.GetRandomChaoRouletteItems(player.ChaoRouletteGroup.ChaoWheelOptions.Rarity, player.GetAllMaxLevelIDs()) // create new wheel items
		//newItems, err := roulette.GetRandomChaoRouletteItems(player.ChaoRouletteGroup.ChaoWheelOptions.Rarity, player.GetAllNonMaxedChaoAndCharacters()) // create new wheel items
		newItems, newRarities, err := roulette.GetRandomChaoRouletteItems(player.ChaoRouletteGroup.ChaoWheelOptions.Rarity, player.GetAllNonMaxedCharacters(), player.GetAllNonMaxedChao(request.Version == "1.1.4"), request.Version == "1.1.4")
		if err != nil {
			helper.InternalErr("Error getting new items", err)
			return
		}
		player.ChaoRouletteGroup.WheelChao = newItems
		player.ChaoRouletteGroup.ChaoWheelOptions.Rarity = newRarities
		helper.DebugOut("Rarities: %v", newRarities)
		if config.CFile.Debug {
			player.ChaoRouletteGroup.WheelChao = []string{enums.CTStrTails, enums.CTStrTails, enums.CTStrTails, enums.CTStrTails, enums.CTStrTails, enums.CTStrTails, enums.CTStrTails, enums.CTStrTails}
		}
	}

	hasTickets := player.PlayerState.NumChaoRouletteTicket >= consts.ChaoRouletteTicketCost*request.Count
	hasAvailableRings := player.PlayerState.NumRedRings >= consts.ChaoRouletteRedRingCost*request.Count

	if hasTickets || player.PlayerState.ChaoEggs >= 10 { // if tickets or chao eggs to spend
		primaryLogic(true)
	} else if hasAvailableRings { // if no tickets, but sufficient red rings
		primaryLogic(false)
	} else { // no tickets nor sufficient red rings
		baseInfo.StatusCode = status.RouletteUseLimit
	}

	helper.DebugOut("POST")
	helper.DebugOut("Items: %s", items)
	helper.DebugOut("Weights: %s", items)
	helper.DebugOut("Chao Eggs (Player): %v", player.PlayerState.ChaoEggs)
	helper.DebugOut("Chao Eggs (ChaoWheelOptions): %v", player.ChaoRouletteGroup.ChaoWheelOptions.NumSpecialEgg)
	helper.DebugOut("Chao Roulette tickets (Player): %v", player.PlayerState.NumChaoRouletteTicket)
	helper.DebugOut("Chao Roulette tickets (ChaoWheelOptions): %v", player.ChaoRouletteGroup.ChaoWheelOptions.NumChaoRouletteToken)
	helper.DebugOut("Chao Roulette spin cost: %v", player.ChaoRouletteGroup.ChaoWheelOptions.SpinCost)

	cState := player.CharacterState
	if request.Version == "1.1.4" { // must send fewer characters
		// only get first 21 characters
		// TODO: enforce order 300000 to 300020?
		//cState = cState[:len(cState)-(len(cState)-10)]
		cState = cState[:16]
		helper.DebugOut("cState length: " + strconv.Itoa(len(cState)))
		helper.DebugOut("Sent character IDs: ")
		for _, char := range cState {
			helper.DebugOut(char.ID)
		}
	}
	response := responses.ChaoWheelSpin(baseInfo, player.PlayerState, cState, player.ChaoState, player.ChaoRouletteGroup.ChaoWheelOptions, spinResults)

	err = db.SavePlayer(player)
	if err != nil {
		helper.InternalErr("Error saving player", err)
		return
	}

	err = helper.SendResponse(response)
	if err != nil {
		helper.InternalErr("Error sending response", err)
		return
	}
	_, err = analytics.Store(player.ID, factors.AnalyticTypeSpinChaoRoulette)
	if err != nil {
		helper.WarnErr("Error storing analytics (AnalyticTypeSpinChaoRoulette)", err)
	}
}
