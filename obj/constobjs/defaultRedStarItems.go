package constobjs

import (
	"github.com/Ramen2X/outrun/obj"
)

var RedStarItemsType0 = rsiDefaultsIT0() // Red Rings
var RedStarItemsType1 = rsiDefaultsIT1() // Rings
var RedStarItemsType2 = rsiDefaultsIT2() // Energies
var RedStarItemsType4 = rsiDefaultsIT4() // Raid boss Energies

func rsiDefaultsIT0() []obj.RedStarItem {
	redstaritems := []obj.RedStarItem{}
	storeItemIds := []string{"900010", "900030", "900060", "900210", "900380"}
	itemIds := []string{"900000", "900000", "900000", "900000", "900000"}
	numItems := []int64{10, 30, 60, 210, 380}
	prices := []int64{0, 0, 0, 0, 0}
	priceDisps := []string{"", "", "", "", ""}
	productIds := []string{"", "", "", "", ""}
	for i, storeItemId := range storeItemIds {
		itemId := itemIds[i]
		numItem := numItems[i]
		price := prices[i]
		priceDisp := priceDisps[i]
		productId := productIds[i]
		rsi := obj.NewRedStarItem(
			storeItemId,
			itemId,
			priceDisp,
			productId,
			numItem,
			price,
			nil,
		)
		redstaritems = append(redstaritems, rsi)
	}
	return redstaritems
}

func rsiDefaultsIT1() []obj.RedStarItem {
	redstaritems := []obj.RedStarItem{}
	storeItemIds := []string{"915001", "915002", "915003", "915004", "915005"}
	itemIds := []string{"910000", "910000", "910000", "910000", "910000"}
	numItems := []int64{20000, 42000, 108000, 224000, 600000}
	prices := []int64{10, 20, 50, 100, 250}
	priceDisps := []string{"", "", "", "", ""}
	productIds := []string{"", "", "", "", ""}
	for i, storeItemId := range storeItemIds {
		itemId := itemIds[i]
		numItem := numItems[i]
		price := prices[i]
		priceDisp := priceDisps[i]
		productId := productIds[i]
		rsi := obj.NewRedStarItem(
			storeItemId,
			itemId,
			priceDisp,
			productId,
			numItem,
			price,
			nil,
		)
		redstaritems = append(redstaritems, rsi)
	}
	return redstaritems
}

func rsiDefaultsIT2() []obj.RedStarItem {
	redstaritems := []obj.RedStarItem{}
	storeItemIds := []string{"920005", "920020", "920030", "920050", "920100"}
	itemIds := []string{"920000", "920000", "920000", "920000", "920000"}
	numItems := []int64{5, 20, 30, 50, 100}
	prices := []int64{10, 35, 50, 80, 150}
	priceDisps := []string{"", "", "", "", ""}
	productIds := []string{"", "", "", "", ""}
	for i, storeItemId := range storeItemIds {
		itemId := itemIds[i]
		numItem := numItems[i]
		price := prices[i]
		priceDisp := priceDisps[i]
		productId := productIds[i]
		//campaign := obj.DefaultCampaign(8, 2000, 0) // WARN: is this a possibility for a memory leak?
		rsi := obj.NewRedStarItem(
			storeItemId,
			itemId,
			priceDisp,
			productId,
			numItem,
			price,
			nil,
			//&campaign,
		)
		redstaritems = append(redstaritems, rsi)
	}
	return redstaritems
}

func rsiDefaultsIT4() []obj.RedStarItem {
	redstaritems := []obj.RedStarItem{}
	storeItemIds := []string{"940005", "940020", "940030", "940050", "940100"}
	itemIds := []string{"940000", "940000", "940000", "940000", "940000"}
	numItems := []int64{5, 20, 30, 50, 100}
	prices := []int64{10, 35, 50, 80, 150}
	priceDisps := []string{"", "", "", "", ""}
	productIds := []string{"", "", "", "", ""}
	for i, storeItemId := range storeItemIds {
		itemId := itemIds[i]
		numItem := numItems[i]
		price := prices[i]
		priceDisp := priceDisps[i]
		productId := productIds[i]
		rsi := obj.NewRedStarItem(
			storeItemId,
			itemId,
			priceDisp,
			productId,
			numItem,
			price,
			nil,
		)
		redstaritems = append(redstaritems, rsi)
	}
	return redstaritems
}
