package consts

import (
	"math/rand"
	"strconv"

	"github.com/fluofoxxo/outrun/enums"
)

type AmountRange struct {
	Min  int64
	Max  int64
	Step int64
}

func (a AmountRange) GetRandom() int64 {
	// construct random list first
	randomSelections := []int64{}
	diff := int64(0)
	currMin := a.Min
	for diff >= 0 {
		randomSelections = append(randomSelections, currMin)
		currMin += a.Step
		diff = a.Max - currMin
	}
	selectionIndex := rand.Intn(len(randomSelections))
	selection := randomSelections[selectionIndex]
	return selection
}

// The game does not support RingBonus, DistanceBonus,
// or AnimalBonus on the normal wheel.

// NOTE: If you remove an item from NormalWheelItemAmountRange
// but don't remove it from RandomItemListNormalWheel, you're going
// to create a memory leak.

var RandomItemListNormalWheel = []string{
	enums.ItemIDStrInvincible,
	enums.ItemIDStrBarrier,
	enums.ItemIDStrMagnet,
	enums.ItemIDStrTrampoline,
	enums.ItemIDStrCombo,
	enums.ItemIDStrLaser,
	enums.ItemIDStrDrill,
	enums.ItemIDStrAsteroid,
	strconv.Itoa(enums.IDTypeRedRing),
	//strconv.Itoa(enums.IDTypeItemRouletteWin),
}

var NormalWheelItemAmountRange = map[string]AmountRange{
	enums.ItemIDStrInvincible:         AmountRange{1, 5, 1},
	enums.ItemIDStrBarrier:            AmountRange{1, 5, 1},
	enums.ItemIDStrMagnet:             AmountRange{1, 5, 1},
	enums.ItemIDStrTrampoline:         AmountRange{1, 5, 1},
	enums.ItemIDStrCombo:              AmountRange{1, 5, 1},
	enums.ItemIDStrLaser:              AmountRange{1, 5, 1},
	enums.ItemIDStrDrill:              AmountRange{1, 5, 1},
	enums.ItemIDStrAsteroid:           AmountRange{1, 5, 1},
	strconv.Itoa(enums.IDTypeRedRing): AmountRange{5, 25, 5},
	//strconv.Itoa(enums.IDTypeItemRouletteWin): AmountRange{1, 1, 1},
}

var RandomItemListBigWheel = []string{
	enums.ItemIDStrInvincible,
	enums.ItemIDStrBarrier,
	enums.ItemIDStrMagnet,
	enums.ItemIDStrTrampoline,
	enums.ItemIDStrCombo,
	enums.ItemIDStrLaser,
	enums.ItemIDStrDrill,
	enums.ItemIDStrAsteroid,
	strconv.Itoa(enums.IDTypeRedRing),
	//strconv.Itoa(enums.IDTypeItemRouletteWin),
}

var BigWheelItemAmountRange = map[string]AmountRange{
	enums.ItemIDStrInvincible:         AmountRange{5, 10, 1},
	enums.ItemIDStrBarrier:            AmountRange{5, 10, 1},
	enums.ItemIDStrMagnet:             AmountRange{5, 10, 1},
	enums.ItemIDStrTrampoline:         AmountRange{5, 10, 1},
	enums.ItemIDStrCombo:              AmountRange{5, 10, 1},
	enums.ItemIDStrLaser:              AmountRange{5, 10, 1},
	enums.ItemIDStrDrill:              AmountRange{5, 10, 1},
	enums.ItemIDStrAsteroid:           AmountRange{5, 10, 1},
	strconv.Itoa(enums.IDTypeRedRing): AmountRange{10, 50, 10},
	//strconv.Itoa(enums.IDTypeItemRouletteWin): AmountRange{1, 1, 1},
}

var RandomItemListSuperWheel = []string{
	enums.ItemIDStrInvincible,
	enums.ItemIDStrBarrier,
	enums.ItemIDStrMagnet,
	enums.ItemIDStrTrampoline,
	enums.ItemIDStrCombo,
	enums.ItemIDStrLaser,
	enums.ItemIDStrDrill,
	enums.ItemIDStrAsteroid,
	strconv.Itoa(enums.IDTypeRedRing),
	//strconv.Itoa(enums.IDTypeItemRouletteWin),
}

var SuperWheelItemAmountRange = map[string]AmountRange{
	enums.ItemIDStrInvincible:         AmountRange{10, 20, 2},
	enums.ItemIDStrBarrier:            AmountRange{10, 20, 2},
	enums.ItemIDStrMagnet:             AmountRange{10, 20, 2},
	enums.ItemIDStrTrampoline:         AmountRange{10, 20, 2},
	enums.ItemIDStrCombo:              AmountRange{10, 20, 2},
	enums.ItemIDStrLaser:              AmountRange{10, 20, 2},
	enums.ItemIDStrDrill:              AmountRange{10, 20, 2},
	enums.ItemIDStrAsteroid:           AmountRange{10, 20, 2},
	strconv.Itoa(enums.IDTypeRedRing): AmountRange{20, 160, 20},
	//strconv.Itoa(enums.IDTypeItemRouletteWin): AmountRange{1, 1, 1},
}
