package netobj

import (
	"github.com/Ramen2X/outrun/logic/roulette"
)

type ChaoRouletteGroup struct {
	ChaoWheelOptions ChaoWheelOptions `json:"ORN_lastChaoWheelOptions"` // actual wheel options for this wheel
	WheelChao        []string         `json:"ORN_wheelChao"`            // what Chao/characters are in this wheel
	ChaoRouletteInfo RouletteInfo     `json:"ORN_chaoRouletteInfo"`     // may not be needed
}

func DefaultChaoRouletteGroup(playerState PlayerState, allowedCharacters, allowedChao []string, oldCharacters bool) ChaoRouletteGroup {
	chaoWheelOptions := DefaultChaoWheelOptions(playerState)
	//wheelChao, err := roulette.GetRandomChaoRouletteItems(chaoWheelOptions.Rarity, exclusions) // populate based on given rarities
	wheelChao, newRarity, err := roulette.GetRandomChaoRouletteItems(chaoWheelOptions.Rarity, allowedCharacters, allowedChao, oldCharacters) // populate based on given rarities
	if err != nil {
		panic(err) // TODO: Find a better way to handle error. Hard to manage since the player creators don't already output errors
	}
	// newRarity is rarity but with any modifications that need to be made
	chaoWheelOptions.Rarity = newRarity
	chaoRouletteInfo := DefaultRouletteInfo()
	return ChaoRouletteGroup{
		chaoWheelOptions,
		wheelChao,
		chaoRouletteInfo,
	}
}
