package responses

import (
	"github.com/fluofoxxo/outrun/netobj"
	"github.com/fluofoxxo/outrun/obj"
	"github.com/fluofoxxo/outrun/obj/constobjs"
	"github.com/fluofoxxo/outrun/responses/responseobjs"
)

type ItemStockNumResponse struct {
	BaseResponse
	ItemStockList []obj.Item `json:"itemStockList"`
}

func ItemStockNum(base responseobjs.BaseInfo, itemStockList []obj.Item) ItemStockNumResponse {
	baseResponse := NewBaseResponse(base)
	return ItemStockNumResponse{
		baseResponse,
		itemStockList,
	}
}

func DefaultItemStockNum(base responseobjs.BaseInfo) ItemStockNumResponse {
	return ItemStockNum(
		base,
		constobjs.DefaultSpinItems,
	)
}

type RaidbossWheelOptionsResponse struct {
	BaseResponse
	RaidbossWheelOptions netobj.RaidbossWheelOptions `json:"raidbossWheelOptions"`
}

func RaidbossWheelOptions(base responseobjs.BaseInfo, raidbossWheelOptions netobj.RaidbossWheelOptions) RaidbossWheelOptionsResponse {
	baseResponse := NewBaseResponse(base)
	out := RaidbossWheelOptionsResponse{
		baseResponse,
		raidbossWheelOptions,
	}
	return out
}
