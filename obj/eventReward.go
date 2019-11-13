package obj

type EventReward struct {
    RewardID int64  `json:"rewardId"` 
    Param    int64  `json:"param"`    
    ItemID   string `json:"itemId"`   
    NumItem  int64  `json:"numItem"`   
}

func NewEventReward(rewardId, param int64, itemId string, numItem int64) EventReward {
    return EventReward{
        rewardId,
        param,
        itemId,
        numItem,
    }
}
