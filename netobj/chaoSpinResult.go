package netobj

import "github.com/fluofoxxo/outrun/obj"

type ChaoSpinResult struct {
	WonPrize ChaoSpinPrize `json:"getChao"` // chao or character
	ItemList []obj.Item    `json:"itemList"`
	ItemWon  int64         `json:"itemWon"` // probably index of item in ItemList
}

type ChaoSpinResult2 struct {
	WonPrize ItemSpinPrize `json:"getItem"` // item??????????
	ItemList []obj.Item    `json:"itemList"`
	ItemWon  int64         `json:"itemWon"` // probably index of item in ItemList
}

func DefaultChaoSpinResultNoItems(wonPrize ChaoSpinPrize) ChaoSpinResult {
	return ChaoSpinResult{
		wonPrize,
		[]obj.Item{},
		-1, // TODO: 1.1.4 doesn't seem to like this. Perhaps something should be there after all?
	}
}

func DefaultChaoSpinResult(wonPrize ChaoSpinPrize, itemList []obj.Item, itemWon int64) ChaoSpinResult {
	return ChaoSpinResult{
		wonPrize,
		itemList,
		itemWon,
	}
}
