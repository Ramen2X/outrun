package netobj

type ItemSpinPrize struct {
	ID     string `json:"itemId"`
	Level  int64  `json:"level"`
	Rarity int64  `json:"rarity"`
}

func ItemIDToItemSpinPrize(itemid string) ItemSpinPrize {
	id := itemid
	level := int64(0) // TODO: check if the game is accepting of this value...
	rarity := int64(100)
	return ItemSpinPrize{
		id,
		level,
		rarity,
	}
}
