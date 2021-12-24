package muxhandlers

import (
	"encoding/json"
	"strconv"

	"github.com/Ramen2X/outrun/analytics"
	"github.com/Ramen2X/outrun/analytics/factors"
	"github.com/Ramen2X/outrun/config/campaignconf"
	"github.com/Ramen2X/outrun/db"
	"github.com/Ramen2X/outrun/emess"
	"github.com/Ramen2X/outrun/enums"
	"github.com/Ramen2X/outrun/helper"
	"github.com/Ramen2X/outrun/logic/conversion"
	"github.com/Ramen2X/outrun/obj"
	"github.com/Ramen2X/outrun/obj/constobjs"
	"github.com/Ramen2X/outrun/requests"
	"github.com/Ramen2X/outrun/responses"
	"github.com/Ramen2X/outrun/status"
)

func GetRedStarExchangeList(helper *helper.Helper) {
	recv := helper.GetGameRequest()
	var request requests.RedStarExchangeListRequest
	err := json.Unmarshal(recv, &request)
	if err != nil {
		helper.Err("Error unmarshalling", err)
		return
	}
	campaignList := []obj.Campaign{}
	if campaignconf.CFile.AllowCampaigns {
		for _, confCampaign := range campaignconf.CFile.CurrentCampaigns {
			newCampaign := conversion.ConfiguredCampaignToCampaign(confCampaign)
			campaignList = append(campaignList, newCampaign)
		}
	}
	helper.DebugOut("Campaign list: %v", campaignList)
	baseInfo := helper.BaseInfo(emess.OK, status.OK)
	var response responses.RedStarExchangeListResponse
	var redStarItems []obj.RedStarItem
	helper.Out("Recv ItemType " + strconv.Itoa(int(request.ItemType)))
	if request.ItemType == 0 { //red star rings
		if request.Version == "2.0.3" {
			redStarItems = []obj.RedStarItem{}
		} else {
			redStarItems = constobjs.RedStarItemsType0
		}
	} else if request.ItemType == 1 { // rings
		redStarItems = constobjs.RedStarItemsType1
	} else if request.ItemType == 2 { // energy
		redStarItems = constobjs.RedStarItemsType2
	} else if request.ItemType == 4 { // raid boss energy
		redStarItems = constobjs.RedStarItemsType4
	} else {
		helper.InvalidRequest()
		return
	}
	index := 0
	campaign := obj.DefaultCampaign(enums.CampaignTypeBankedRingBonus, 2000, 0)
	campaignActive := false
	for index < len(campaignList) {
		if obj.IsCampaignActive(campaignList[index]) {
			switch request.ItemType {
			case 0: //red star rings
				if campaignList[index].Type == enums.CampaignTypePurchaseAddRedRings || campaignList[index].Type == enums.CampaignTypePurchaseAddRedRingsNoChargeUser {
					campaign = campaignList[index]
					campaignActive = true
				}
			case 1: //rings
				if campaignList[index].Type == enums.CampaignTypePurchaseAddRings {
					campaign = campaignList[index]
					campaignActive = true
				}
			case 2: //energy
				if campaignList[index].Type == enums.CampaignTypePurchaseAddEnergies {
					campaign = campaignList[index]
					campaignActive = true
				}
			case 4: //raid boss energy
				if campaignList[index].Type == enums.CampaignTypePurchaseAddRaidEnergies {
					campaign = campaignList[index]
					campaignActive = true
				}
			}
		}
		index++
	}
	index = 0
	for index < len(redStarItems) {
		if campaignActive {
			redStarItems[index].Campaign = &campaign
		}
		index++
	}
	response = responses.RedStarExchangeList(baseInfo, redStarItems, 0, "1900-1-1")
	err = helper.SendResponse(response)
	if err != nil {
		helper.InternalErr("Error sending response", err)
	}
}

func RedStarExchange(helper *helper.Helper) {
	recv := helper.GetGameRequest()
	var request requests.RedStarExchange
	err := json.Unmarshal(recv, &request)
	if err != nil {
		helper.Err("Error unmarshalling", err)
		return
	}

	player, err := helper.GetCallingPlayer(true)
	if err != nil {
		helper.InternalErr("Error getting calling player", err)
		return
	}

	baseInfo := helper.BaseInfo(emess.OK, status.OK)
	if player.Suspended {
		baseInfo.StatusCode = status.MissingPlayer
		err = helper.SendResponse(responses.NewBaseResponse(baseInfo))
		if err != nil {
			helper.InternalErr("Error sending response", err)
			return
		}
		return
	}
	itemID := request.ItemID
	getItemType := func(iid string) (string, int64, bool) {
		var itemType string
		itemPrice, ok := constobjs.ShopRingPrices[iid]
		if !ok {
			// it's not a ring item
			itemPrice, ok = constobjs.ShopEnergyPrices[iid]
			if !ok {
				// it's not an energy item
				itemPrice, ok = constobjs.ShopRaidbossEnergyPrices[iid]
				if !ok {
					// it's not ring, energy, nor raidboss energy item
					// unrecognized item!
					return "", 0, false
				}
				itemType = "raidbossEnergy"
				return itemType, itemPrice, true
			}
			itemType = "energy"
			return itemType, itemPrice, true
		}
		itemType = "ring"
		return itemType, itemPrice, true
	}

	itemType, itemPrice, found := getItemType(itemID)
	if found {
		switch itemType {
		case "ring":
			if player.PlayerState.NumRedRings-itemPrice < 0 {
				baseInfo.StatusCode = status.NotEnoughRedRings
				return
			}
			player.PlayerState.NumRedRings -= itemPrice
			player.PlayerState.NumRings += constobjs.ShopRingAmounts[itemID]
			//player.PlayerState.NumBuyRings += constobjs.ShopRingAmounts[itemID]
			db.SavePlayer(player)
			_, err = analytics.Store(player.ID, factors.AnalyticTypePurchaseRings)
			if err != nil {
				helper.WarnErr("Error storing analytics (AnalyticTypePurchaseRings)", err)
			}
		case "energy":
			if player.PlayerState.NumRedRings-itemPrice < 0 {
				baseInfo.StatusCode = status.NotEnoughRedRings
				return
			}
			player.PlayerState.NumRedRings -= itemPrice
			//player.PlayerState.Energy += constobjs.ShopEnergyAmounts[itemID]
			player.PlayerState.EnergyBuy += constobjs.ShopEnergyAmounts[itemID]
			db.SavePlayer(player)
			_, err = analytics.Store(player.ID, factors.AnalyticTypePurchaseEnergy)
			if err != nil {
				helper.WarnErr("Error storing analytics (AnalyticTypePurchaseEnergy)", err)
			}
		case "raidbossEnergy":
			if player.PlayerState.NumRedRings-itemPrice < 0 {
				baseInfo.StatusCode = status.NotEnoughRedRings
				return
			}
			player.PlayerState.NumRedRings -= itemPrice
			//player.PlayerState.Energy += constobjs.ShopEnergyAmounts[itemID]
			player.EventUserRaidbossState.RaidBossEnergyBuy += constobjs.ShopRaidbossEnergyAmounts[itemID]
			db.SavePlayer(player)
		default:
			// this should never execute!
			baseInfo.StatusCode = status.MasterDataMismatch
			helper.Out("Default case executed... Something went wrong!")
		}
	} else {
		baseInfo.StatusCode = status.MasterDataMismatch
	}

	response := responses.DefaultRedStarExchange(baseInfo, player)
	err = helper.SendResponse(response)
	if err != nil {
		helper.InternalErr("Error sending response", err)
	}
}
