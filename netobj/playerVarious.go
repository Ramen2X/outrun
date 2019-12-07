package netobj

type PlayerVarious struct {
	CmSkipCount          int64 `json:"cmSkipCount"`       // no clear purpose
	EnergyRecoveryMax    int64 `json:"energyRecoveryMax"` // max time-recoverable energy
	EnergyRecoveryTime   int64 `json:"energyRecveryTime"` // time until energy recovery
	OnePlayCmCount       int64 `json:"onePlayCmCount"`
	OnePlayContinueCount int64 `json:"onePlayContinueCount"` // number of continues allowed
	IsPurchased          int64 `json:"isPurchased"`
}

func DefaultPlayerVarious() PlayerVarious {
	cmSkipCount := int64(0)
	energyRecoveryMax := int64(5)    // max lives should be five
	energyRecoveryTime := int64(600) // ten minutes
	onePlayCmCount := int64(0)
	onePlayContinueCount := int64(5)
	isPurchased := int64(0)
	return PlayerVarious{
		cmSkipCount,
		energyRecoveryMax,
		energyRecoveryTime,
		onePlayCmCount,
		onePlayContinueCount,
		isPurchased,
	}
}
