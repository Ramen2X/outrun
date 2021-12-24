package conversion

import (
	"time"

	"github.com/Ramen2X/outrun/enums"
	"github.com/Ramen2X/outrun/netobj"
	"github.com/Ramen2X/outrun/obj"
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
	mainCharaLevel := int64(0)
	subCharaID := player.PlayerState.SubCharaID
	subCharaLevel := int64(0)
	mainChaoID := player.PlayerState.MainChaoID
	mainChaoLevel := int64(0)
	subChaoID := player.PlayerState.SubChaoID
	subChaoLevel := int64(0)
	if player.IndexOfChara(mainCharaID) != -1 {
		mainCharaLevel = player.CharacterState[player.IndexOfChara(mainCharaID)].Level
	}
	if player.IndexOfChara(subCharaID) != -1 {
		subCharaLevel = player.CharacterState[player.IndexOfChara(subCharaID)].Level
	}
	if player.IndexOfChao(mainChaoID) != -1 {
		mainChaoLevel = player.ChaoState[player.IndexOfChao(mainChaoID)].Level
	}
	if player.IndexOfChao(subChaoID) != -1 {
		subChaoLevel = player.ChaoState[player.IndexOfChao(subChaoID)].Level
	}
	language := int64(enums.LangEnglish)
	league := player.PlayerState.RankingLeague
	maxScore := player.PlayerState.HighScore
	if mode == 1 {
		rankingScore = player.PlayerState.TimedHighScore
		league = player.PlayerState.QuickRankingLeague
		maxScore = player.PlayerState.TimedHighScore
	}
	if time.Now().UTC().Unix() > player.PlayerState.WeeklyScoresExpireAt {
		//if expired, show 0 for total scores
		if player.PlayerState.TotalScoreThisPeriod > player.PlayerState.TotalScore {
			player.PlayerState.TotalScore = player.PlayerState.TotalScoreThisPeriod
		}
		if player.PlayerState.TimedTotalScoreThisPeriod > player.PlayerState.TimedTotalScore {
			player.PlayerState.TimedTotalScore = player.PlayerState.TimedTotalScoreThisPeriod
		}
		player.PlayerState.HighScoreThisPeriod = 0
		player.PlayerState.TimedHighScoreThisPeriod = 0
		player.PlayerState.TotalScoreThisPeriod = 0
		player.PlayerState.TimedTotalScoreThisPeriod = 0
	}
	switch lbtype {
	case 0:
		// Friends High Score?
		if mode == 1 {
			rankingScore = player.PlayerState.TimedHighScoreThisPeriod
		} else {
			rankingScore = player.PlayerState.HighScoreThisPeriod
		}
	case 1:
		// Friends Total Score?
		if mode == 1 {
			rankingScore = player.PlayerState.TimedTotalScoreThisPeriod
		} else {
			rankingScore = player.PlayerState.TotalScoreThisPeriod
		}
	case 2:
		// World High Score
		if mode == 1 {
			rankingScore = player.PlayerState.TimedHighScoreThisPeriod
		} else {
			rankingScore = player.PlayerState.HighScoreThisPeriod
		}
	case 3:
		// World Total Score
		if mode == 1 {
			rankingScore = player.PlayerState.TimedTotalScoreThisPeriod
		} else {
			rankingScore = player.PlayerState.TotalScoreThisPeriod
		}
	case 4:
		// Runners' League High Score
		if mode == 1 {
			rankingScore = player.PlayerState.TimedHighScoreThisPeriod
		} else {
			rankingScore = player.PlayerState.HighScoreThisPeriod
		}
	case 5:
		// Runners' League Total Score
		if mode == 1 {
			rankingScore = player.PlayerState.TimedTotalScoreThisPeriod
		} else {
			rankingScore = player.PlayerState.TotalScoreThisPeriod
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
