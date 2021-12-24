package constnetobjs

import (
	"math/rand"
	"strconv"

	"github.com/Ramen2X/outrun/config/eventconf"
	"github.com/Ramen2X/outrun/config/gameconf"
	"github.com/Ramen2X/outrun/consts"
	"github.com/Ramen2X/outrun/enums"
	"github.com/Ramen2X/outrun/netobj"
	"github.com/Ramen2X/outrun/obj"
)

var BlankPlayer = func() netobj.Player {
	randChar := func(charset string, length int64) string {
		runes := []rune(charset)
		final := make([]rune, 10)
		for i := range final {
			final[i] = runes[rand.Intn(len(runes))]
		}
		return string(final)
	}
	// create ID
	uid := ""
	for i := range make([]byte, 10) {
		if i == 0 { // if first character
			uid += strconv.Itoa(rand.Intn(9) + 1)
		} else {
			uid += strconv.Itoa(rand.Intn(10))
		}
	}
	username := ""
	password := randChar("abcdefghijklmnopqrstuvwxyz1234567890", 10)
	migrationPassword := randChar("abcdefghijklmnopqrstuvwxyz1234567890", 10)
	userPassword := ""
	key := randChar("abcdefghijklmnopqrstuvwxyz1234567890", 10)
	playerState := netobj.DefaultPlayerState()
	characterState := netobj.DefaultCharacterState()
	chaoState := GetAllNetChaoList() // needed for chaoRouletteAllowed
	eventState := netobj.DefaultEventState()
	eventUserRaidbossState := netobj.DefaultUserRaidbossState()
	optionUserResult := netobj.DefaultOptionUserResult()
	mileageMapState := netobj.DefaultMileageMapState()
	mileageFriends := []netobj.MileageFriend{}
	playerVarious := netobj.DefaultPlayerVarious()
	rouletteInfo := netobj.DefaultRouletteInfo()
	wheelOptions := netobj.DefaultWheelOptions(playerState.NumRouletteTicket, rouletteInfo.RouletteCountInPeriod, enums.WheelRankNormal, consts.RouletteFreeSpins)
	// TODO: get rid of logic here?
	allowedCharacters := []string{}
	allowedChao := []string{}
	for _, chao := range chaoState {
		if chao.Level < 10 { // not max level
			allowedChao = append(allowedChao, chao.ID)
		}
	}
	for _, character := range characterState {
		if character.Star < 10 { // not max star
			allowedCharacters = append(allowedCharacters, character.ID)
		}
	}
	chaoRouletteGroup := netobj.DefaultChaoRouletteGroup(playerState, allowedCharacters, allowedChao, true)
	personalEvents := []eventconf.ConfiguredEvent{}
	suspended := false
	operatorMessages := []obj.OperatorMessage{
		obj.NewOperatorMessage(
			1,
			"A welcome gift", // TODO: make this configurable
			obj.NewMessageItem(
				enums.ItemIDRedRing,
				gameconf.CFile.StartingRedRings,
				0,
				0,
			),
			2592000,
		),
		obj.NewOperatorMessage(
			2,
			"A welcome gift", // TODO: make this configurable
			obj.NewMessageItem(
				enums.ItemIDRing,
				gameconf.CFile.StartingRings,
				0,
				0,
			),
			2592000,
		),
		obj.NewOperatorMessage(
			3,
			"A welcome gift", // TODO: make this configurable
			obj.NewMessageItem(
				enums.IDRouletteTicketPremium,
				5,
				0,
				0,
			),
			2592000,
		),
		obj.NewOperatorMessage(
			4,
			"A welcome gift", // TODO: make this configurable
			obj.NewMessageItem(
				enums.IDRouletteTicketItem,
				5,
				0,
				0,
			),
			2592000,
		),
	}
	battleState := netobj.DefaultBattleState()
	loginBonusState := netobj.DefaultLoginBonusState(0)
	return netobj.NewPlayer(
		uid,
		username,
		password,
		migrationPassword,
		userPassword,
		key,
		playerState,
		characterState,
		chaoState,
		eventState,
		eventUserRaidbossState,
		optionUserResult,
		mileageMapState,
		mileageFriends,
		playerVarious,
		wheelOptions,
		rouletteInfo,
		chaoRouletteGroup,
		personalEvents,
		suspended,
		operatorMessages,
		battleState,
		loginBonusState,
	)
}() // TODO: Solve duplication requirement with db/assistants.go
