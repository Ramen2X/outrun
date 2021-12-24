package responses

import (
	"strconv"

	"github.com/Ramen2X/outrun/consts"
	"github.com/Ramen2X/outrun/enums"
	"github.com/Ramen2X/outrun/netobj"
	"github.com/Ramen2X/outrun/obj"
	"github.com/Ramen2X/outrun/responses/responseobjs"
)

type ChaoWheelOptionsResponse struct {
	BaseResponse
	ChaoWheelOptions netobj.ChaoWheelOptions `json:"chaoWheelOptions"`
}

func ChaoWheelOptions(base responseobjs.BaseInfo, chaoWheelOptions netobj.ChaoWheelOptions) ChaoWheelOptionsResponse {
	baseResponse := NewBaseResponse(base)
	out := ChaoWheelOptionsResponse{
		baseResponse,
		chaoWheelOptions,
	}
	return out
}

func DefaultChaoWheelOptions(base responseobjs.BaseInfo, player netobj.Player) ChaoWheelOptionsResponse {
	// TODO: Assess if needed
	chaoWheelOptions := player.ChaoRouletteGroup.ChaoWheelOptions
	return ChaoWheelOptions(
		base,
		chaoWheelOptions,
	)
}

type PrizeChaoWheelResponse struct {
	BaseResponse
	PrizeList []obj.ChaoPrize `json:"prizeList"`
}

func PrizeChaoWheel(base responseobjs.BaseInfo, prizeList []obj.ChaoPrize) PrizeChaoWheelResponse {
	baseResponse := NewBaseResponse(base)
	out := PrizeChaoWheelResponse{
		baseResponse,
		prizeList,
	}
	return out
}

func DefaultPrizeChaoWheel(base responseobjs.BaseInfo) PrizeChaoWheelResponse {
	//prizeList := constobjs.DefaultChaoPrizeWheelPrizeList
	prizeList := []obj.ChaoPrize{}
	chaoids := []string{}
	chaorarities := []int64{}
	for chid := range consts.RandomChaoWheelChaoPrizes {
		rarity, _ := strconv.Atoi(string(chid[2]))
		chaoids = append(chaoids, chid)
		chaorarities = append(chaorarities, int64(rarity))
	}
	for chid := range consts.RandomChaoWheelCharacterPrizes {
		chaoids = append(chaoids, chid)
		chaorarities = append(chaorarities, int64(100))
	}
	for index := range chaoids {
		prizeList = append(prizeList, obj.NewChaoPrize(chaoids[index], chaorarities[index]))
	}
	return PrizeChaoWheel(base, prizeList)
}

type EquipChaoResponse struct {
	BaseResponse
	PlayerState netobj.PlayerState `json:"playerState"`
}

func EquipChao(base responseobjs.BaseInfo, playerState netobj.PlayerState) EquipChaoResponse {
	baseResponse := NewBaseResponse(base)
	return EquipChaoResponse{
		baseResponse,
		playerState,
	}
}

type ChaoWheelSpinResponse struct {
	BaseResponse
	PlayerState      netobj.PlayerState      `json:"playerState"`
	CharacterState   []netobj.Character      `json:"characterState"`
	ChaoState        []netobj.Chao           `json:"chaoState"` // also works with json:"chaoStatus"
	ChaoWheelOptions netobj.ChaoWheelOptions `json:"chaoWheelOptions"`
	ChaoSpinResults  []netobj.ChaoSpinResult `json:"chaoSpinResultList"` // Should only contain one element! Otherwise, ItemWon is interpreted as -1
	ChaoSpinResult   netobj.ChaoSpinResult   `json:"chaoSpinResult"`     // This is only there for 1.1.4 backwards compatibility
}

func ChaoWheelSpin(base responseobjs.BaseInfo, playerState netobj.PlayerState, characterState []netobj.Character, chaoState []netobj.Chao, chaoWheelOptions netobj.ChaoWheelOptions, chaoSpinResults []netobj.ChaoSpinResult) ChaoWheelSpinResponse {
	baseResponse := NewBaseResponse(base)
	return ChaoWheelSpinResponse{
		baseResponse,
		playerState,
		characterState,
		chaoState,
		chaoWheelOptions,
		chaoSpinResults,
		chaoSpinResults[0],
	}
}

func DefaultChaoWheelSpin(base responseobjs.BaseInfo, player netobj.Player) ChaoWheelSpinResponse {
	// WARN: Do not use for normal purposes!! This should only be used for debugging
	dummyPrize := netobj.CharacterIDToChaoSpinPrize(enums.CTStrShadow)
	chaoSpinResults := netobj.DefaultChaoSpinResultNoItems(dummyPrize)
	return ChaoWheelSpin(
		base,
		player.PlayerState,
		player.CharacterState,
		player.ChaoState,
		player.ChaoRouletteGroup.ChaoWheelOptions,
		[]netobj.ChaoSpinResult{chaoSpinResults},
	)
}
