package obj

type RaidbossPrize struct {
	ItemID     string `json:"itemId"`
	NumItem    int64  `json:"numItem,string"`
	ItemWeight int64  `json:"itemRate,string"`
	SpinID     int64  `json:"spinId,string"`
}

func NewRaidbossPrize(id string, numitem, weight, spinid int64) RaidbossPrize {
	return RaidbossPrize{
		id,
		numitem,
		weight,
		spinid,
	}
}
