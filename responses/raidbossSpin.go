package responses

import (
	"strconv"

	"github.com/fluofoxxo/outrun/enums"
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

type PrizeRaidbossWheelResponse struct {
	BaseResponse
	PrizeList []obj.RaidbossPrize `json:"raidbossPrizeList"`
}

func PrizeRaidbossWheel(base responseobjs.BaseInfo, prizeList []obj.RaidbossPrize) PrizeRaidbossWheelResponse {
	baseResponse := NewBaseResponse(base)
	out := PrizeRaidbossWheelResponse{
		baseResponse,
		prizeList,
	}
	return out
}

func DefaultPrizeRaidbossWheel(base responseobjs.BaseInfo) PrizeRaidbossWheelResponse {
	prizeList := []obj.RaidbossPrize{
		obj.NewRaidbossPrize(strconv.Itoa(int(enums.ItemIDAsteroid)), 5, 10, 0),
		obj.NewRaidbossPrize(strconv.Itoa(int(enums.ItemIDRedRing)), 15, 1, 0),
	}
	return PrizeRaidbossWheel(base, prizeList)
}
