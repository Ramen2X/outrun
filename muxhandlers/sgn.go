package muxhandlers

import (
	"encoding/json"

	"github.com/fluofoxxo/outrun/emess"
	"github.com/fluofoxxo/outrun/helper"
	"github.com/fluofoxxo/outrun/requests"
	"github.com/fluofoxxo/outrun/responses"
	"github.com/fluofoxxo/outrun/status"
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
	err = helper.SendResponse(response)
	if err != nil {
		helper.InternalErr("Error sending response", err)
	}
}
