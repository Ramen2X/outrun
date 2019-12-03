package netobj

import (
	"time"
)

type EventUserRaidbossState struct {
	NumRaidbossRings       int64 `json:"numRaidbossRings"` // NOTE: these are different from normal rings!
	RaidBossEnergy         int64 `json:"raidbossEnergy"`
	RaidBossEnergyBuy      int64 `json:"raidbossEnergyBuy"`
	NumBeatedEncounter     int64 `json:"numBeatedEncounter"`  // number of times the boss has been defeated?
	NumBeatedEnterprise    int64 `json:"numBeatedEnterprise"` // ?
	NumRaidBossEncountered int64 `json:"numTotalEncountered"`
	EnergyRenewsAt         int64 `json:"raidbossEnergyRenewsAt"`
}

func DefaultUserRaidbossState() EventUserRaidbossState {
	// TODO: establish as constants
	numRaidbossRings := int64(0)
	raidbossEnergy := int64(20)
	raidbossEnergyBuy := int64(0)
	numBeatedEncounter := int64(0)
	numBeatedEnterprise := int64(0)
	numRaidBossEncountered := int64(0)
	energyRenewsAt := int64(time.Now().Unix() + 600) // in ten minutes
	return EventUserRaidbossState{
		numRaidbossRings,
		raidbossEnergy,
		raidbossEnergyBuy,
		numBeatedEncounter,
		numBeatedEnterprise,
		numRaidBossEncountered,
		energyRenewsAt,
	}
}
