package netobj

import (
	"github.com/jinzhu/now"
)

type LoginBonusState struct {
	CurrentFirstLoginBonusDay int64 `json:"currentFirstLoginBonusDay"` // this doesn't get reset when the login bonus resets
	CurrentLoginBonusDay      int64 `json:"currentLoginBonusDay"`
	LastLoginBonusTime        int64 `json:"lastLoginBonusTime"`
	NextLoginBonusTime        int64 `json:"nextLoginBonusTime"`
	LoginBonusStartTime       int64 `json:"loginBonusStartTime"`
	LoginBonusEndTime         int64 `json:"loginBonusEndTime"`
}

func NewLoginBonusState(cflbd, clbd, llbt, nlbt, lbst, lbet int64) LoginBonusState {
	return LoginBonusState{
		cflbd,
		clbd,
		llbt,
		nlbt,
		lbst,
		lbet,
	}
}

func DefaultLoginBonusState(currentFirstLoginBonusDay int64) LoginBonusState {
	currentLoginBonusDay := int64(0)
	lastLoginBonusTime := int64(0)
	nextLoginBonusTime := int64(0)
	loginBonusStartTime := now.BeginningOfWeek().UTC().Unix()
	loginBonusEndTime := now.EndOfWeek().UTC().Unix()
	return NewLoginBonusState(
		currentFirstLoginBonusDay,
		currentLoginBonusDay,
		lastLoginBonusTime,
		nextLoginBonusTime,
		loginBonusStartTime,
		loginBonusEndTime,
	)
}
