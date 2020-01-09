package obj

type CostItem struct {
	ID        string `json:"itemId"`
	Amount    int64  `json:"numItem,string"`
	ItemStock int64  `json:"itemStock,string`
}

func NewCostItem(iid string, amount, stock int64) CostItem {
	return CostItem{
		iid,
		amount,
		stock,
	}
}
