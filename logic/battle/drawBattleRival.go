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

func DrawBattleRival(player netobj.Player, limit int) netobj.BattleState {
	if !player.BattleState.MatchedUpWithRival { // are we not matched up yet?
		playerIDs := []string{}
		dbaccess.ForEachKey(consts.DBBucketPlayers, func(k, v []byte) error {
			playerIDs = append(playerIDs, string(k))
			return nil
		})
		currentTime := time.Now().UTC().Unix()
		rivalID := ""
		iterations := 0
		for len(playerIDs) > 0 {
			index := rand.Intn(len(playerIDs))
			potentialRival, err := db.GetPlayer(playerIDs[index])
			if err != nil {
				log.Printf("[WARN] (battle.DrawBattleRival) Unable to get player '%s': %s", playerIDs[index], err.Error())
			} else {
				if player.ID != playerIDs[index] && potentialRival.BattleState.ScoreRecordedToday && !potentialRival.BattleState.MatchedUpWithRival && currentTime < potentialRival.BattleState.BattleEndsAt {
					rivalID = playerIDs[index]
					break
				}
			}
			playerIDs[index] = playerIDs[len(playerIDs)-1]
			playerIDs = playerIDs[:len(playerIDs)-1]
			iterations++
			if iterations >= limit && limit > 0 {
				break
			}
		}
		if len(rivalID) > 0 {
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
	if player.ID == player.BattleState.RivalID && player.BattleState.MatchedUpWithRival {
		log.Printf("[WARN] (battle.DrawBattleRival) Somehow matched up with self! Removing match...")
		player.BattleState.MatchedUpWithRival = false
	}
	return player.BattleState
}
