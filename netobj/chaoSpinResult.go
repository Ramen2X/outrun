package netobj

import "github.com/fluofoxxo/outrun/obj"

type ChaoSpinResult struct {
	WonPrize ChaoSpinPrize `json:"getChao"` // chao or character
	ItemList []obj.Item    `json:"itemList"`
	ItemWon  int64         `json:"itemWon"` // probably index of item in ItemList
}

type ChaoSpinResult2 struct { //TODO: Research this.
	WonPrize ItemSpinPrize `json:"getItem"` // item??????????
	ItemList []obj.Item    `json:"itemList"`
	ItemWon  int64         `json:"itemWon"` // probably index of item in ItemList
}

func DefaultChaoSpinResultNoItems(wonPrize ChaoSpinPrize) ChaoSpinResult {
	return ChaoSpinResult{
		wonPrize,
		[]obj.Item{},
		-1,
	}
}
