package obj

import (
	"time"

	"github.com/Ramen2X/outrun/enums"
	"github.com/jinzhu/now"
)

type Campaign struct {
	Type       int64 `json:"campaignType"`
	Content    int64 `json:"campaignContent"`
	SubContent int64 `json:"campaignSubContent"`
	StartTime  int64 `json:"campaignStartTime"`
	EndTime    int64 `json:"campaignEndTime"`
}

func NewCampaign(ctype, content, subcontent, startTime, endTime int64) Campaign {
	return Campaign{
		ctype,
		content,
		subcontent,
		startTime,
		endTime,
	}
}

func DefaultCampaign(ctype, content, subcontent int64) Campaign {
	return NewCampaign(
		ctype,
		content,
		subcontent,
		now.BeginningOfDay().UTC().Unix(),
		now.EndOfDay().UTC().Unix(),
	)
}

func DefaultCampaigns() []Campaign {
	return []Campaign{
		DefaultCampaign(enums.CampaignTypeBankedRingBonus, 250, 0), // 25 percent ring boost
	}
}

func IsCampaignActive(campaign Campaign) bool {
	return time.Now().UTC().Unix() >= campaign.StartTime && time.Now().UTC().Unix() < campaign.EndTime
}
