package muxhandlers

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/Ramen2X/outrun/analytics"
	"github.com/Ramen2X/outrun/analytics/factors"
	"github.com/Ramen2X/outrun/config/campaignconf"
	"github.com/Ramen2X/outrun/consts"
	"github.com/Ramen2X/outrun/db"
	"github.com/Ramen2X/outrun/emess"
	"github.com/Ramen2X/outrun/enums"
	"github.com/Ramen2X/outrun/helper"
	"github.com/Ramen2X/outrun/logic"
	"github.com/Ramen2X/outrun/logic/conversion"
	"github.com/Ramen2X/outrun/netobj"
	"github.com/Ramen2X/outrun/obj"
	"github.com/Ramen2X/outrun/requests"
	"github.com/Ramen2X/outrun/responses"
	"github.com/Ramen2X/outrun/status"
)

func GetWheelOptions(helper *helper.Helper) {
	player, err := helper.GetCallingPlayer(true)
	if err != nil {
		helper.InternalErr("Error getting calling player", err)
		return
	}
	baseInfo := helper.BaseInfo(emess.OK, status.OK)
	if player.Suspended {
		baseInfo.StatusCode = status.MissingPlayer
		err = helper.SendCompatibleResponse(responses.NewBaseResponse(baseInfo), true)
		if err != nil {
			helper.InternalErr("Error sending response", err)
			return
		}
		return
	}

	//player.LastWheelOptions = netobj.DefaultWheelOptions(player.PlayerState) // generate new wheel for 'reroll' mechanic
	helper.DebugOut("Time now: %v", time.Now().Unix())
	helper.DebugOut("RoulettePeriodEnd: %v", player.RouletteInfo.RoulettePeriodEnd)
	// check if we need to reset the end period
	endPeriod := player.RouletteInfo.RoulettePeriodEnd
	if time.Now().Unix() > endPeriod {
		player.RouletteInfo = netobj.DefaultRouletteInfo() // Effectively reset everything, set new end time
	}

	// refresh wheel
	player.LastWheelOptions = logic.WheelRefreshLogic(player, player.LastWheelOptions)

	response := responses.WheelOptions(baseInfo, player.LastWheelOptions)
	err = helper.SendCompatibleResponse(response, true)
	if err != nil {
		helper.InternalErr("Error sending response", err)
	}
}

func CommitWheelSpin(helper *helper.Helper) {
	recv := helper.GetGameRequest()
	var request requests.CommitWheelSpinRequest
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
		err = helper.SendCompatibleResponse(responses.NewBaseResponse(baseInfo), true)
		if err != nil {
			helper.InternalErr("Error sending response", err)
			return
		}
		return
	}
	helper.DebugOut("request.Count: %v", request.Count)

	freeSpins := consts.RouletteFreeSpins
	campaignList := []obj.Campaign{}
	if campaignconf.CFile.AllowCampaigns {
		for _, confCampaign := range campaignconf.CFile.CurrentCampaigns {
			newCampaign := conversion.ConfiguredCampaignToCampaign(confCampaign)
			campaignList = append(campaignList, newCampaign)
		}
	}
	index := 0
	for index < len(campaignList) {
		if obj.IsCampaignActive(campaignList[index]) && campaignList[index].Type == enums.CampaignTypeFreeWheelSpinCount {
			freeSpins = campaignList[index].Content
		}
		index++
	}

	endPeriod := player.RouletteInfo.RoulettePeriodEnd
	helper.DebugOut("Time now: %v", time.Now().Unix())
	helper.DebugOut("End period: %v", endPeriod)
	if time.Now().Unix() > endPeriod {
		player.RouletteInfo = netobj.DefaultRouletteInfo() // Effectively reset everything, set new end time
		helper.DebugOut("New roulette period")
		helper.DebugOut("RouletteCountInPeriod: %v", player.RouletteInfo.RouletteCountInPeriod)
	}

	hasTickets := player.PlayerState.NumRouletteTicket > 0
	hasFreeSpins := player.RouletteInfo.RouletteCountInPeriod < freeSpins
	helper.DebugOut("Has tickets: %v", hasTickets)
	helper.DebugOut("Number of tickets: %v", player.PlayerState.NumRouletteTicket)
	helper.DebugOut("Has free spins: %v", hasFreeSpins)
	helper.DebugOut("Roulette count: %v", player.RouletteInfo.RouletteCountInPeriod)
	landedOnUpgrade := false
	if hasTickets || hasFreeSpins {
		//if player.LastWheelOptions.NumRemainingRoulette > 0 {
		wonItem := player.LastWheelOptions.Items[player.LastWheelOptions.ItemWon]
		itemExists := player.IndexOfItem(wonItem) != -1
		if itemExists {
			amountOfItemWon := player.LastWheelOptions.Item[player.LastWheelOptions.ItemWon]
			helper.DebugOut("wonItem: %v", wonItem)
			helper.DebugOut("amountOfItemWon: %v", amountOfItemWon)
			itemIndex := player.IndexOfItem(wonItem)
			helper.DebugOut("Amount of item player has: %v", player.PlayerState.Items[itemIndex].Amount)
			player.PlayerState.Items[itemIndex].Amount += amountOfItemWon
			helper.DebugOut("New amount of item player has: %v", player.PlayerState.Items[itemIndex].Amount)
		} else {
			if wonItem == strconv.Itoa(enums.IDTypeItemRouletteWin) {
				// BIG/SUPER/Jackpot
				if player.LastWheelOptions.RouletteRank == enums.WheelRankSuper { // Don't award jackpot unless on super
					player.PlayerState.NumRings += player.LastWheelOptions.NumJackpotRing
				} else {
					landedOnUpgrade = true
				}
			} else if wonItem == strconv.Itoa(enums.IDTypeRedRing) {
				// Red rings
				player.PlayerState.NumRedRings += player.LastWheelOptions.Item[player.LastWheelOptions.ItemWon]
			} else if wonItem[:2] == "40" {
				// Chao
				amountOfItemWon := player.LastWheelOptions.Item[player.LastWheelOptions.ItemWon]
				helper.DebugOut("wonItem: %v", wonItem)
				helper.DebugOut("amountOfItemWon: %v", amountOfItemWon)

				chaoIndex := player.IndexOfChao(wonItem)
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
				player.ChaoState[chaoIndex].Level += amountOfItemWon
				maxChaoLevel := int64(10)
				if request.Version == "1.0.0" {
					maxChaoLevel = int64(5)
				}
				if player.ChaoState[chaoIndex].Level > maxChaoLevel { // if max chao level
					player.ChaoState[chaoIndex].Level = maxChaoLevel              // reset to maximum
					player.ChaoState[chaoIndex].Status = enums.ChaoStatusMaxLevel // set status to MaxLevel
				}
			} else {
				helper.Warn("item '" + wonItem + "' not found")
				helper.DebugOut(wonItem[:2])
			}
		}

		helper.DebugOut("Time now: %v", time.Now().Unix())
		helper.DebugOut("RoulettePeriodEnd: %v", player.RouletteInfo.RoulettePeriodEnd)
		endPeriod := player.RouletteInfo.RoulettePeriodEnd
		helper.DebugOut("Time now (passed): %v", time.Now().Unix())
		helper.DebugOut("End period (passed): %v", endPeriod)
		if time.Now().Unix() > endPeriod { // TODO: Do we still need this?
			player.RouletteInfo = netobj.DefaultRouletteInfo() // Effectively reset everything, set new end time
			helper.DebugOut("New roulette period")
			helper.DebugOut("RouletteCountInPeriod: %v", player.RouletteInfo.RouletteCountInPeriod)
		}

		// generate NEXT! wheel
		if !landedOnUpgrade {
			//don't use up spin if we landed on an upgrade
			player.RouletteInfo.RouletteCountInPeriod++ // we've spun an additional time
			if player.RouletteInfo.RouletteCountInPeriod > freeSpins {
				// we've run out of free spins for the period
				player.PlayerState.NumRouletteTicket--
			}
		}
		numRouletteTicket := player.PlayerState.NumRouletteTicket
		player.OptionUserResult.NumItemRoulette++
		rouletteCount := player.RouletteInfo.RouletteCountInPeriod // get amount of times we've spun the wheel today
		//player.LastWheelOptions = netobj.DefaultWheelOptions(numRouletteTicket, rouletteCount) // create wheel
		oldRanking := player.LastWheelOptions.RouletteRank
		player.LastWheelOptions = netobj.UpgradeWheelOptions(player.LastWheelOptions, numRouletteTicket, rouletteCount, freeSpins) // create wheel
		if player.RouletteInfo.GotJackpotThisPeriod {
			player.LastWheelOptions.NumJackpotRing = 1
		}
		if wonItem == strconv.Itoa(enums.IDTypeItemRouletteWin) && oldRanking == enums.WheelRankSuper { // won jackpot in super wheel
			helper.DebugOut("Won jackpot in super wheel")
			player.RouletteInfo.GotJackpotThisPeriod = true
			player.OptionUserResult.NumJackpot++
			if player.LastWheelOptions.NumJackpotRing > player.OptionUserResult.NumMaximumJackpotRings {
				player.OptionUserResult.NumMaximumJackpotRings = player.LastWheelOptions.NumJackpotRing
			}
		}
	} else {
		// do not modify the wheel, set error status
		baseInfo.StatusCode = status.RouletteUseLimit
	}

	cState := player.CharacterState
	if request.Version == "1.0.0" { // must send fewer characters
		// only get first 21 characters
		// TODO: enforce order 300000 to 300020?
		//cState = cState[:len(cState)-(len(cState)-10)]
		cState = cState[:15]
		helper.DebugOut("cState length: " + strconv.Itoa(len(cState)))
		helper.DebugOut("Sent character IDs: ")
		for _, char := range cState {
			helper.DebugOut(char.ID)
		}
	}
	response := responses.WheelSpin(baseInfo, player.PlayerState, cState, player.ChaoState, player.LastWheelOptions)

	err = db.SavePlayer(player)
	if err != nil {
		helper.InternalErr("Error saving player", err)
		return
	}

	err = helper.SendCompatibleResponse(response, true)
	if err != nil {
		helper.InternalErr("Error sending response", err)
		return
	}
	_, err = analytics.Store(player.ID, factors.AnalyticTypeSpinItemRoulette)
	if err != nil {
		helper.WarnErr("Error storing analytics (AnalyticTypeSpinItemRoulette)", err)
	}
}

// 1.1.4 support
func GetWheelSpinInfo(helper *helper.Helper) {
	baseInfo := helper.BaseInfo(emess.OK, status.OK)
	response := responses.DefaultWheelSpinInfo(baseInfo)
	err := helper.SendCompatibleResponse(response, true)
	if err != nil {
		helper.InternalErr("Error sending response", err)
	}
}
