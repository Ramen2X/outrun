package battle

import (
	"github.com/fluofoxxo/outrun/netobj"
)

func DrawBattleRival(player netobj.Player) netobj.BattleState {
	if !player.BattleState.MatchedUpWithRival { // are we not matched up yet?
		potentialRivals := []int64{}
		// TODO: finish this code
	}
	return player.BattleState
}
