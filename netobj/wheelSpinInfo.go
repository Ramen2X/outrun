package netobj

import (
	"time"
)

type WheelSpinInfo struct {
	ID    string `json:"id"`
	Start int64  `json:"start"`
	End   int64  `json:"end"`
	Param string `json:"param"`
}

func NewWheelSpinInfo(id, param string) WheelSpinInfo {
	return WheelSpinInfo{
		id,
		time.Now().UTC().Unix(),
		time.Now().UTC().Unix() + 7300, // 2 hours + 100s from now
		param,
	}
}

func DefaultWheelSpinInfoList() []WheelSpinInfo {
	//TODO: Should this be specifiable in a roulette configuration file?
	return []WheelSpinInfo{
		/*NewWheelSpinInfo("1", "This"),
		NewWheelSpinInfo("2", "is"),
		NewWheelSpinInfo("3", "a"),
		NewWheelSpinInfo("4", "test"),
		NewWheelSpinInfo("5", "message,"),
		NewWheelSpinInfo("6", "but"),
		NewWheelSpinInfo("7", "not"),
		NewWheelSpinInfo("8", "joined!"),*/
		NewWheelSpinInfo("1337", "Welcome to the Item Roulette, where you can spend your roulette tickets to win items, rings, or red rings! You get 5 free spins per day! Can you win the jackpot?"),
	}
}
