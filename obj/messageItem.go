package obj

type MessageItem struct {
	ID              int64 `json:"itemId,string"`
	Amount          int64 `json:"numItem"`
	AdditionalInfo1 int64 `json:"additionalInfo1"`
	AdditionalInfo2 int64 `json:"additionalInfo2"`
}

func NewMessageItem(id, amount, ai1, ai2 int64) MessageItem {
	return MessageItem{
		id,
		amount,
		ai1,
		ai2,
	}
}

func MessageItemToPresent(item MessageItem) Present {
	return NewPresent(
		item.ID,
		item.Amount,
		item.AdditionalInfo1,
		item.AdditionalInfo2,
	)
}
