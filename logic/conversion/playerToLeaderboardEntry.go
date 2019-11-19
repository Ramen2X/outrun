package conversion

import (
	"time"

	"github.com/fluofoxxo/outrun/enums"
	"github.com/fluofoxxo/outrun/netobj"
	"github.com/fluofoxxo/outrun/obj"
	"github.com/jinzhu/now"
)

func PlayerToLeaderboardEntry(player netobj.Player, place, mode, lbtype int64) obj.LeaderboardEntry {
	friendID := player.ID
	name := player.Username
	url := player.Username + "_findme" // TODO: only used for testing right now
	grade := place
	exposeOnline := int64(0)
	rankingScore := player.PlayerState.HighScore
	rankChanged := int64(0)
	isSentEnergy := int64(0)
	expireTime := now.EndOfWeek().UTC().Unix()
	numRank := player.PlayerState.Rank
	loginTime := player.LastLogin
	mainCharaID := player.PlayerState.MainCharaID
	mainCharaLevel := player.CharacterState[0].Level // TODO: Is this right?
	subCharaID := player.PlayerState.SubCharaID
	subCharaLevel := player.CharacterState[1].Level
	mainChaoID := player.PlayerState.MainChaoID
	mainChaoLevel := player.ChaoState[0].Level
	subChaoID := player.PlayerState.SubChaoID
	subChaoLevel := player.ChaoState[1].Level
	language := int64(enums.LangEnglish)
	league := player.PlayerState.RankingLeague
	maxScore := player.PlayerState.HighScore
	if mode == 1 {
		rankingScore = player.PlayerState.TimedHighScore
		league = player.PlayerState.QuickRankingLeague
		maxScore = player.PlayerState.TimedHighScore
	}
	if time.Now().UTC().Unix() > player.PlayerState.TotalScoreExpiresAt {
		//if expired, show 0 for total scores
		player.PlayerState.TotalScore = 0
		player.PlayerState.TimedTotalScore = 0
	}
	switch lbtype {
	case 0:
		// Friends High Score?
	case 1:
		// Friends Total Score?
		if mode == 1 {
			rankingScore = player.PlayerState.TimedTotalScore
		} else {
			rankingScore = player.PlayerState.TotalScore
		}
	case 2:
		// World High Score
	case 3:
		// World Total Score
		if mode == 1 {
			rankingScore = player.PlayerState.TimedTotalScore
		} else {
			rankingScore = player.PlayerState.TotalScore
		}
	case 4:
		// Runners' League High Score
	case 5:
		// Runners' League Total Score
		if mode == 1 {
			rankingScore = player.PlayerState.TimedTotalScore
		} else {
			rankingScore = player.PlayerState.TotalScore
		}
	case 6:
		// History High Score
	case 7:
		// History Total Score
		if mode == 1 {
			rankingScore = player.PlayerState.TimedTotalScore
		} else {
			rankingScore = player.PlayerState.TotalScore
		}
	case 8:
		// Event High Score?
	case 9:
		// Event Total Score?
		rankingScore = player.EventState.Param //TODO: is this right?
	}
	return obj.LeaderboardEntry{
		friendID,
		name,
		url,
		grade,
		exposeOnline,
		rankingScore,
		rankChanged,
		isSentEnergy,
		expireTime,
		numRank,
		loginTime,
		mainCharaID,
		mainCharaLevel,
		subCharaID,
		subCharaLevel,
		mainChaoID,
		mainChaoLevel,
		subChaoID,
		subChaoLevel,
		language,
		league,
		maxScore,
	}
}
