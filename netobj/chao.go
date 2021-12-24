package netobj

import (
	"github.com/Ramen2X/outrun/obj"
)

type Chao struct {
	obj.Chao
	Status    int64 `json:"status"` // enums.ChaoStatus*
	Level     int64 `json:"level"`
	Dealing   int64 `json:"setStatus"` // enums.ChaoDealing*
	NumInvite int64 `json:"numInvite"` // ?
	Acquired  int64 `json:"acquired"`  // appears in the game code as NumAcquired
}

func NewNetChao(chao obj.Chao, status, level, dealing, acquired int64) Chao {
	return Chao{
		chao,
		status,
		level,
		dealing,
		int64(0),
		acquired,
	}
}
