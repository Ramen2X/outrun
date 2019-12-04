package obj

import (
	"github.com/fluofoxxo/outrun/config/campaignconf"
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
	campaignList := []Campaign{}
	if campaignconf.CFile.AllowCampaigns {
		for _, confCampaign := range campaignconf.CFile.CurrentCampaigns {
			newCampaign := conversion.ConfiguredCampaignToCampaign(confCampaign)
			campaignList = append(campaignList, newCampaign)
		}
	}
	helper.DebugOut("Campaign list: %v", campaignList)
	return campaignList
	/*return []Campaign{
		DefaultCampaign(enums.CampaignTypeBankedRingBonus, 250, 0), // 25 percent ring boost
	}*/
}
