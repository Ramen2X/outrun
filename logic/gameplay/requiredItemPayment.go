package gameplay

import (
	"github.com/Ramen2X/outrun/netobj"
	"github.com/Ramen2X/outrun/obj"
	"github.com/Ramen2X/outrun/obj/constobjs"
)

func findItem(id string) obj.ConsumedItem {
	var result obj.ConsumedItem
	for _, citem := range constobjs.DefaultCostList {
		if citem.ID == id {
			result = citem
			break
		}
	}
	return result
}

func GetRequiredItemPayment(items []string, player netobj.Player) int64 {
	totalRingPayment := int64(0)
	for _, itemID := range items {
		citem := findItem(itemID)
		if itemID[:2] == "11" { // boosts, not items
			totalRingPayment += citem.Item.Amount
		} else {
			index := player.IndexOfItem(itemID)
			if player.PlayerState.Items[index].Amount < 1 {
				totalRingPayment += citem.Item.Amount
			}
		}
	}
	return totalRingPayment
}
