package muxhandlers

import (
	"encoding/json"

	"github.com/Ramen2X/outrun/emess"
	"github.com/Ramen2X/outrun/helper"
	"github.com/Ramen2X/outrun/requests"
	"github.com/Ramen2X/outrun/responses"
	"github.com/Ramen2X/outrun/status"
)

func SendApollo(helper *helper.Helper) {
	recv := helper.GetGameRequest()
	var request requests.SendApolloRequest
	err := json.Unmarshal(recv, &request)
	if err != nil {
		helper.Err("Error unmarshalling", err)
		return
	}
	helper.DebugOut("Data type: %v", request.Type)
	if len(request.Value) > 0 {
		index := 0
		for index < len(request.Value) {
			helper.DebugOut("Value %v: \"%s\"", index+1, request.Value[index])
			index++
		}
	} else {
		helper.DebugOut("No data.")
	}
	baseInfo := helper.BaseInfo(emess.OK, status.OK)
	response := responses.NewBaseResponse(baseInfo)
	err = helper.SendCompatibleResponse(response, true)
	if err != nil {
		helper.InternalErr("Error sending response", err)
	}
}

func SetNoahID(helper *helper.Helper) {
	recv := helper.GetGameRequest()
	var request requests.SetNoahIDRequest
	err := json.Unmarshal(recv, &request)
	if err != nil {
		helper.Err("Error unmarshalling", err)
		return
	}
	helper.DebugOut("Noah ID: %v", request.NoahID)
	baseInfo := helper.BaseInfo(emess.OK, status.OK)
	response := responses.NewBaseResponse(baseInfo)
	err = helper.SendCompatibleResponse(response, true)
	if err != nil {
		helper.InternalErr("Error sending response", err)
	}
}
