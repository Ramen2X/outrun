package obj

import (
	"strconv"
	"time"
)

type RaidbossUserState struct {
	WrestleID         string `json:"friendId"`
	Name              string `json:"name"`
	Grade             int64  `json:"grade"`
	NumRank           int64  `json:"numRank"`
	LoginTime         int64  `json:"loginTime"`
	CharacterID       string `json:"charaId"`
	CharacterLevel    int64  `json:"characterLevel"`
	SubcharacterID    string `json:"subCharaId"`
	SubcharacterLevel int64  `json:"subCharaLevel"`
	MainChaoID        string `json:"mainChaoId"`
	MainChaoLevel     int64  `json:"mainChaoLevel"`
	SubChaoID         string `json:"subChaoId"`
	SubChaoLevel      int64  `json:"subChaoLevel"`
	Language          int64  `json:"language"` // enums.Lang*
	League            int64  `json:"league"`
	WrestleCount      int64  `json:"wrestleCount"`
	WrestleDamage     int64  `json:"wrestleDamage"`
	WrestleBeatFlg    int64  `json:"wrestleBeatFlg"`
}

func NewRaidbossUserState(wid, n string, g, nr, lt, cid, cl, schid, schl, mcid, mcl, scid, scl, lang, league, wc, wd, wbf int64) RaidbossUserState {
	return RaidbossUserState{
		wid,
		n,
		g,
		nr,
		lt,
		strconv.Itoa(int(cid)),
		cl,
		strconv.Itoa(int(schid)),
		schl,
		strconv.Itoa(int(mcid)),
		mcl,
		strconv.Itoa(int(scid)),
		scl,
		lang,
		league,
		wc,
		wd,
		wbf,
	}
}

func DefaultRaidbossUserState(uid string) RaidbossUserState {
	return NewRaidbossUserState(
		uid,
		"",
		1,
		0,
		time.Now().Unix(),
		0,
		0,
		0,
		0,
		0,
		0,
		0,
		0,
		1,
		0,
		0,
		0,
		0,
	)
}
