package netobj

import "time"

type EventRaidbossState struct {
	ID               int64  `json:"raidbossId"`
	Level            int64  `json:"raidbossLevel"`
	Rarity           int64  `json:"raidbossRarity"` // 0: Normal, 1: Rare, 2: Super-rare
	HP               int64  `json:"raidbossHitPoint"`
	MaxHP            int64  `json:"raidbossMaxHitPoint"`
	Status           int64  `json:"raidbossStatus"`   // TODO: research
	EscapeAt         int64  `json:"raidbossEscapeAt"` // the time when the raid boss expires
	EncounterName    string `json:"encounterName"`
	EncounterFlg     int64  `json:"encounterFlg"`
	CrowdedFlg       int64  `json:"crowdedFlg"`       // raid boss is full?
	ParticipateCount int64  `json:"participateCount"` // number of people taking on this raid boss?
}

func NewRaidbossState(id, level, rarity int64, encounterName string) EventRaidbossState {
	maxHp := int64((15*(level+3) - 50) * (rarity + 1)) //TODO: Preliminary algorithm
	hp := maxHp
	status := int64(0)
	escapeAt := int64(time.Now().Unix() + 3600) // raid boss expires after 1 hour
	encounterFlg := int64(0)
	crowdedFlg := int64(0)
	participateCount := int64(0)
	return EventRaidbossState{
		id,
		level,
		rarity,
		hp,
		maxHp,
		status,
		escapeAt,
		encounterName,
		encounterFlg,
		crowdedFlg,
		participateCount,
	}
}

func DefaultRaidbossState() EventRaidbossState {
	// TODO: establish as constants
	id := int64(0)
	level := int64(8)
	rarity := int64(2)
	encounterName := "System"
	return NewRaidbossState(
		id,
		level,
		rarity,
		encounterName,
	)
}
