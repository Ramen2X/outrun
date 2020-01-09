package netobj

import (
	"log"
	"math/rand"
	"strconv"

	"github.com/fluofoxxo/outrun/config"
	"github.com/fluofoxxo/outrun/consts"
	"github.com/fluofoxxo/outrun/enums"
	"github.com/fluofoxxo/outrun/logic/roulette"
	"github.com/fluofoxxo/outrun/obj"
	"github.com/jinzhu/now"
)

type RaidbossWheelOptions struct {
	Items                []string       `json:"items"`
	Item                 []int64        `json:"item"`
	ItemWeight           []int64        `json:"itemWeight"`
	ItemWon              int64          `json:"itemWon"`
	NextFreeSpin         int64          `json:"nextFreeSpin"` // midnight (start of next day)
	RouletteRank         int64          `json:"rouletteRank"`
	NumSpecialEgg        int64          `json:"numSpecialEgg"`
	NumRouletteToken     int64          `json:"numRouletteToken"`
	NumJackpotRing       int64          `json:"numJackpotRing"`
	NumRemainingRoulette int64          `json:"numRemainingRoulette"`
	SpinID               int64          `json:"spinID"`
	CostItemList         []obj.CostItem `json:"costItemList"`
}

func DefaultRaidbossWheelOptions(numRouletteTicket, numSpecialEgg, rouletteCountInPeriod, rouletteRank, freeSpins int64) RaidbossWheelOptions {
	items := []string{strconv.Itoa(enums.IDTypeItemRouletteWin)} // first item is always jackpot/big/super
	item := []int64{1}
	itemWeight := []int64{1250, 1250, 1250, 1250, 1250, 1250, 1250, 1250}
	chaoIDs, err := roulette.GetRandomChaoWheelChao(2, 7, true)
	if err != nil {
		log.Printf("error generating chao list: %s\n", err.Error())
	}
	if config.CFile.Debug {
		for _ = range make([]byte, 7) { // loop 7 times
			items = append(items, chaoIDs[len(items)-1])
			item = append(item, 1)
		}
	} else {
		for _ = range make([]byte, 7) { // loop 7 times
			items = append(items, strconv.Itoa(enums.ItemIDRing))
			item = append(item, 1000)
		}
	}
	itemWon := int64(rand.Intn(len(items)))   //TODO: adjust this to accurately represent item weights
	nextFreeSpin := now.EndOfDay().Unix() + 1 // midnight
	numJackpotRing := int64(consts.RouletteJackpotRings)
	// TODO: get rid of logic here!
	costItemList := []obj.CostItem{
		obj.NewCostItem(
			strconv.Itoa(enums.ItemIDRaidbossRing),
			5000,
			27272,
		),
	}
	numRouletteToken := numRouletteTicket
	if numRouletteToken > 0 {
		costItemList[0] = obj.NewCostItem(
			strconv.Itoa(enums.IDRouletteTicketRaid),
			1,
			numRouletteToken,
		)
	}
	numRemainingRoulette := numRouletteToken + freeSpins - rouletteCountInPeriod // TODO: is this proper?
	if numRemainingRoulette < numRouletteToken {
		numRemainingRoulette = numRouletteToken
	}
	spinID := int64(0)
	out := RaidbossWheelOptions{
		items,
		item,
		itemWeight,
		itemWon,
		nextFreeSpin,
		rouletteRank,
		numSpecialEgg,
		numRouletteToken,
		numJackpotRing,
		numRemainingRoulette,
		spinID,
		costItemList,
	}
	return out
}
