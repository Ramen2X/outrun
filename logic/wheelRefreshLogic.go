package logic

import (
	"github.com/fluofoxxo/outrun/config/campaignconf"
	"github.com/fluofoxxo/outrun/consts"
	"github.com/fluofoxxo/outrun/enums"
	"github.com/fluofoxxo/outrun/logic/conversion"
	"github.com/fluofoxxo/outrun/netobj"
	"github.com/fluofoxxo/outrun/obj"
)

func WheelRefreshLogic(player netobj.Player, wheel netobj.WheelOptions) netobj.WheelOptions {
	// TODO: Find a more standard way of refreshing the wheel status, because this is scary code
	freeSpins := consts.RouletteFreeSpins
	campaignList := []obj.Campaign{}
	if campaignconf.CFile.AllowCampaigns {
		for _, confCampaign := range campaignconf.CFile.CurrentCampaigns {
			newCampaign := conversion.ConfiguredCampaignToCampaign(confCampaign)
			campaignList = append(campaignList, newCampaign)
		}
	}
	index := 0
	for index < len(campaignList) {
		if obj.IsCampaignActive(campaignList[index]) && campaignList[index].Type == enums.CampaignTypeFreeWheelSpinCount {
			freeSpins = campaignList[index].Content
		}
		index++
	}
	numRouletteTicket := player.PlayerState.NumRouletteTicket  // get roulette tickets
	rouletteCount := player.RouletteInfo.RouletteCountInPeriod // get amount of times we've spun the wheel today
	if player.RouletteInfo.GotJackpotThisPeriod {
		wheel.NumJackpotRing = 1
	}
	wheel.NumRouletteToken = numRouletteTicket
	wheel.NumRemainingRoulette = wheel.NumRouletteToken + -rouletteCount // TODO: is this proper?
	if wheel.NumRemainingRoulette < wheel.NumRouletteToken {
		wheel.NumRemainingRoulette = wheel.NumRouletteToken
	}

	return wheel
}
