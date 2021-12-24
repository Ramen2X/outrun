package netobj

import (
	"fmt"
	"strconv"
	"time"

	"github.com/Ramen2X/outrun/config/eventconf"
	"github.com/Ramen2X/outrun/enums"
	"github.com/Ramen2X/outrun/obj"
	"github.com/Ramen2X/outrun/obj/constobjs"
)

type Player struct {
	ID                     string                      `json:"userID"`
	Username               string                      `json:"username"`
	Password               string                      `json:"password"`
	MigrationPassword      string                      `json:"migrationPassword"` // used in migration
	UserPassword           string                      `json:"userPassword"`      // used in migration
	Key                    string                      `json:"key"`
	LastLogin              int64                       // TODO: use `json:"lastLogin"`
	PlayerState            PlayerState                 `json:"playerState"`
	CharacterState         []Character                 `json:"characterState"`
	ChaoState              []Chao                      `json:"chaoState"`
	EventState             EventState                  `json:"eventState"`
	EventUserRaidbossState EventUserRaidbossState      `json:"eventUserRaidboss"`
	OptionUserResult       OptionUserResult            `json:"optionUserResult"`
	MileageMapState        MileageMapState             `json:"mileageMapState"`
	MileageFriends         []MileageFriend             `json:"mileageFriendList"`
	PlayerVarious          PlayerVarious               `json:"playerVarious"`
	LastWheelOptions       WheelOptions                `json:"ORN_wheelOptions"` // TODO: Make RouletteGroup to hold LastWheelOptions and RouletteInfo?
	RouletteInfo           RouletteInfo                `json:"ORN_rouletteInfo"`
	ChaoRouletteGroup      ChaoRouletteGroup           `json:"ORN_chaoRouletteGroup"`
	PersonalEvents         []eventconf.ConfiguredEvent `json:"ORN_personalEvents"`
	Suspended              bool                        `json:"ORN_suspended"`
	OperatorMessages       []obj.OperatorMessage       `json:"operatorMessageList"`
	BattleState            BattleState                 `json:"battleState"`
	LoginBonusState        LoginBonusState             `json:"loginBonusState"`
}

func NewPlayer(id, username, password, migrationPassword, userPassword, key string, playerState PlayerState, characterState []Character, chaoState []Chao, eventState EventState, eventUserRaidbossState EventUserRaidbossState, optionUserResult OptionUserResult, mileageMapState MileageMapState, mf []MileageFriend, playerVarious PlayerVarious, wheelOptions WheelOptions, rouletteInfo RouletteInfo, chaoRouletteGroup ChaoRouletteGroup, personalEvents []eventconf.ConfiguredEvent, suspended bool, operatorMessages []obj.OperatorMessage, battleState BattleState, loginBonusState LoginBonusState) Player {
	return Player{
		id,
		username,
		password,
		migrationPassword,
		userPassword,
		key,
		time.Now().Unix(),
		playerState,
		characterState,
		chaoState,
		eventState,
		eventUserRaidbossState,
		optionUserResult,
		mileageMapState,
		mf,
		playerVarious,
		wheelOptions,
		rouletteInfo,
		chaoRouletteGroup,
		personalEvents,
		suspended,
		operatorMessages,
		battleState,
		loginBonusState,
	}
}

/*
func (p *Player) Save() {

}
*/
// TODO: remove any functions that access p.PlayerState since we are not calling from a pointer anyways and it will not modify the object
func (p *Player) AddRings(amount int64) {
	ps := p.PlayerState
	ps.NumRings += amount
}
func (p *Player) SubRings(amount int64) {
	ps := p.PlayerState
	ps.NumRings -= amount
}
func (p *Player) AddRedRings(amount int64) {
	ps := p.PlayerState
	ps.NumRedRings += amount
}
func (p *Player) SubRedRings(amount int64) {
	ps := p.PlayerState
	ps.NumRedRings -= amount
}
func (p *Player) SetUsername(username string) {
	p.Username = username
}
func (p *Player) SetPassword(password string) {
	p.Password = password
}
func (p *Player) AddEnergy(amount int64) {
	ps := p.PlayerState
	ps.Energy += amount
}
func (p *Player) SubEnergy(amount int64) {
	ps := p.PlayerState
	ps.Energy -= amount
}
func (p *Player) SetMainCharacter(cid string) {
	ps := p.PlayerState
	ps.MainCharaID = cid
}
func (p *Player) SetSubCharacter(cid string) {
	ps := p.PlayerState
	ps.SubCharaID = cid
}
func (p *Player) SetMainChao(chid string) {
	ps := p.PlayerState
	ps.MainChaoID = chid
}
func (p *Player) SetSubChao(chid string) {
	ps := p.PlayerState
	ps.SubChaoID = chid
}
func (p *Player) AddItem(item obj.Item) {
	ps := p.PlayerState
	ps.Items = append(ps.Items, item)
}
func (p *Player) RemoveItemOf(iid string) bool {
	newItems := []obj.Item{}
	foundItem := false
	ps := p.PlayerState
	for _, item := range ps.Items {
		if item.ID != iid || foundItem {
			newItems = append(newItems, item)
		} else if !foundItem {
			foundItem = true
		}
	}
	ps.Items = newItems
	return foundItem
}
func (p *Player) IndexOfItem(iid string) int {
	for i, item := range p.PlayerState.Items {
		if item.ID == iid {
			return i
		}
	}
	return -1
}
func (p *Player) RemoveAllItemsOf(iid string) {
	for p.RemoveItemOf(iid) {
	}
}
func (p *Player) AddAnimals(amount int64) {
	ps := p.PlayerState
	ps.Animals += amount
}
func (p *Player) SubAnimals(amount int64) {
	ps := p.PlayerState
	ps.Animals -= amount
}
func (p *Player) ApplyHighScore(score int64) bool {
	ps := p.PlayerState
	if ps.HighScore < score {
		ps.HighScore = score
		return true
	}
	return false
}
func (p *Player) AddDistance(amount int64) {
	ps := p.PlayerState
	ps.TotalDistance += amount
	p.ApplyHighDistance(amount)
}
func (p *Player) ApplyHighDistance(amount int64) {
	ps := p.PlayerState
	ps.HighDistance = amount
}
func (p *Player) AddNewChaoByID(chid string) bool {
	chao := constobjs.Chao[chid]
	netchao := NewNetChao(
		chao,
		enums.ChaoStatusOwned, // TODO: does the idea that a chao is owned mean that it's possible to send chao that are not owned?
		1,
		enums.ChaoDealingNone,
		1, // implies that adding means acquired. This may not be the case if we can send non-owned chao.
	)
	return p.AddNetChao(netchao)
}
func (p *Player) AddNewChao(chao obj.Chao) bool {
	netchao := NewNetChao(
		chao,
		enums.ChaoStatusOwned,
		1,
		enums.ChaoDealingNone,
		1,
	)
	return p.AddNetChao(netchao)
}
func (p *Player) AddNetChao(netchao Chao) bool {
	// Returns whether or not the Chao was already found.
	// It will not add Chao already in the ChaoState.
	if !p.HasChao(netchao.Chao.ID) {
		p.ChaoState = append(p.ChaoState, netchao)
		return false
	}
	return true
}
func (p *Player) HasChao(chid string) bool {
	for _, netchao := range p.ChaoState {
		if netchao.Chao.ID == chid {
			return true
		}
	}
	return false
}
func (p *Player) GetChara(cid string) (Character, error) {
	var char Character
	found := false
	for _, c := range p.CharacterState {
		if c.ID == cid {
			char = c
			found = true
		}
	}
	if !found {
		return char, fmt.Errorf("character not found")
	}
	return char, nil
}
func (p *Player) IndexOfChara(cid string) int {
	for i, char := range p.CharacterState {
		if char.ID == cid {
			return i
		}
	}
	return -1
}
func (p *Player) GetChao(chid string) (Chao, error) {
	var chao Chao
	found := false
	for _, c := range p.ChaoState {
		if c.ID == chid {
			chao = c
			found = true
		}
	}
	if !found {
		return chao, fmt.Errorf("chao not found")
	}
	return chao, nil
}
func (p *Player) IndexOfChao(chid string) int {
	for i, chao := range p.ChaoState {
		if chao.ID == chid {
			return i
		}
	}
	return -1
}
func (p *Player) GetMainChara() (Character, error) {
	ps := p.PlayerState
	cid := ps.MainCharaID
	char, err := p.GetChara(cid)
	return char, err
}
func (p *Player) GetSubChara() (Character, error) {
	ps := p.PlayerState
	cid := ps.SubCharaID
	char, err := p.GetChara(cid)
	return char, err
}
func (p *Player) GetMainChao() (Chao, error) {
	ps := p.PlayerState
	chid := ps.MainChaoID
	chao, err := p.GetChao(chid)
	return chao, err
}
func (p *Player) GetSubChao() (Chao, error) {
	ps := p.PlayerState
	chid := ps.SubChaoID
	chao, err := p.GetChao(chid)
	return chao, err
}
func (p *Player) GetMaxLevelChao(legacyCap bool) []Chao {
	maxLevel := int64(10)
	if legacyCap {
		maxLevel = int64(5)
	}
	mxlvl := []Chao{}
	for _, c := range p.ChaoState {
		if c.Level >= maxLevel { // if max level (or above)
			mxlvl = append(mxlvl, c)
		}
	}
	return mxlvl
}
func (p *Player) GetMaxLevelChaoIDs(legacyCap bool) []string {
	maxLevel := int64(10)
	if legacyCap {
		maxLevel = int64(5)
	}
	chao := p.GetMaxLevelChao(legacyCap)
	ids := []string{}
	for _, c := range chao {
		if c.Level >= maxLevel { // if max level (or above)
			ids = append(ids, c.ID)
		}
	}
	return ids
}
func (p *Player) GetMaxLevelCharacters() []Character {
	mxlvl := []Character{}
	for _, ch := range p.CharacterState {
		if ch.Star >= 10 { // if max stars (or above)
			mxlvl = append(mxlvl, ch)
		}
	}
	return mxlvl
}
func (p *Player) GetMaxLevelCharacterIDs() []string {
	chars := p.GetMaxLevelCharacters()
	ids := []string{}
	for _, ch := range chars {
		//if ch.Level >= 100 { // if max level (or above)
		if ch.Star >= 10 {
			ids = append(ids, ch.ID)
		}
	}
	return ids
}
func (p *Player) GetAllMaxLevelIDs(legacyCap bool) []string {
	chars := p.GetMaxLevelCharacterIDs()
	chao := p.GetMaxLevelChaoIDs(legacyCap)
	combined := []string{}
	combined = append(combined, chars...)
	combined = append(combined, chao...)
	return combined
}
func (p *Player) AllChaoMaxLevel(legacyCap bool) bool {
	maxLevel := int64(10)
	if legacyCap {
		maxLevel = int64(5)
	}
	for _, c := range p.ChaoState {
		if c.Level < maxLevel {
			return false
		}
	}
	return true
}
func (p *Player) AllCharactersMaxLevel() bool {
	for _, ch := range p.CharacterState {
		if ch.Star < 10 {
			return false
		}
	}
	return true
}
func (p *Player) GetAllNonMaxedChaoAndCharacters(legacyCap bool) []string {
	result := append(p.GetAllNonMaxedChao(legacyCap), p.GetAllNonMaxedCharacters()...) // combine two
	return result
}
func (p *Player) GetAllNonMaxedChao(legacyCap bool) []string {
	result := []string{}
	maxLevel := int64(10)
	if legacyCap {
		maxLevel = int64(5)
	}
	for _, chao := range p.ChaoState {
		if chao.Level < maxLevel { // not max level
			result = append(result, chao.ID)
		}
	}
	return result
}
func (p *Player) GetAllNonMaxedCharacters() []string {
	result := []string{}
	for _, character := range p.CharacterState {
		if character.Star < 10 { // not max star
			result = append(result, character.ID)
		}
	}
	return result
}

func (p *Player) AcceptOperatorMessage(id int64) interface{} {
	for index, message := range p.OperatorMessages {
		if strconv.Itoa(int(id)) == message.ID {
			p.RemoveFromOperatorMessages(index)
			if time.Now().UTC().Unix() < message.ExpireTime {
				return obj.MessageItemToPresent(message.Item)
			}
		}
	}
	return nil
}

func (p *Player) RemoveFromOperatorMessages(index int) {
	// don't care about order; it shouldn't really matter
	p.OperatorMessages[index] = p.OperatorMessages[len(p.OperatorMessages)-1]
	p.OperatorMessages = p.OperatorMessages[:len(p.OperatorMessages)-1]
}

func (p *Player) GetAllOperatorMessageIDs() []int64 {
	result := []int64{}
	for _, message := range p.OperatorMessages {
		messageid, _ := strconv.Atoi(message.ID)
		result = append(result, int64(messageid))
	}
	return result
}

func (p *Player) CleanUpExpiredOperatorMessages() {
	removals := -1
	for removals != 0 {
		removals = 0
		for index, message := range p.OperatorMessages {
			if time.Now().UTC().Unix() >= message.ExpireTime {
				p.RemoveFromOperatorMessages(index)
				removals++
			}
		}
	}
}

func (p *Player) AddOperatorMessage(messageContents string, item obj.MessageItem, expiresAfter int64) {
	// A function to add an operator message, automatically determining its ID
	// TODO: Optimize this code. It's not pretty.
	index := 0
	foundPreferredID := false
	preferredID := 1
	for !foundPreferredID {
		foundPreferredID = true
		index = 0
		for index < len(p.OperatorMessages) {
			if p.OperatorMessages[index].ID == strconv.Itoa(preferredID) {
				foundPreferredID = false
			}
			index++
		}
		preferredID++
	}
	preferredID--
	p.OperatorMessages = append(
		p.OperatorMessages,
		obj.NewOperatorMessage(
			int64(preferredID),
			messageContents,
			item,
			expiresAfter,
		),
	)
}
