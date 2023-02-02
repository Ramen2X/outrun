package muxhandlers

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"

	"github.com/fluofoxxo/outrun/analytics"
	"github.com/fluofoxxo/outrun/analytics/factors"
	"github.com/fluofoxxo/outrun/config"
	"github.com/fluofoxxo/outrun/config/gameconf"
	"github.com/fluofoxxo/outrun/consts"
	"github.com/fluofoxxo/outrun/db"
	"github.com/fluofoxxo/outrun/emess"
	"github.com/fluofoxxo/outrun/enums"
	"github.com/fluofoxxo/outrun/helper"
	"github.com/fluofoxxo/outrun/logic/campaign"
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
	helper.DebugOut(fmt.Sprintf("%v", player.PlayerState.Items)) // TODO: get rid of this
	baseInfo := helper.BaseInfo(emess.OK, status.OK)
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
	response := responses.DefaultMileageData(baseInfo, player)
	err = helper.SendResponse(response)
	if err != nil {
		helper.InternalErr("Error sending response", err)
	}
}

func GetCampaignList(helper *helper.Helper) {
	baseInfo := helper.BaseInfo(emess.OK, status.OK)
	response := responses.DefaultCampaignList(baseInfo)
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
	responseStatus := status.OK
	// consume items
	modToStringSlice := func(ns []int64) []string {
		result := []string{}
		for _, n := range ns {
			result = append(result, fmt.Sprintf("%v", n))
		}
		return result
	}
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
				player.PlayerState.Items[index].Amount -= 1
			} else {
				if player.PlayerState.NumRings < consumedRings { // not enough rings
					responseStatus = status.NotEnoughRings
					break
				}
				player.PlayerState.NumRings -= consumedRings
			}
		}
	}
	helper.DebugOut(fmt.Sprintf("%v", player.PlayerState.Items))
	baseInfo := helper.BaseInfo(emess.OK, responseStatus)
	response := responses.DefaultQuickActStart(baseInfo, player)
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
				player.PlayerState.Items[index].Amount -= 1
			} else {
				if player.PlayerState.NumRings < consumedRings { // not enough rings
					responseStatus = status.NotEnoughRings
					break
				}
				player.PlayerState.NumRings -= consumedRings
			}
		}
	}
	helper.DebugOut(fmt.Sprintf("%v", player.PlayerState.Items))
	baseInfo := helper.BaseInfo(emess.OK, responseStatus)
	response := responses.DefaultActStart(baseInfo, player)
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
	responseStatus := status.OK
	if player.PlayerState.NumRedRings >= 5 { // enough red rings
		player.PlayerState.NumRedRings -= 5
		err = db.SavePlayer(player)
		if err != nil {
			helper.InternalErr("Error saving player", err)
			return
		}
	} else {
		responseStatus = status.NotEnoughRedRings
	}
	baseInfo := helper.BaseInfo(emess.OK, responseStatus)
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
		playerTimedHighScore := player.PlayerState.TimedHighScore
		if request.Score > playerTimedHighScore {
			player.PlayerState.TimedHighScore = request.Score
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
		if mainC.Level < 100 {
			mainC.Exp += expIncrease
			for mainC.Exp >= mainC.Cost {
				// more exp than cost = level up
				mainC.Level++                                   // increase level
				mainC.AbilityLevel[abilityIndex]++              // increase ability level
				mainC.Exp -= mainC.Cost                         // remove cost from exp
				mainC.Cost += consts.UpgradeIncreases[mainC.ID] // increase cost
				playCharacters = []netobj.Character{            // TODO: check if this redefinition is needed
					mainC,
				}
			}
		}
		if hasSubCharacter {
			if subC.Level < 100 {
				subC.Exp += expIncrease
				for subC.Exp >= subC.Cost {
					// more exp than cost = level up
					subC.Level++                                  // increase level
					subC.AbilityLevel[abilityIndex]++             // increase ability level
					subC.Exp -= subC.Cost                         // remove cost from exp
					subC.Cost += consts.UpgradeIncreases[subC.ID] // increase cost
				}
			}
			playCharacters = []netobj.Character{ // TODO: check if this redefinition is needed
				mainC,
				subC,
			}
		}
	}

	mainCIndex := player.IndexOfChara(mainC.ID) // TODO: check if -1
	subCIndex := -1
	if hasSubCharacter {
		subCIndex = player.IndexOfChara(subC.ID) // TODO: check if -1
	}

	baseInfo := helper.BaseInfo(emess.OK, status.OK)
	response := responses.DefaultQuickPostGameResults(baseInfo, player, playCharacters)
	// apply the save after the response so that we don't break the leveling
	player.CharacterState[mainCIndex] = mainC
	if hasSubCharacter {
		player.CharacterState[subCIndex] = subC
	}
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
		player.PlayerState.NumRings += request.Rings
		player.PlayerState.NumRedRings += request.RedRings
		player.PlayerState.NumRouletteTicket += request.RedRings // TODO: URGENT! Remove as soon as possible!
		player.PlayerState.Animals += request.Animals
		playerHighScore := player.PlayerState.HighScore
		if request.Score > playerHighScore {
			player.PlayerState.HighScore = request.Score
		}
		player.PlayerState.TotalDistance += request.Distance
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
		if mainC.Level < 100 {
			mainC.Exp += expIncrease
			for mainC.Exp >= mainC.Cost {
				// more exp than cost = level up
				mainC.Level++                                   // increase level
				mainC.AbilityLevel[abilityIndex]++              // increase ability level
				mainC.Exp -= mainC.Cost                         // remove cost from exp
				mainC.Cost += consts.UpgradeIncreases[mainC.ID] // increase cost
			}
			playCharacters = []netobj.Character{ // TODO: check if this redefinition is needed
				mainC,
			}
		}
		if hasSubCharacter {
			if subC.Level < 100 {
				subC.Exp += expIncrease
				for subC.Exp >= subC.Cost {
					// more exp than cost = level up
					subC.Level++                                  // increase level
					subC.AbilityLevel[abilityIndex]++             // increase ability level
					subC.Exp -= subC.Cost                         // remove cost from exp
					subC.Cost += consts.UpgradeIncreases[subC.ID] // increase cost
				}
			}
			playCharacters = []netobj.Character{ // TODO: check if this redefinition is needed
				mainC,
				subC,
			}
		}

		player.MileageMapState.StageTotalScore += request.Score

		goToNextChapter := request.ChapterClear == 1
		//chaoEggs := request.GetChaoEgg
		// TODO: Add chao eggs to player
		newPoint := request.ReachPoint

		goToNextEpisode := true
		if goToNextChapter {
			// Assumed this just means next episode...
			maxChapters, episodeHasMultipleChapters := consts.EpisodeWithChapters[player.MileageMapState.Episode]
			if episodeHasMultipleChapters {
				goToNextEpisode = false
				player.MileageMapState.Chapter++
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
				player.PlayerState.Rank++
				if player.PlayerState.Rank > 998 { // rank going above 999
					player.PlayerState.Rank = 998
				}
				helper.DebugOut("goToNextEpisode -> Episode: %v", player.MileageMapState.Episode)
				if config.CFile.Debug {
					player.MileageMapState.Episode = 11
				}
			}
			if player.MileageMapState.Episode > 50 { // if beat game, reset to 50-1
				player.MileageMapState.Episode = 50
				player.MileageMapState.Chapter = 1
				player.MileageMapState.Point = 0
				player.MileageMapState.StageTotalScore = 0
				player.PlayerState.Rank++
				if player.PlayerState.Rank > 998 { // rank going above 999
					player.PlayerState.Rank = 998
				}
				helper.DebugOut("goToNextEpisode: Player (%s) beat the game!", player.ID)
			}
		} else {
			player.MileageMapState.Point = newPoint
		}
		if config.CFile.Debug {
			if player.MileageMapState.Episode < 11 {
				//player.MileageMapState.Episode = 11
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
			if reward.ItemID[:2] == "12" { // ID is an item
				// check if the item is already in the player's inventory
				itemIndex := player.IndexOfItem(reward.ItemID)
				if itemIndex != -1 {
					player.PlayerState.Items[itemIndex].Amount += reward.NumItem
				} else {
					helper.Warn("Unknown item reward '%s'", reward.ItemID)
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
				helper.Warn("Unknown reward '%s', ignoring", reward.ItemID)
			}
			// TODO: allow for any character joining the cast
		}
		helper.DebugOut("Current rings: %v", player.PlayerState.NumRings)
		player.PlayerState.Items = newItems
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

	baseInfo := helper.BaseInfo(emess.OK, status.OK)
	response := responses.DefaultPostGameResults(baseInfo, player, playCharacters, incentives)
	// apply the save after the response so that we don't break the leveling
	player.CharacterState[mainCIndex] = mainC
	if hasSubCharacter {
		player.CharacterState[subCIndex] = subC
	}
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
