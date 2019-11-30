package conversion

import (
	"github.com/fluofoxxo/outrun/enums"
	"github.com/fluofoxxo/outrun/netobj"
	"github.com/fluofoxxo/outrun/obj"
)

func PlayerToRaidbossUserState(player netobj.Player, place int64) obj.RaidbossUserState {
	wrestleID := player.ID
	name := player.Username
	grade := place
	numRank := player.PlayerState.Rank
	loginTime := player.LastLogin
	mainCharaID := player.PlayerState.MainCharaID
	mainCharaLevel := player.CharacterState[0].Level
	subCharaID := player.PlayerState.SubCharaID
	subCharaLevel := player.CharacterState[1].Level
	mainChaoID := player.PlayerState.MainChaoID
	mainChaoLevel := player.ChaoState[0].Level
	subChaoID := player.PlayerState.SubChaoID
	subChaoLevel := player.ChaoState[1].Level
	language := int64(enums.LangEnglish)
	league := player.PlayerState.RankingLeague
	wrestleCount := player.EventUserRaidbossState.NumRaidBossEncountered // TODO: is this right?
	wrestleDamage := int64(0)
	wrestleBeatFlg := int64(0)
	return obj.RaidbossUserState{
		wrestleID,
		name,
		grade,
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
		wrestleCount,
		wrestleDamage,
		wrestleBeatFlg,
	}
}
