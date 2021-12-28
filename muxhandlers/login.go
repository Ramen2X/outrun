package muxhandlers

import (
	"encoding/json"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/Ramen2X/outrun/enums"

	"github.com/Ramen2X/outrun/netobj"
	"github.com/jinzhu/now"

	"github.com/Ramen2X/outrun/analytics"
	"github.com/Ramen2X/outrun/analytics/factors"
	"github.com/Ramen2X/outrun/config/infoconf"
	"github.com/Ramen2X/outrun/db"
	"github.com/Ramen2X/outrun/emess"
	"github.com/Ramen2X/outrun/helper"
	"github.com/Ramen2X/outrun/logic"
	"github.com/Ramen2X/outrun/logic/conversion"
	"github.com/Ramen2X/outrun/obj"
	"github.com/Ramen2X/outrun/obj/constobjs"
	"github.com/Ramen2X/outrun/requests"
	"github.com/Ramen2X/outrun/responses"
	"github.com/Ramen2X/outrun/status"
)

func Login(helper *helper.Helper) {
	recv := helper.GetGameRequest()
	var request requests.LoginRequest
	err := json.Unmarshal(recv, &request)
	if err != nil {
		helper.Err("Error unmarshalling", err)
		return
	}

	uid := request.LineAuth.UserID
	password := request.LineAuth.Password

	baseInfo := helper.BaseInfo(emess.OK, status.OK)
	if uid == "0" && password == "" {
		helper.Out("Entering LoginAlpha (registration)")
		newPlayer := db.NewAccount()
		err = db.SavePlayer(newPlayer)
		if err != nil {
			helper.InternalErr("Error saving player", err)
			return
		}
		baseInfo.StatusCode = status.InvalidPassword
		baseInfo.SetErrorMessage(emess.BadPassword)
		response := responses.LoginRegister(
			baseInfo,
			newPlayer.ID,
			newPlayer.Password,
			newPlayer.Key,
		)
		err = helper.SendCompatibleResponse(response, true)
		if err != nil {
			helper.InternalErr("Error responding", err)
		}
		return
	} else if uid == "0" && password != "" {
		helper.Out("Entering LoginBravo (INVALID)")
		// invalid request
		helper.InvalidRequest()
		return
	} else if uid != "0" && password == "" {
		helper.Out("Entering LoginCharlie (initial login attempt)")
		// game wants to log in
		baseInfo.StatusCode = status.InvalidPassword
		baseInfo.SetErrorMessage(emess.BadPassword)
		player, err := db.GetPlayer(uid)
		if err != nil {
			helper.InternalErr("Error getting player", err)
			// likely account that wasn't found, so let's tell them that:
			response := responses.LoginCheckKey(baseInfo, "")
			baseInfo.StatusCode = status.MissingPlayer
			helper.SendCompatibleResponse(response, true)
			return
		}
		response := responses.LoginCheckKey(baseInfo, player.Key)
		err = helper.SendCompatibleResponse(response, true)
		if err != nil {
			helper.InternalErr("Error sending response", err)
			return
		}
		return
	} else if uid != "0" && password != "" {
		helper.Out("Entering LoginDelta (login with passkey)")
		// game is attempting to log in using key

		baseInfo.StatusCode = status.OK
		baseInfo.SetErrorMessage(emess.OK)
		player, err := db.GetPlayer(uid)
		if err != nil {
			helper.InternalErr("Error getting player", err)
			return
		}
		if player.Suspended {
			baseInfo.StatusCode = status.MissingPlayer
			err = helper.SendResponse(responses.NewBaseResponse(baseInfo))
			if err != nil {
				helper.InternalErr("Error sending response", err)
				return
			}
			return
		}
		if logic.GenerateLoginPassword(player) == request.Password {
			sid, err := db.AssignSessionID(uid)
			if err != nil {
				helper.InternalErr("Error assigning session ID", err)
				return
			}
			player.LastLogin = time.Now().UTC().Unix()
			if player.PlayerState.RankingLeague == int64(enums.RankingLeagueNone) {
				player.PlayerState.RankingLeague = int64(enums.RankingLeagueF_M)
			}
			if player.PlayerState.QuickRankingLeague == int64(enums.RankingLeagueNone) {
				player.PlayerState.QuickRankingLeague = int64(enums.RankingLeagueF_M)
			}
			err = db.SavePlayer(player)
			if err != nil {
				helper.InternalErr("Error saving player", err)
				return
			}
			var response interface{}
			response = responses.LoginSuccess(baseInfo, sid, player.Username)
			if infoconf.CFile.EOLMessageEnabled {
				baseInfo.StatusCode = status.ServerNextVersion
				response = responses.NewNextVersionResponse(baseInfo,
					player.PlayerState.NumRedRings,
					player.PlayerState.NumBuyRedRings,
					player.Username,
					infoconf.CFile.EOLMessageJP,
					infoconf.CFile.EOLMessageEN,
					infoconf.CFile.EOLMessageURL,
				)
			}
			err = helper.SendCompatibleResponse(response, true)
			if err != nil {
				helper.InternalErr("Error sending response", err)
				return
			}
			analytics.Store(player.ID, factors.AnalyticTypeLogins)
			return
		} else {
			baseInfo.StatusCode = status.InvalidPassword
			baseInfo.SetErrorMessage(emess.BadPassword)
			helper.DebugOut("Incorrect passkey sent: \"%s\"", request.Password)
			err = helper.SendResponse(responses.NewBaseResponse(baseInfo))
			if err != nil {
				helper.InternalErr("Error sending response", err)
				return
			}
			return
		}
	}
}

func GetVariousParameter(helper *helper.Helper) {
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
	response := responses.VariousParameter(baseInfo, player)
	err = helper.SendCompatibleResponse(response, true)
	if err != nil {
		helper.InternalErr("Error sending response", err)
		return
	}
}

func GetInformation(helper *helper.Helper) {
	baseInfo := helper.BaseInfo(emess.OK, status.OK)
	infos := []obj.Information{}
	helper.DebugOut("%v", infoconf.CFile.EnableInfos)
	if infoconf.CFile.EnableInfos {
		for _, ci := range infoconf.CFile.Infos {
			newInfo := conversion.ConfiguredInfoToInformation(ci)
			infos = append(infos, newInfo)
			helper.DebugOut(newInfo.Param)
		}
	}
	operatorInfos := []obj.OperatorInformation{}
	numOpUnread := int64(len(operatorInfos))
	response := responses.Information(baseInfo, infos, operatorInfos, numOpUnread)
	err := helper.SendCompatibleResponse(response, true)
	if err != nil {
		helper.InternalErr("Error sending response", err)
	}
}

func GetTicker(helper *helper.Helper) {
	player, err := helper.GetCallingPlayer(true)
	if err != nil {
		helper.InternalErr("Error getting calling player", err)
		return
	}
	baseInfo := helper.BaseInfo(emess.OK, status.OK)
	response := responses.DefaultTicker(baseInfo, player)
	err = helper.SendCompatibleResponse(response, true)
	if err != nil {
		helper.InternalErr("Error sending response", err)
	}
}

func LoginBonus(helper *helper.Helper) {
	player, err := helper.GetCallingPlayer(true)
	if err != nil {
		helper.InternalErr("Error getting calling player", err)
		return
	}
	if time.Now().UTC().Unix() > player.LoginBonusState.LoginBonusEndTime {
		player.LoginBonusState = netobj.DefaultLoginBonusState(player.LoginBonusState.CurrentFirstLoginBonusDay)
	}
	doLoginBonus := false
	if time.Now().UTC().Unix() > player.LoginBonusState.NextLoginBonusTime {
		doLoginBonus = true
		player.LoginBonusState.LastLoginBonusTime = time.Now().UTC().Unix()
		player.LoginBonusState.NextLoginBonusTime = now.EndOfDay().UTC().Unix()
		player.LoginBonusState.CurrentFirstLoginBonusDay++
		player.LoginBonusState.CurrentLoginBonusDay++
	}
	err = db.SavePlayer(player)
	if err != nil {
		helper.InternalErr("Error saving player", err)
		return
	}
	baseInfo := helper.BaseInfo(emess.OK, status.OK)
	response := responses.DefaultLoginBonus(baseInfo, player, doLoginBonus)
	err = helper.SendCompatibleResponse(response, true)
	if err != nil {
		helper.InternalErr("Error sending response", err)
	}
}

func LoginBonusSelect(helper *helper.Helper) {
	recv := helper.GetGameRequest()
	var request requests.LoginBonusSelectRequest
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
	rewardList := []obj.Item{}
	firstRewardList := []obj.Item{}
	if request.FirstRewardDays != -1 && int(request.FirstRewardDays) < len(constobjs.DefaultFirstLoginBonusRewardList) {
		firstRewardList = constobjs.DefaultFirstLoginBonusRewardList[request.FirstRewardDays].SelectRewardList[request.FirstRewardSelect].ItemList
	}
	if request.RewardDays != -1 && int(request.RewardDays) < len(constobjs.DefaultLoginBonusRewardList) {
		rewardList = constobjs.DefaultLoginBonusRewardList[request.RewardDays].SelectRewardList[request.RewardSelect].ItemList
	}
	for _, item := range rewardList {
		itemid, _ := strconv.Atoi(item.ID)
		player.AddOperatorMessage(
			"A Login Bonus.",
			obj.MessageItem{
				int64(itemid),
				item.Amount,
				0,
				0,
			},
			2592000,
		)
		helper.DebugOut("Sent %s x %v to gift box (Login Bonus)", item.ID, item.Amount)
	}
	for _, item := range firstRewardList {
		itemid, _ := strconv.Atoi(item.ID)
		player.AddOperatorMessage(
			"A Debut Dash Login Bonus.",
			obj.MessageItem{
				int64(itemid),
				item.Amount,
				0,
				0,
			},
			2592000,
		)
		helper.DebugOut("Sent %s x %v to gift box (Start Dash Login Bonus)", item.ID, item.Amount)
	}
	err = db.SavePlayer(player)
	if err != nil {
		helper.InternalErr("Error saving player", err)
		return
	}
	baseInfo := helper.BaseInfo(emess.OK, status.OK)
	response := responses.LoginBonusSelect(baseInfo, rewardList, firstRewardList)
	err = helper.SendCompatibleResponse(response, true)
	if err != nil {
		helper.InternalErr("Error sending response", err)
	}
}

func GetCountry(helper *helper.Helper) {
	// TODO: Should get correct country code!
	baseInfo := helper.BaseInfo(emess.OK, status.OK)
	response := responses.DefaultGetCountry(baseInfo)
	err := helper.SendCompatibleResponse(response, true)
	if err != nil {
		helper.InternalErr("Error sending response", err)
	}
}

func GetMigrationPassword(helper *helper.Helper) {
	randChar := func(charset string, length int64) string {
		runes := []rune(charset)
		final := make([]rune, 10)
		for i := range final {
			final[i] = runes[rand.Intn(len(runes))]
		}
		return string(final)
	}
	recv := helper.GetGameRequest()
	var request requests.GetMigrationPasswordRequest
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
	if len(player.MigrationPassword) != 10 {
		player.MigrationPassword = randChar("abcdefghijklmnopqrstuvwxyz1234567890", 10)
	}
	player.UserPassword = request.UserPassword
	db.SavePlayer(player)
	response := responses.MigrationPassword(baseInfo, player)
	err = helper.SendCompatibleResponse(response, true)
	if err != nil {
		helper.InternalErr("Error sending response", err)
	}
}

func Migration(helper *helper.Helper) {
	randChar := func(charset string, length int64) string {
		runes := []rune(charset)
		final := make([]rune, 10)
		for i := range final {
			final[i] = runes[rand.Intn(len(runes))]
		}
		return string(final)
	}

	recv := helper.GetGameRequest()
	var request requests.LoginRequest
	err := json.Unmarshal(recv, &request)
	if err != nil {
		helper.Err("Error unmarshalling", err)
		return
	}
	password := request.LineAuth.MigrationPassword
	migrationUserPassword := request.LineAuth.MigrationUserPassword

	password = strings.Replace(password, "-", "", -1)

	baseInfo := helper.BaseInfo(emess.OK, status.OK)

	helper.DebugOut("Transfer ID: %s", password)
	foundPlayers, err := logic.FindPlayersByMigrationPassword(password, false)
	if err != nil {
		helper.Err("Error finding players by password", err)
		return
	}
	playerIDs := []string{}
	for _, player := range foundPlayers {
		playerIDs = append(playerIDs, player.ID)
	}
	if len(playerIDs) > 0 {
		migratePlayer, err := db.GetPlayer(playerIDs[0])
		if err != nil {
			helper.InternalErr("Error getting player", err)
			return
		}
		if migrationUserPassword == migratePlayer.UserPassword {
			if migratePlayer.Suspended {
				baseInfo.StatusCode = status.MissingPlayer
				err = helper.SendResponse(responses.NewBaseResponse(baseInfo))
				if err != nil {
					helper.InternalErr("Error sending response", err)
					return
				}
				return
			}
			baseInfo.StatusCode = status.OK
			baseInfo.SetErrorMessage(emess.OK)
			migratePlayer.MigrationPassword = randChar("abcdefghijklmnopqrstuvwxyz1234567890", 10) //generate a brand new transfer ID                                                        //clear user password
			migratePlayer.LastLogin = time.Now().UTC().Unix()
			err = db.SavePlayer(migratePlayer)
			if err != nil {
				helper.InternalErr("Error saving player", err)
				return
			}
			sid, err := db.AssignSessionID(migratePlayer.ID)
			if err != nil {
				helper.InternalErr("Error assigning session ID", err)
				return
			}
			helper.DebugOut("User ID: %s", migratePlayer.ID)
			helper.DebugOut("Username: %s", migratePlayer.Username)
			helper.DebugOut("New Transfer ID: %s", migratePlayer.MigrationPassword)
			response := responses.MigrationSuccess(baseInfo, sid, migratePlayer.ID, migratePlayer.Username, migratePlayer.Password)
			helper.SendCompatibleResponse(response, true)
		} else {
			baseInfo.StatusCode = status.InvalidPassword
			baseInfo.SetErrorMessage(emess.BadPassword)
			helper.DebugOut("Incorrect password for user ID %s", migratePlayer.ID)
			response := responses.NewBaseResponse(baseInfo)
			helper.SendCompatibleResponse(response, true)
		}
	} else {
		helper.DebugOut("Failed to find player")
		baseInfo.StatusCode = status.InvalidPassword
		response := responses.NewBaseResponse(baseInfo)
		helper.SendCompatibleResponse(response, true)
	}
}
