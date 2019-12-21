package conversion

import (
	"github.com/fluofoxxo/outrun/enums"
	"github.com/fluofoxxo/outrun/netobj"
	"github.com/fluofoxxo/outrun/obj"
)

func DebugPlayerToBattleData(player netobj.Player) obj.BattleData {
	uid := player.ID
	username := player.Username
	maxScore := player.BattleState.DailyBattleHighScore
	league := player.PlayerState.RankingLeague
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
	rank := player.PlayerState.Rank
	goOnWin := player.BattleState.WinStreak
	isSentEnergy := int64(0)
	language := int64(enums.LangEnglish)
	return obj.BattleData{
		uid,
		username,
		maxScore,
		league,
		loginTime,
		mainChaoID,
		mainChaoLevel,
		subChaoID,
		subChaoLevel,
		rank,
		mainCharaID,
		mainCharaLevel,
		subCharaID,
		subCharaLevel,
		goOnWin,
		isSentEnergy,
		language,
	}
}
