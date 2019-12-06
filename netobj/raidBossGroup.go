package netobj

type RaidBossGroup struct {
	CurrentRaidBoss      EventRaidbossState   `json:"ORN_currentRaidBoss"`
	PreviousRaidBosses   []EventRaidbossState `json:"ORN_previousRaidBosses`
	AverageRaidBossLevel int64                `json:"ORN_averageRaidBossLevel"`
	NextID               int64                `json:"ORN_nextId"`
}
