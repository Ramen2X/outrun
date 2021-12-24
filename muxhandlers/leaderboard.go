package muxhandlers

import (
	"encoding/json"
	"strconv"

	"github.com/Ramen2X/outrun/emess"
	"github.com/Ramen2X/outrun/helper"
	"github.com/Ramen2X/outrun/requests"
	"github.com/Ramen2X/outrun/responses"
	"github.com/Ramen2X/outrun/status"
)

func GetWeeklyLeaderboardOptions(helper *helper.Helper) {
	recv := helper.GetGameRequest()
	var request requests.LeaderboardRequest
	err := json.Unmarshal(recv, &request)
	if err != nil {
		helper.Err("Error unmarshalling", err)
		return
	}
	mode := request.Mode
	baseInfo := helper.BaseInfo(emess.OK, status.OK)
	response := responses.DefaultWeeklyLeaderboardOptions(baseInfo, mode)
	err = helper.SendResponse(response)
	if err != nil {
		helper.InternalErr("Error sending response", err)
	}
}

func GetWeeklyLeaderboardEntries(helper *helper.Helper) {
	recv := helper.GetGameRequest()
	var request requests.LeaderboardEntriesRequest
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
	mode := request.Mode
	first := request.First
	lbtype := request.Type
	baseInfo := helper.BaseInfo(emess.OK, status.OK)
	helper.DebugOut("Start from: " + strconv.Itoa(int(first)))
	helper.DebugOut("Mode: " + strconv.Itoa(int(mode)))
	helper.DebugOut("Type: " + strconv.Itoa(int(lbtype)))
	response := responses.DefaultWeeklyLeaderboardEntries(baseInfo, player, mode, lbtype, first)
	err = helper.SendResponse(response)
	if err != nil {
		helper.InternalErr("Error sending response", err)
	}
}

func GetLeagueData(helper *helper.Helper) {
	recv := helper.GetGameRequest()
	var request requests.LeaderboardRequest
	err := json.Unmarshal(recv, &request)
	if err != nil {
		helper.Err("Error unmarshalling", err)
		return
	}
	mode := request.Mode
	baseInfo := helper.BaseInfo(emess.OK, status.OK)
	response := responses.DefaultLeagueData(baseInfo, mode)
	err = helper.SendResponse(response)
	if err != nil {
		helper.InternalErr("Error sending response", err)
	}
}
