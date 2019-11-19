package conversion

import (
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
	rankingScore := player.PlayerState.HighScore // TODO: this probably differs based on mode...
	rankChanged := int64(2)
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
	league := player.PlayerState.RankingLeague // TODO: This should be changed to QuickRankingLeague when in that mode
	maxScore := player.PlayerState.HighScore
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
