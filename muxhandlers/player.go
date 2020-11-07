package muxhandlers

import (
	"encoding/json"
	"math/rand"
	"strconv"
	"time"

	"github.com/fluofoxxo/outrun/consts"

	"github.com/fluofoxxo/outrun/db"
	"github.com/fluofoxxo/outrun/emess"
	"github.com/fluofoxxo/outrun/helper"
	"github.com/fluofoxxo/outrun/requests"
	"github.com/fluofoxxo/outrun/responses"
	"github.com/fluofoxxo/outrun/status"
	"github.com/jinzhu/now"
)

func GetPlayerState(helper *helper.Helper) {
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
	if player.PlayerState.NumRings < 0 {
		player.PlayerState.NumRings = 0
	}
	if player.PlayerState.NumRedRings < 0 {
		player.PlayerState.NumRedRings = 0
	}
	for time.Now().UTC().Unix() >= player.PlayerState.EnergyRenewsAt && player.PlayerState.Energy < player.PlayerVarious.EnergyRecoveryMax {
		player.PlayerState.Energy++
		player.PlayerState.EnergyRenewsAt += player.PlayerVarious.EnergyRecoveryTime
	}
	if player.PlayerState.NextNumDailyChallenge <= 0 || int(player.PlayerState.NextNumDailyChallenge) > len(consts.DailyMissionRewards) {
		player.PlayerState.NumDailyChallenge = int64(0)
		player.PlayerState.NextNumDailyChallenge = int64(1)
		player.PlayerState.DailyChalCatNum = int64(rand.Intn(5))
	}
	if time.Now().UTC().Unix() >= player.PlayerState.DailyMissionEndTime {
		if player.PlayerState.DailyChallengeComplete == 1 && player.PlayerState.DailyChalSetNum < 10 {
			helper.DebugOut("Advancing to next daily mission...")
			player.PlayerState.DailyChalSetNum++
		} else {
			player.PlayerState.DailyChalCatNum = int64(rand.Intn(5))
			player.PlayerState.DailyChalSetNum = int64(0)
		}
		if player.PlayerState.DailyChallengeComplete == 0 {
			player.PlayerState.NumDailyChallenge = int64(0)
			player.PlayerState.NextNumDailyChallenge = int64(1)
		} else {
			player.PlayerState.NextNumDailyChallenge++
			if int(player.PlayerState.NextNumDailyChallenge) > len(consts.DailyMissionRewards) {
				player.PlayerState.NumDailyChallenge = int64(0)
				player.PlayerState.NextNumDailyChallenge = int64(1) //restart from beginning
				player.PlayerState.DailyChalCatNum = int64(rand.Intn(5))
				player.PlayerState.DailyChalSetNum = int64(0)
			}
		}
		player.PlayerState.DailyChalPosNum = int64(1 + rand.Intn(2))
		player.PlayerState.DailyMissionID = int64((player.PlayerState.DailyChalCatNum * 33) + (player.PlayerState.DailyChalSetNum * 3) + player.PlayerState.DailyChalPosNum)
		player.PlayerState.DailyChallengeValue = int64(0)
		player.PlayerState.DailyChallengeComplete = int64(0)
		player.PlayerState.DailyMissionEndTime = now.EndOfDay().UTC().Unix() + 1
		helper.DebugOut("New daily mission ID: %v", player.PlayerState.DailyMissionID)
	}
	err = db.SavePlayer(player)
	if err != nil {
		helper.InternalErr("Error saving player", err)
		return
	}
	response := responses.PlayerState(baseInfo, player.PlayerState)
	helper.SendResponse(response)
}

func GetCharacterState(helper *helper.Helper) {
	src := helper.GetGameRequest()
	var request requests.Base
	err := json.Unmarshal(src, &request)
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
	response := responses.CharacterState(baseInfo, cState)
	helper.SendResponse(response)
}

func GetChaoState(helper *helper.Helper) {
	src := helper.GetGameRequest()
	var request requests.Base
	err := json.Unmarshal(src, &request)
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
	//cap max levels to prevent hang
	chaoState := player.ChaoState
	maxLevel := 10
	if request.Version == "1.1.4" {
		maxLevel = 5
	}
	for index, chao := range chaoState {
		if int(chao.Level) > maxLevel {
			chaoState[index].Level = int64(maxLevel)
		}
	}
	response := responses.ChaoState(baseInfo, chaoState)
	helper.SendResponse(response)
}

func SetUsername(helper *helper.Helper) {
	recv := helper.GetGameRequest()
	var request requests.SetUsernameRequest
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
	// TODO: check if username is already taken
	player.Username = request.Username
	err = db.SavePlayer(player)
	if err != nil {
		helper.InternalErr("Error saving player", err)
		return
	}
	baseInfo := helper.BaseInfo(emess.OK, status.OK)
	response := responses.NewBaseResponse(baseInfo)
	err = helper.SendResponse(response)
	if err != nil {
		helper.InternalErr("Error sending response", err)
		return
	}
}
