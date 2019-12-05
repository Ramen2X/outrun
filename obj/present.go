package obj

type Present struct {
	ItemID          int64 `json:"itemId"`
	NumItem         int64 `json:"numItem"`
	AdditionalInfo1 int64 `json:"additionalInfo1"`
	AdditionalInfo2 int64 `json:"additionalInfo2"`
}

func NewPresent(itemId, numItem, additionalInfo1, additionalInfo2 int64) Present {
	return Present{
		itemId,
		numItem,
		additionalInfo1,
		additionalInfo2,
	}
}
