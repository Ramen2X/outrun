package netobj

type OptionUserResult struct {
	TotalSumHighScore      int64 `json:"totalSumHightScore"`      // highest total score recorded
	QuickTotalSumHighScore int64 `json:"quickTotalSumHightScore"` // same as above but for timed mode
	NumTakeAllRings        int64 `json:"numTakeAllRings"`         // total number of rings acquired ever
	NumTakeAllRedRings     int64 `json:"numTakeAllRedRings"`      // total number of red rings acquired ever
	NumChaoRoulette        int64 `json:"numChaoRoulette"`         // total times the chao roulette was spun
	NumItemRoulette        int64 `json:"numItemRoulette"`         // total times the item roulette was spun
	NumJackpot             int64 `json:"numJackPot"`              // total number of jackpots won ever
	NumMaximumJackpotRings int64 `json:"numMaximumJackPotRings"`  // biggest jackpot won
	NumSupport             int64 `json:"numSupport"`              // ?
}

func DefaultOptionUserResult() OptionUserResult {
	totalSumHighScore := int64(0)
	quickTotalSumHighScore := int64(0)
	numTakeAllRings := int64(0) // TODO: should the starting rings and red rings count?
	numTakeAllRedRings := int64(0)
	numChaoRoulette := int64(0)
	numItemRoulette := int64(0)
	numJackpot := int64(0)
	numMaximumJackpotRings := int64(0)
	numSupport := int64(19191)
	return OptionUserResult{
		totalSumHighScore,
		quickTotalSumHighScore,
		numTakeAllRings,
		numTakeAllRedRings,
		numChaoRoulette,
		numItemRoulette,
		numJackpot,
		numMaximumJackpotRings,
		numSupport,
	}
}
