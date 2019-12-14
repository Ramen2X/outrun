package muxhandlers

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/jinzhu/now"

	"github.com/fluofoxxo/outrun/analytics"
	"github.com/fluofoxxo/outrun/analytics/factors"
	"github.com/fluofoxxo/outrun/config"
	"github.com/fluofoxxo/outrun/config/campaignconf"
	"github.com/fluofoxxo/outrun/config/gameconf"
	"github.com/fluofoxxo/outrun/consts"
	"github.com/fluofoxxo/outrun/db"
	"github.com/fluofoxxo/outrun/emess"
	"github.com/fluofoxxo/outrun/enums"
	"github.com/fluofoxxo/outrun/helper"
	"github.com/fluofoxxo/outrun/logic/campaign"
	"github.com/fluofoxxo/outrun/logic/conversion"
	"github.com/fluofoxxo/outrun/logic/gameplay"
	"github.com/fluofoxxo/outrun/netobj"
	"github.com/fluofoxxo/outrun/obj"
	"github.com/fluofoxxo/outrun/obj/constobjs"
	"github.com/fluofoxxo/outrun/requests"
	"github.com/fluofoxxo/outrun/responses"
	"github.com/fluofoxxo/outrun/status"
)

func GetDailyChallengeData(helper *helper.Helper) {
	player, err := helper.GetCallingPlayer()
	if err != nil {
		helper.InternalErr("Error getting player", err)
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
	response := responses.DailyChallengeData(baseInfo)
	err = helper.SendResponse(response)
	if err != nil {
		helper.InternalErr("Error sending response", err)
	}
}

func GetCostList(helper *helper.Helper) {
	// no player, agonstic
	baseInfo := helper.BaseInfo(emess.OK, status.OK)
	response := responses.DefaultCostList(baseInfo)
	err := helper.SendResponse(response)
	if err != nil {
		helper.InternalErr("Error sending response", err)
	}
}

func GetMileageData(helper *helper.Helper) {
	player, err := helper.GetCallingPlayer()
	if err != nil {
		helper.InternalErr("Error getting player", err)
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
	response := responses.DefaultMileageData(baseInfo, player)
	err = helper.SendResponse(response)
	if err != nil {
		helper.InternalErr("Error sending response", err)
	}
}

func GetCampaignList(helper *helper.Helper) {
	campaignList := []obj.Campaign{}
	if campaignconf.CFile.AllowCampaigns {
		for _, confCampaign := range campaignconf.CFile.CurrentCampaigns {
			newCampaign := conversion.ConfiguredCampaignToCampaign(confCampaign)
			campaignList = append(campaignList, newCampaign)
		}
	}
	helper.DebugOut("Campaign list: %v", campaignList)
	baseInfo := helper.BaseInfo(emess.OK, status.OK)
	response := responses.CampaignList(baseInfo, campaignList)
	err := helper.SendResponse(response)
	if err != nil {
		helper.InternalErr("Error sending response", err)
	}
}

func QuickActStart(helper *helper.Helper) {
	recv := helper.GetGameRequest()
	var request requests.QuickActStartRequest
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
	if player.Suspended {
		baseInfo.StatusCode = status.MissingPlayer
		err = helper.SendResponse(responses.NewBaseResponse(baseInfo))
		if err != nil {
			helper.InternalErr("Error sending response", err)
			return
		}
		return
	}
	campaignList := []obj.Campaign{}
	if campaignconf.CFile.AllowCampaigns {
		for _, confCampaign := range campaignconf.CFile.CurrentCampaigns {
			newCampaign := conversion.ConfiguredCampaignToCampaign(confCampaign)
			campaignList = append(campaignList, newCampaign)
		}
	}
	helper.DebugOut("Campaign list: %v", campaignList)
	// consume items
	modToStringSlice := func(ns []int64) []string {
		result := []string{}
		for _, n := range ns {
			result = append(result, fmt.Sprintf("%v", n))
		}
		return result
	}
	for time.Now().UTC().Unix() >= player.PlayerState.EnergyRenewsAt && player.PlayerState.Energy < player.PlayerVarious.EnergyRecoveryMax {
		player.PlayerState.Energy++
		player.PlayerState.EnergyRenewsAt += player.PlayerVarious.EnergyRecoveryTime
	}
	if player.PlayerState.Energy+player.PlayerState.EnergyBuy > 0 {
		if gameconf.CFile.EnableEnergyConsumption {
			if player.PlayerState.EnergyBuy > 0 {
				player.PlayerState.EnergyBuy--
			} else {
				if player.PlayerState.Energy >= player.PlayerVarious.EnergyRecoveryMax {
					player.PlayerState.EnergyRenewsAt = time.Now().UTC().Unix() + player.PlayerVarious.EnergyRecoveryTime
				}
				player.PlayerState.Energy--
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
						baseInfo.StatusCode = status.NotEnoughRings
						break
					}
					player.PlayerState.NumRings -= consumedRings
				}
			}
		}
	} else {
		baseInfo.StatusCode = status.NotEnoughEnergy
	}
	helper.DebugOut(fmt.Sprintf("%v", player.PlayerState.Items))
	response := responses.DefaultQuickActStart(baseInfo, player, campaignList)
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
	_, err = analytics.Store(player.ID, factors.AnalyticTypeTimedStarts)
	if err != nil {
		helper.WarnErr("Error storing analytics (AnalyticTypeTimedStarts)", err)
	}
}

func ActStart(helper *helper.Helper) {
	recv := helper.GetGameRequest()
	var request requests.ActStartRequest
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
	if player.Suspended {
		baseInfo.StatusCode = status.MissingPlayer
		err = helper.SendResponse(responses.NewBaseResponse(baseInfo))
		if err != nil {
			helper.InternalErr("Error sending response", err)
			return
		}
		return
	}
	campaignList := []obj.Campaign{}
	if campaignconf.CFile.AllowCampaigns {
		for _, confCampaign := range campaignconf.CFile.CurrentCampaigns {
			newCampaign := conversion.ConfiguredCampaignToCampaign(confCampaign)
			campaignList = append(campaignList, newCampaign)
		}
	}
	helper.DebugOut("Campaign list: %v", campaignList)
	helper.DebugOut(fmt.Sprintf("%v", player.PlayerState.Items))
	// consume items
	modToStringSlice := func(ns []int64) []string {
		result := []string{}
		for _, n := range ns {
			result = append(result, fmt.Sprintf("%v", n))
		}
		return result
	}
	helper.DebugOut(fmt.Sprintf("%v", player.PlayerState.Items))
	for time.Now().UTC().Unix() >= player.PlayerState.EnergyRenewsAt && player.PlayerState.Energy < player.PlayerVarious.EnergyRecoveryMax {
		player.PlayerState.Energy++
		player.PlayerState.EnergyRenewsAt += player.PlayerVarious.EnergyRecoveryTime
	}
	if player.PlayerState.Energy+player.PlayerState.EnergyBuy > 0 {
		if gameconf.CFile.EnableEnergyConsumption {
			if player.PlayerState.EnergyBuy > 0 {
				player.PlayerState.EnergyBuy--
			} else {
				player.PlayerState.Energy--
				if player.PlayerState.Energy < player.PlayerVarious.EnergyRecoveryMax {
					player.PlayerState.EnergyRenewsAt = time.Now().UTC().Unix() + player.PlayerVarious.EnergyRecoveryTime
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
						baseInfo.StatusCode = status.NotEnoughRings
						break
					}
					player.PlayerState.NumRings -= consumedRings
				}
			}
		}
	} else {
		baseInfo.StatusCode = status.NotEnoughEnergy
	}
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
	response := responses.DefaultActStart(baseInfo, respPlayer, campaignList)
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
	_, err = analytics.Store(player.ID, factors.AnalyticTypeStoryStarts)
	if err != nil {
		helper.WarnErr("Error storing analytics (AnalyticTypeStoryStarts)", err)
	}
}

func ActRetry(helper *helper.Helper) {
	player, err := helper.GetCallingPlayer()
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
	if player.PlayerState.NumRedRings >= 5 { //does the player actually have enough red rings to be revived?
		player.PlayerState.NumRedRings -= 5
		err = db.SavePlayer(player)
		if err != nil {
			helper.InternalErr("Error saving player", err)
			return
		}
	} else {
		baseInfo.StatusCode = status.NotEnoughRedRings
	}
	response := responses.NewBaseResponse(baseInfo)
	err = helper.SendResponse(response)
	if err != nil {
		helper.InternalErr("Error sending response", err)
		return
	}
	_, err = analytics.Store(player.ID, factors.AnalyticTypeRevives)
	if err != nil {
		helper.WarnErr("Error storing analytics (AnalyticTypeRevives)", err)
	}
}

func ActRetryFree(helper *helper.Helper) {
	player, err := helper.GetCallingPlayer()
	if err != nil {
		helper.InternalErr("Error getting calling player", err)
		return
	}
	// more than likely used for ad-based revives
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
	response := responses.NewBaseResponse(baseInfo)
	err = helper.SendResponse(response)
	if err != nil {
		helper.InternalErr("Error sending response", err)
		return
	}
	_, err = analytics.Store(player.ID, factors.AnalyticTypeRevives)
	if err != nil {
		helper.WarnErr("Error storing analytics (AnalyticTypeRevives)", err)
	}
}

func QuickPostGameResults(helper *helper.Helper) {
	recv := helper.GetGameRequest()
	var request requests.QuickPostGameResultsRequest
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
	if player.Suspended {
		baseInfo.StatusCode = status.MissingPlayer
		err = helper.SendResponse(responses.NewBaseResponse(baseInfo))
		if err != nil {
			helper.InternalErr("Error sending response", err)
			return
		}
		return
	}

	//update energy counter
	for time.Now().UTC().Unix() >= player.PlayerState.EnergyRenewsAt && player.PlayerState.Energy < player.PlayerVarious.EnergyRecoveryMax {
		player.PlayerState.Energy++
		player.PlayerState.EnergyRenewsAt += player.PlayerVarious.EnergyRecoveryTime
	}

	hasSubCharacter := player.PlayerState.SubCharaID != "-1"
	var subC netobj.Character
	mainC, err := player.GetMainChara()
	if err != nil {
		helper.InternalErr("Error getting main character", err)
		return
	}
	playCharacters := []netobj.Character{ // assume only main character active right now
		mainC,
	}
	if hasSubCharacter {
		subC, err = player.GetSubChara()
		if err != nil {
			helper.InternalErr("Error getting sub character", err)
			return
		}
		playCharacters = []netobj.Character{ // add sub character to playCharacters
			mainC,
			subC,
		}
	}
	if request.Closed == 0 { // If the game wasn't exited out of
		player.PlayerState.NumRings += request.Rings
		player.PlayerState.NumRedRings += request.RedRings
		player.PlayerState.NumRouletteTicket += request.RedRings // TODO: URGENT! Remove as soon as possible!
		player.PlayerState.Animals += request.Animals
		player.OptionUserResult.NumTakeAllRings += request.Rings
		player.OptionUserResult.NumTakeAllRedRings += request.RedRings
		playerTimedHighScore := player.PlayerState.TimedHighScore
		if request.Score > playerTimedHighScore {
			player.PlayerState.TimedHighScore = request.Score
		}
		if time.Now().UTC().Unix() > player.PlayerState.WeeklyScoresExpireAt {
			if player.PlayerState.TotalScoreThisPeriod > player.PlayerState.TotalScore {
				player.PlayerState.TotalScore = player.PlayerState.TotalScoreThisPeriod
			}
			if player.PlayerState.TimedTotalScoreThisPeriod > player.PlayerState.TimedTotalScore {
				player.PlayerState.TimedTotalScore = player.PlayerState.TimedTotalScoreThisPeriod
			}
			player.PlayerState.HighScoreThisPeriod = 0
			player.PlayerState.TimedHighScoreThisPeriod = 0
			player.PlayerState.TotalScoreThisPeriod = 0
			player.PlayerState.TimedTotalScoreThisPeriod = 0
			player.PlayerState.WeeklyScoresExpireAt = now.EndOfWeek().UTC().Unix()
		}
		playerTimedHighScoreThisPeriod := player.PlayerState.TimedHighScoreThisPeriod
		if request.Score > playerTimedHighScoreThisPeriod {
			player.PlayerState.TimedHighScoreThisPeriod = request.Score
		}
		player.PlayerState.TimedTotalScoreThisPeriod += request.Score
		if player.PlayerState.TimedTotalScoreThisPeriod > player.OptionUserResult.QuickTotalSumHighScore {
			player.OptionUserResult.QuickTotalSumHighScore = player.PlayerState.TimedTotalScoreThisPeriod
		}
		//player.PlayerState.TotalDistance += request.Distance  // We don't do this in timed mode!
		// increase character(s)'s experience
		expIncrease := request.Rings + request.FailureRings // all rings collected
		abilityIndex := 1
		for abilityIndex == 1 { // unused ability is at index 1
			abilityIndex = rand.Intn(len(mainC.AbilityLevel))
		}
		// check that increases exist
		_, ok := consts.UpgradeIncreases[mainC.ID]
		if !ok {
			helper.InternalErr("Error getting upgrade increase", fmt.Errorf("no key '%s' in consts.UpgradeIncreases", mainC.ID))
			return
		}
		if hasSubCharacter {
			_, ok = consts.UpgradeIncreases[subC.ID]
			if !ok {
				helper.InternalErr("Error getting upgrade increase for sub character", fmt.Errorf("no key '%s' in consts.UpgradeIncreases", subC.ID))
				return
			}
		}

		playCharacters[0].AbilityLevelUp = []int64{}
		playCharacters[0].AbilityLevelUpExp = []int64{}
		if playCharacters[0].Level < 100 {
			playCharacters[0].Exp += expIncrease
			for playCharacters[0].Exp >= playCharacters[0].Cost {
				// more exp than cost = level up
				if playCharacters[0].Level < 100 {
					abilityIndex = 1
					for abilityIndex == 1 || playCharacters[0].AbilityLevel[abilityIndex] >= 10 { // unused ability is at index 1
						abilityIndex = rand.Intn(len(mainC.AbilityLevel))
					}
					playCharacters[0].Level++                                               // increase level
					playCharacters[0].AbilityLevel[abilityIndex]++                          // increase ability level
					playCharacters[0].Exp -= playCharacters[0].Cost                         // remove cost from exp
					playCharacters[0].Cost += consts.UpgradeIncreases[playCharacters[0].ID] // increase cost
					playCharacters[0].AbilityLevelUp = append(playCharacters[0].AbilityLevelUp, int64(abilityIndex))
					playCharacters[0].AbilityLevelUpExp = append(playCharacters[0].AbilityLevelUpExp, playCharacters[0].Cost)
				} else {
					playCharacters[0].Exp -= playCharacters[0].Cost
				}
			}
		}

		if hasSubCharacter {
			playCharacters[1].AbilityLevelUp = []int64{}
			playCharacters[1].AbilityLevelUpExp = []int64{}
			if playCharacters[1].Level < 100 {
				playCharacters[1].Exp += expIncrease
				for playCharacters[1].Exp >= playCharacters[1].Cost {
					// more exp than cost = level up
					if playCharacters[1].Level < 100 {
						abilityIndex = 1
						for abilityIndex == 1 || playCharacters[1].AbilityLevel[abilityIndex] >= 10 { // unused ability is at index 1
							abilityIndex = rand.Intn(len(playCharacters[1].AbilityLevel))
						}
						playCharacters[1].Level++                                               // increase level
						playCharacters[1].AbilityLevel[abilityIndex]++                          // increase ability level
						playCharacters[1].Exp -= playCharacters[1].Cost                         // remove cost from exp
						playCharacters[1].Cost += consts.UpgradeIncreases[playCharacters[1].ID] // increase cost
						playCharacters[1].AbilityLevelUp = append(playCharacters[1].AbilityLevelUp, int64(abilityIndex))
						playCharacters[1].AbilityLevelUpExp = append(playCharacters[1].AbilityLevelUpExp, playCharacters[1].Cost)
					} else {
						playCharacters[1].Exp -= playCharacters[1].Cost
					}
				}
			}
		}

		helper.DebugOut("Old mainC Exp: %v / %v", mainC.Exp, mainC.Cost)
		helper.DebugOut("Old mainC Level: %v", mainC.Level)
		if hasSubCharacter {
			helper.DebugOut("Old subC Exp: %v / %v", subC.Exp, subC.Cost)
			helper.DebugOut("Old subC Level: %v", subC.Level)
		}
		helper.DebugOut("New mainC Exp: %v / %v", playCharacters[0].Exp, playCharacters[0].Cost)
		helper.DebugOut("New mainC Level: %v", playCharacters[0].Level)
		helper.DebugOut("mainC AbilityLevelUp: %v", playCharacters[0].AbilityLevelUp)
		helper.DebugOut("mainC AbilityLevelUpExp: %v", playCharacters[0].AbilityLevelUpExp)
		if hasSubCharacter {
			helper.DebugOut("New subC Exp: %v / %v", playCharacters[1].Exp, playCharacters[1].Cost)
			helper.DebugOut("New subC Level: %v", playCharacters[1].Level)
			helper.DebugOut("subC AbilityLevelUp: %v", playCharacters[1].AbilityLevelUp)
			helper.DebugOut("subC AbilityLevelUpExp: %v", playCharacters[1].AbilityLevelUpExp)
		}

		/*playCharacters = []netobj.Character{ // TODO: check if this redefinition is needed
			mainC,
			subC,
		}*/
	}

	mainCIndex := player.IndexOfChara(mainC.ID) // TODO: check if -1
	subCIndex := -1
	if hasSubCharacter {
		subCIndex = player.IndexOfChara(subC.ID) // TODO: check if -1
	}

	response := responses.DefaultQuickPostGameResults(baseInfo, player, playCharacters)
	// apply the save after the response so that we don't break the leveling
	mainC = playCharacters[0]
	if hasSubCharacter {
		subC = playCharacters[1]
	}
	player.CharacterState[mainCIndex] = mainC
	if hasSubCharacter {
		player.CharacterState[subCIndex] = subC
	}
	helper.DebugOut("CheatResult: " + request.CheatResult)
	err = db.SavePlayer(player)
	if err != nil {
		helper.InternalErr("Error saving player", err)
		return
	}
	helper.DebugOut(fmt.Sprintf("%v", player.PlayerState.Items))

	err = helper.SendResponse(response)
	if err != nil {
		helper.InternalErr("Error sending response", err)
		return
	}
	_, err = analytics.Store(player.ID, factors.AnalyticTypeTimedEnds)
	if err != nil {
		helper.WarnErr("Error storing analytics (AnalyticTypeTimedEnds)", err)
	}
}

func PostGameResults(helper *helper.Helper) {
	recv := helper.GetGameRequest()
	var request requests.PostGameResultsRequest
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
	if player.Suspended {
		baseInfo.StatusCode = status.MissingPlayer
		err = helper.SendResponse(responses.NewBaseResponse(baseInfo))
		if err != nil {
			helper.InternalErr("Error sending response", err)
			return
		}
		return
	}

	//update energy counter
	for time.Now().UTC().Unix() >= player.PlayerState.EnergyRenewsAt && player.PlayerState.Energy < player.PlayerVarious.EnergyRecoveryMax {
		player.PlayerState.Energy++
		player.PlayerState.EnergyRenewsAt += player.PlayerVarious.EnergyRecoveryTime
	}

	hasSubCharacter := player.PlayerState.SubCharaID != "-1"
	var subC netobj.Character
	mainC, err := player.GetMainChara()
	if err != nil {
		helper.InternalErr("Error getting main character", err)
		return
	}
	playCharacters := []netobj.Character{ // assume only main character active right now
		mainC,
	}
	if hasSubCharacter {
		subC, err = player.GetSubChara()
		if err != nil {
			helper.InternalErr("Error getting sub character", err)
			return
		}
		playCharacters = []netobj.Character{ // add sub character to playCharacters
			mainC,
			subC,
		}
	}
	helper.DebugOut("Pre-function")
	helper.DebugOut("Chapter: %v", player.MileageMapState.Chapter)
	helper.DebugOut("Episode: %v", player.MileageMapState.Episode)
	helper.DebugOut("StageTotalScore: %v", player.MileageMapState.StageTotalScore)
	helper.DebugOut("Point: %v", player.MileageMapState.Point)
	helper.DebugOut("request.Score: %v", request.Score)

	incentives := constobjs.GetMileageIncentives(player.MileageMapState.Episode, player.MileageMapState.Chapter) // Game wants incentives in _current_ episode-chapter
	var oldRewardEpisode, newRewardEpisode int64
	var oldRewardChapter, newRewardChapter int64
	var oldRewardPoint, newRewardPoint int64

	if request.Closed == 0 { // If the game wasn't exited out of
		oldRewardEpisode = player.MileageMapState.Episode
		oldRewardChapter = player.MileageMapState.Chapter
		oldRewardPoint = player.MileageMapState.Point
		helper.DebugOut("Old player ring count: %v", player.PlayerState.NumRings)
		player.PlayerState.NumRings += request.Rings
		player.OptionUserResult.NumTakeAllRings += request.Rings
		helper.DebugOut("Old player red ring count: %v", player.PlayerState.NumRedRings)
		player.PlayerState.NumRedRings += request.RedRings
		helper.DebugOut("New player ring count: %v", player.PlayerState.NumRings)
		helper.DebugOut("New player red ring count: %v", player.PlayerState.NumRedRings)
		player.OptionUserResult.NumTakeAllRedRings += request.RedRings
		player.PlayerState.NumRouletteTicket += request.RedRings // TODO: URGENT! Remove as soon as possible!
		player.PlayerState.Animals += request.Animals
		playerHighScore := player.PlayerState.HighScore
		if request.Score > playerHighScore {
			player.PlayerState.HighScore = request.Score
		}
		playerHighDistance := player.PlayerState.HighDistance
		if request.Distance > playerHighDistance {
			player.PlayerState.HighDistance = request.Distance
		}
		player.PlayerState.TotalDistance += request.Distance
		if time.Now().UTC().Unix() > player.PlayerState.WeeklyScoresExpireAt {
			if player.PlayerState.TotalScoreThisPeriod > player.PlayerState.TotalScore {
				player.PlayerState.TotalScore = player.PlayerState.TotalScoreThisPeriod
			}
			if player.PlayerState.TimedTotalScoreThisPeriod > player.PlayerState.TimedTotalScore {
				player.PlayerState.TimedTotalScore = player.PlayerState.TimedTotalScoreThisPeriod
			}
			player.PlayerState.HighScoreThisPeriod = 0
			player.PlayerState.TimedHighScoreThisPeriod = 0
			player.PlayerState.TotalScoreThisPeriod = 0
			player.PlayerState.TimedTotalScoreThisPeriod = 0
			player.PlayerState.WeeklyScoresExpireAt = now.EndOfWeek().UTC().Unix()
		}
		playerHighScoreThisPeriod := player.PlayerState.HighScoreThisPeriod
		if request.Score > playerHighScoreThisPeriod {
			player.PlayerState.HighScoreThisPeriod = request.Score
		}
		player.PlayerState.TotalScoreThisPeriod += request.Score
		if player.PlayerState.TotalScoreThisPeriod > player.OptionUserResult.TotalSumHighScore {
			player.OptionUserResult.TotalSumHighScore = player.PlayerState.TotalScoreThisPeriod
		}
		// increase character(s)'s experience
		expIncrease := request.Rings + request.FailureRings // all rings collected
		abilityIndex := 1
		for abilityIndex == 1 { // unused ability is at index 1
			abilityIndex = rand.Intn(len(mainC.AbilityLevel))
		}
		// check that increases exist
		_, ok := consts.UpgradeIncreases[mainC.ID]
		if !ok {
			helper.InternalErr("Error getting upgrade increase for main character", fmt.Errorf("no key '%s' in consts.UpgradeIncreases", mainC.ID))
			return
		}
		if hasSubCharacter {
			_, ok = consts.UpgradeIncreases[subC.ID]
			if !ok {
				helper.InternalErr("Error getting upgrade increase for sub character", fmt.Errorf("no key '%s' in consts.UpgradeIncreases", subC.ID))
				return
			}
		}
		playCharacters[0].AbilityLevelUp = []int64{}
		playCharacters[0].AbilityLevelUpExp = []int64{}
		if playCharacters[0].Level < 100 {
			playCharacters[0].Exp += expIncrease
			for playCharacters[0].Exp >= playCharacters[0].Cost {
				// more exp than cost = level up
				if playCharacters[0].Level < 100 {
					abilityIndex = 1
					for abilityIndex == 1 || playCharacters[0].AbilityLevel[abilityIndex] >= 10 { // unused ability is at index 1
						abilityIndex = rand.Intn(len(mainC.AbilityLevel))
					}
					playCharacters[0].Level++                                               // increase level
					playCharacters[0].AbilityLevel[abilityIndex]++                          // increase ability level
					playCharacters[0].Exp -= playCharacters[0].Cost                         // remove cost from exp
					playCharacters[0].Cost += consts.UpgradeIncreases[playCharacters[0].ID] // increase cost
					playCharacters[0].AbilityLevelUp = append(playCharacters[0].AbilityLevelUp, int64(abilityIndex))
					playCharacters[0].AbilityLevelUpExp = append(playCharacters[0].AbilityLevelUpExp, playCharacters[0].Cost)
				} else {
					playCharacters[0].Exp -= playCharacters[0].Cost
				}
			}
		}

		if hasSubCharacter {
			playCharacters[1].AbilityLevelUp = []int64{}
			playCharacters[1].AbilityLevelUpExp = []int64{}
			if playCharacters[1].Level < 100 {
				playCharacters[1].Exp += expIncrease
				for playCharacters[1].Exp >= playCharacters[1].Cost {
					// more exp than cost = level up
					if playCharacters[1].Level < 100 {
						abilityIndex = 1
						for abilityIndex == 1 || playCharacters[1].AbilityLevel[abilityIndex] >= 10 { // unused ability is at index 1
							abilityIndex = rand.Intn(len(playCharacters[1].AbilityLevel))
						}
						playCharacters[1].Level++                                               // increase level
						playCharacters[1].AbilityLevel[abilityIndex]++                          // increase ability level
						playCharacters[1].Exp -= playCharacters[1].Cost                         // remove cost from exp
						playCharacters[1].Cost += consts.UpgradeIncreases[playCharacters[1].ID] // increase cost
						playCharacters[1].AbilityLevelUp = append(playCharacters[1].AbilityLevelUp, int64(abilityIndex))
						playCharacters[1].AbilityLevelUpExp = append(playCharacters[1].AbilityLevelUpExp, playCharacters[1].Cost)
					} else {
						playCharacters[1].Exp -= playCharacters[1].Cost
					}
				}
			}
		}

		helper.DebugOut("Old mainC Exp: %v / %v", mainC.Exp, mainC.Cost)
		helper.DebugOut("Old mainC Level: %v", mainC.Level)
		if hasSubCharacter {
			helper.DebugOut("Old subC Exp: %v / %v", subC.Exp, subC.Cost)
			helper.DebugOut("Old subC Level: %v", subC.Level)
		}
		helper.DebugOut("New mainC Exp: %v / %v", playCharacters[0].Exp, playCharacters[0].Cost)
		helper.DebugOut("New mainC Level: %v", playCharacters[0].Level)
		helper.DebugOut("mainC AbilityLevelUp: %v", playCharacters[0].AbilityLevelUp)
		helper.DebugOut("mainC AbilityLevelUpExp: %v", playCharacters[0].AbilityLevelUpExp)
		if hasSubCharacter {
			helper.DebugOut("New subC Exp: %v / %v", playCharacters[1].Exp, playCharacters[1].Cost)
			helper.DebugOut("New subC Level: %v", playCharacters[1].Level)
			helper.DebugOut("subC AbilityLevelUp: %v", playCharacters[1].AbilityLevelUp)
			helper.DebugOut("subC AbilityLevelUpExp: %v", playCharacters[1].AbilityLevelUpExp)
		}
		/*playCharacters = []netobj.Character{ // TODO: check if this redefinition is needed
			mainC,
			subC,
		}*/
		doStoryProgression := true
		if request.EventId != 0 { // Is this an event stage?
			if strconv.Itoa(int(request.EventId))[1:] == "1" {
				// This is a special stage; don't do story progression since it'll screw with the current point.
				doStoryProgression = false
			}
			helper.DebugOut("Event ID: %v", request.EventId)
			helper.DebugOut("Player got %v event object(s)", request.EventValue)
			player.EventState.Param += request.EventValue
			//TODO: Send rewards to gift box
		}
		if doStoryProgression {
			player.MileageMapState.StageTotalScore += request.Score

			goToNextChapter := request.ChapterClear == 1
			chaoEggs := request.GetChaoEgg
			if chaoEggs > 0 {
				player.PlayerState.ChaoEggs += chaoEggs
				if player.PlayerState.ChaoEggs > 10 {
					player.PlayerState.ChaoEggs = 10
				}
				player.ChaoRouletteGroup.ChaoWheelOptions = netobj.DefaultChaoWheelOptions(player.PlayerState)
			}
			// TODO: Add chao eggs to player
			newPoint := request.ReachPoint

			goToNextEpisode := true
			if goToNextChapter {
				// Assumed this just means next episode...
				if player.PlayerState.Rank < 998 {
					player.PlayerState.Rank++ //TODO: This should be looked into more.
					if player.PlayerState.Energy < player.PlayerVarious.EnergyRecoveryMax {
						player.PlayerState.Energy = player.PlayerVarious.EnergyRecoveryMax //restore energy
					}
				}
				maxChapters, episodeHasMultipleChapters := consts.EpisodeWithChapters[player.MileageMapState.Episode]
				if episodeHasMultipleChapters {
					goToNextEpisode = false
					player.MileageMapState.Chapter++
					player.MileageMapState.Point = 0
					player.MileageMapState.StageTotalScore = 0
					if player.MileageMapState.Chapter > maxChapters {
						// there's no more chapters for this episode!
						goToNextEpisode = true
					}
				}
				if goToNextEpisode {
					player.MileageMapState.Episode++
					player.MileageMapState.Chapter = 1
					player.MileageMapState.Point = 0
					player.MileageMapState.StageTotalScore = 0
					helper.DebugOut("goToNextEpisode -> Episode: %v", player.MileageMapState.Episode)
					if config.CFile.Debug {
						player.MileageMapState.Episode = 15
					}
				}
				if player.MileageMapState.Episode > 50 { // if beat game, reset to 50-1
					player.MileageMapState.Episode = 50
					player.MileageMapState.Chapter = 1
					player.MileageMapState.Point = 0
					player.MileageMapState.StageTotalScore = 0
					helper.DebugOut("goToNextEpisode: Player (%s) beat the game!", player.ID)
				}
			} else {
				player.MileageMapState.Point = newPoint
			}
			if config.CFile.Debug {
				if player.MileageMapState.Episode < 14 {
					player.MileageMapState.Episode = 14
				}
			}
			newRewardEpisode = player.MileageMapState.Episode
			newRewardChapter = player.MileageMapState.Chapter
			newRewardPoint = player.MileageMapState.Point
			// add rewards to PlayerState
			wonRewards := campaign.GetWonRewards(oldRewardEpisode, oldRewardChapter, oldRewardPoint, newRewardEpisode, newRewardChapter, newRewardPoint)
			helper.DebugOut("wonRewards length: %v", wonRewards)
			helper.DebugOut("Previous rings: %v", player.PlayerState.NumRings)
			newItems := player.PlayerState.Items
			for _, reward := range wonRewards { // TODO: This is O(n^2). Maybe alleviate this?
				helper.DebugOut("Reward: %s", reward.ItemID)
				helper.DebugOut("Reward amount: %v", reward.NumItem)
				if reward.ItemID[2:] == "12" { // ID is an item
					// check if the item is already in the player's inventory
					for _, item := range player.PlayerState.Items {
						if item.ID == reward.ItemID { // item found, increment amount
							item.Amount += reward.NumItem
							break
						}
					}
				} else if reward.ItemID == strconv.Itoa(enums.ItemIDRing) { // Rings
					player.PlayerState.NumRings += reward.NumItem
				} else if reward.ItemID == strconv.Itoa(enums.ItemIDRedRing) { // Red rings
					player.PlayerState.NumRedRings += reward.NumItem
				} else if reward.ItemID == enums.CTStrTails { // Tails node
					tailsIndex := player.IndexOfChara(enums.CTStrTails)
					player.CharacterState[tailsIndex].Status = enums.CharacterStatusUnlocked
				} else if reward.ItemID == enums.CTStrKnuckles { // Knuckles node
					knucklesIndex := player.IndexOfChara(enums.CTStrKnuckles)
					player.CharacterState[knucklesIndex].Status = enums.CharacterStatusUnlocked
				} else {
					helper.Out("Unknown reward '" + reward.ItemID + "', ignoring")
				}
			}
			helper.DebugOut("Current rings: %v", player.PlayerState.NumRings)
			player.PlayerState.Items = newItems
		}
	}

	helper.DebugOut("Chapter: %v", player.MileageMapState.Chapter)
	helper.DebugOut("Episode: %v", player.MileageMapState.Episode)
	helper.DebugOut("StageTotalScore: %v", player.MileageMapState.StageTotalScore)
	helper.DebugOut("Point: %v", player.MileageMapState.Point)
	helper.DebugOut("request.Score: %v", request.Score)

	mainCIndex := player.IndexOfChara(mainC.ID) // TODO: check if -1
	subCIndex := -1
	if hasSubCharacter {
		subCIndex = player.IndexOfChara(subC.ID) // TODO: check if -1
	}

	respPlayer := player
	if request.Version == "1.1.4" { // must send fewer characters
		// only get first 21 characters
		// TODO: enforce order 300000 to 300020?
		//cState = cState[:len(cState)-(len(cState)-10)]
		cState := respPlayer.CharacterState
		cState = cState[:16]
		helper.DebugOut("cState length: %v", len(cState))
		helper.DebugOut("Sent character IDs: ")
		for _, char := range cState {
			helper.DebugOut(char.ID)
		}
		respPlayer.CharacterState = cState
	}
	response := responses.DefaultPostGameResults(baseInfo, respPlayer, playCharacters, incentives, respPlayer.EventState)
	err = helper.SendResponse(response)
	if err != nil {
		helper.InternalErr("Error sending response", err)
		return
	}
	// apply the save after the response so that we don't break the leveling
	mainC = playCharacters[0]
	if hasSubCharacter {
		subC = playCharacters[1]
	}
	player.CharacterState[mainCIndex] = mainC
	if hasSubCharacter {
		player.CharacterState[subCIndex] = subC
	}
	helper.DebugOut("CheatResult: %s", request.CheatResult)
	err = db.SavePlayer(player)
	if err != nil {
		helper.InternalErr("Error saving player", err)
		return
	}
	helper.DebugOut(fmt.Sprintf("%v", player.PlayerState.Items))

	_, err = analytics.Store(player.ID, factors.AnalyticTypeStoryEnds)
	if err != nil {
		helper.WarnErr("Error storing analytics (AnalyticTypeStoryEnds)", err)
	}
}

func GetFreeItemList(helper *helper.Helper) {
	baseInfo := helper.BaseInfo(emess.OK, status.OK)
	var response responses.FreeItemListResponse
	if gameconf.CFile.AllItemsFree {
		response = responses.DefaultFreeItemList(baseInfo)
	} else {
		response = responses.FreeItemList(baseInfo, []obj.Item{}) // No free items
	}
	err := helper.SendResponse(response)
	if err != nil {
		helper.InternalErr("Error sending response", err)
	}
}

func GetMileageReward(helper *helper.Helper) {
	recv := helper.GetGameRequest()
	var request requests.MileageRewardRequest
	err := json.Unmarshal(recv, &request)
	if err != nil {
		helper.Err("Error unmarshalling", err)
		return
	}
	/*
		player, err := helper.GetCallingPlayer()
		if err != nil {
			helper.InternalErr("Error getting calling player", err)
			return
		}
	*/
	baseInfo := helper.BaseInfo(emess.OK, status.OK)
	response := responses.DefaultMileageReward(baseInfo, request.Chapter, request.Episode)
	err = helper.SendResponse(response)
	if err != nil {
		helper.InternalErr("Error sending response", err)
	}
}

func DrawRaidBoss(helper *helper.Helper) {
	recv := helper.GetGameRequest()
	var request requests.DrawRaidBossRequest
	err := json.Unmarshal(recv, &request)
	if err != nil {
		helper.Err("Error unmarshalling", err)
		return
	}
	baseInfo := helper.BaseInfo(emess.OK, status.OK)
	response := responses.DrawRaidBoss(baseInfo, netobj.DefaultRaidbossState())
	err = helper.SendCompatibleResponse(response)
	if err != nil {
		helper.InternalErr("Error sending response", err)
	}
}
