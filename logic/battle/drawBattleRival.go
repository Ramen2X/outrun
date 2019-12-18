package battle

import (
	"log"
	"math/rand"
	"time"

	"github.com/fluofoxxo/outrun/consts"
	"github.com/fluofoxxo/outrun/db"
	"github.com/fluofoxxo/outrun/db/dbaccess"
	"github.com/fluofoxxo/outrun/netobj"
)

func DrawBattleRival(player netobj.Player) netobj.BattleState {
	if !player.BattleState.MatchedUpWithRival { // are we not matched up yet?
		playerIDs := []string{}
		dbaccess.ForEachKey(consts.DBBucketPlayers, func(k, v []byte) error {
			playerIDs = append(playerIDs, string(k))
			return nil
		})
		potentialRivalIDs := []string{}
		for _, pid := range playerIDs {
			potentialRival, err := db.GetPlayer(pid)
			if err != nil {
				log.Printf("[WARN] (battle.DrawBattleRival) Unable to get player '%s': %s", pid, err.Error())
			} else {
				if player.ID != pid && potentialRival.BattleState.ScoreRecordedToday && !potentialRival.BattleState.MatchedUpWithRival && time.Now().UTC().Unix() < potentialRival.BattleState.BattleEndsAt {
					potentialRivalIDs = append(potentialRivalIDs, potentialRival.ID)
				}
			}
		}
		if len(potentialRivalIDs) > 0 {
			rivalID := potentialRivalIDs[rand.Intn(len(potentialRivalIDs))]
			rival, err := db.GetPlayer(rivalID)
			if err != nil {
				log.Printf("[WARN] (battle.DrawBattleRival) Unable to get player '%s': %s", rivalID, err.Error())
			} else {
				rival.BattleState.RivalID = player.ID
				rival.BattleState.MatchedUpWithRival = true
				err = db.SavePlayer(rival)
				if err != nil {
					log.Printf("[WARN] (battle.DrawBattleRival) Unable to save rival data: %s", err.Error())
				} else {
					player.BattleState.RivalID = rivalID
					player.BattleState.MatchedUpWithRival = true
				}
			}
		}
	}
	return player.BattleState
}
