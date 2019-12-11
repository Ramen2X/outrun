package campaignconf

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/fluofoxxo/outrun/enums"
)

// defaults
// default variable names MUST be "D" + (nameOfVariable)
var Defaults = map[string]interface{}{
	"DAllowCampaigns":   false,
	"DCurrentCampaigns": []ConfiguredCampaign{},
}

var CampaignTypes = map[string]int64{
	"bankedRingBonus":                 enums.CampaignTypeBankedRingBonus,
	"chaoRouletteCost":                enums.CampaignTypeChaoRouletteCost,
	"characterUpgradeCost":            enums.CampaignTypeCharacterUpgradeCost,
	"continueCost":                    enums.CampaignTypeContinueCost,
	"dailyMissionBonus":               enums.CampaignTypeDailyMissionBonus,
	"freeWheelSpinCount":              enums.CampaignTypeFreeWheelSpinCount,
	"gameItemCost":                    enums.CampaignTypeGameItemCost,
	"inviteCount":                     enums.CampaignTypeInviteCount,
	"jackpotValue":                    enums.CampaignTypeJackpotValueBonus,
	"mileagePassingRingBonus":         enums.CampaignTypeMileagePassingRingBonus,
	"premiumRouletteOdds":             enums.CampaignTypePremiumRouletteOdds,
	"purchaseAddEnergy":               enums.CampaignTypePurchaseAddEnergies,
	"purchaseAddRaidEnergy":           enums.CampaignTypePurchaseAddRaidEnergies,
	"purchaseAddRedRings":             enums.CampaignTypePurchaseAddRedRings,
	"purchaseAddRedRingsNoChargeUser": enums.CampaignTypePurchaseAddRedRingsNoChargeUser,
	"purchaseAddRings":                enums.CampaignTypePurchaseAddRings,
	"sendAddEnergy":                   enums.CampaignTypeSendAddEnergies,
}

type ConfiguredCampaign struct {
	Type       string `json:"type"`
	Content    int64  `json:"content"`
	SubContent int64  `json:"subContent"`
	StartTime  int64  `json:"startTime"` // *
	EndTime    int64  `json:"endTime"`   // *
}

func (c ConfiguredCampaign) RealType() int64 {
	return CampaignTypes[c.Type]
}
func (c ConfiguredCampaign) HasValidType() bool {
	_, ok := CampaignTypes[c.Type]
	return ok
}

var CFile ConfigFile

type ConfigFile struct {
	AllowCampaigns   bool                 `json:"allowCampaigns,omitempty"`
	CurrentCampaigns []ConfiguredCampaign `json:"currentCampaigns,omitempty"`
}

func Parse(filename string) error {
	CFile = ConfigFile{
		Defaults["DAllowCampaigns"].(bool),
		Defaults["DCurrentCampaigns"].([]ConfiguredCampaign),
	}
	file, err := loadFile(filename)
	if err != nil {
		return err
	}
	err = json.Unmarshal(file, &CFile)
	if err != nil {
		return err
	}
	newCampaigns := []ConfiguredCampaign{}
	for i, cc := range CFile.CurrentCampaigns {
		if !cc.HasValidType() {
			log.Printf("[WARN] Invalid campaign type %s at index %v, ignoring\n", cc.Type, i)
			continue
		}
		newCampaigns = append(newCampaigns, cc)
	}
	CFile.CurrentCampaigns = newCampaigns
	return nil
}

func loadFile(filename string) ([]byte, error) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return []byte{}, err
	}
	return b, err
}
