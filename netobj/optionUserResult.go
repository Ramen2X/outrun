package netobj

type OptionUserResult struct {
	TotalSumHighScore      int64 `json:"totalSumHightScore"`     //
	NumTakeAllRings        int64 `json:"numTakeAllRings"`        //
	NumTakeAllRedRings     int64 `json:"numTakeAllRedRings"`     //
	NumChaoRoulette        int64 `json:"numChaoRoulette"`        //
	NumItemRoulette        int64 `json:"numItemRoulette"`        //
	NumJackpot             int64 `json:"numJackPot"`             //
	NumMaximumJackpotRings int64 `json:"numMaximumJackPotRings"` //
	NumSupport             int64 `json:"numSupport"`             //
}

func DefaultOptionUserResult() OptionUserResult {
	totalSumHighScore := int64(0)
	numTakeAllRings := int64(10000)
	numTakeAllRedRings := int64(50)
	numChaoRoulette := int64(0)
	numItemRoulette := int64(0)
	numJackpot := int64(0)
	numMaximumJackpotRings := int64(0)
	numSupport := int64(19191)
	return OptionUserResult{
		totalSumHighScore,
		numTakeAllRings,
		numTakeAllRedRings,
		numChaoRoulette,
		numItemRoulette,
		numJackpot,
		numMaximumJackpotRings,
		numSupport,
	}
}
