package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Ramen2X/outrun/bgtasks"
	"github.com/Ramen2X/outrun/config"
	"github.com/Ramen2X/outrun/config/campaignconf"
	"github.com/Ramen2X/outrun/config/eventconf"
	"github.com/Ramen2X/outrun/config/gameconf"
	"github.com/Ramen2X/outrun/config/infoconf"
	"github.com/Ramen2X/outrun/cryption"
	"github.com/Ramen2X/outrun/inforeporters"
	"github.com/Ramen2X/outrun/meta"
	"github.com/Ramen2X/outrun/muxhandlers"
	"github.com/Ramen2X/outrun/muxhandlers/muxobj"
	"github.com/Ramen2X/outrun/orpc"
	"github.com/gorilla/mux"
)

const UNKNOWN_REQUEST_DIRECTORY = "logging/unknown_requests/"

var (
	LogExecutionTime   = true
	LogUnknownRequests = false
)

func HandleUnknownRequest(w http.ResponseWriter, r *http.Request) {
	recv, _ := cryption.GetReceivedMessage(r)
	if LogUnknownRequests {
		// make a new logging path
		timeStr := strconv.Itoa(int(time.Now().Unix()))
		os.MkdirAll(UNKNOWN_REQUEST_DIRECTORY, 0644)
		normalizedReq := strings.ReplaceAll(r.URL.Path, "/", "-")
		path := UNKNOWN_REQUEST_DIRECTORY + normalizedReq + "_" + timeStr + ".txt"
		err := ioutil.WriteFile(path, recv, 0644)
		if err != nil {
			log.Println("[OUT] UNABLE TO WRITE UNKNOWN REQUEST: " + err.Error())
			w.Write([]byte(""))
			return
		}
		log.Println("[OUT] !!!!!!!!!!!! Unknown request, output to " + path)
	}
	w.Write([]byte(""))
}

func HandleDSFile(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(""))
}

func HandlePPAdsFile(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(""))
}

func removePrependingSlashes(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for len(r.URL.Path) != 0 && string(r.URL.Path[0]) == "/" {
			r.URL.Path = r.URL.Path[1:]
		}
		r.URL.Path = "/" + r.URL.Path
		next.ServeHTTP(w, r)
	})
}

func checkArgs() bool {
	// TODO: _VERY_ dirty command line argument checking. This should be
	// changed into something more robust and less hacky!
	args := os.Args[1:] // drop executable
	amt := len(args)
	if amt >= 1 {
		if args[0] == "--version" {
			fmt.Printf("Outrun %s\n", meta.Version)
			return true
		}
		fmt.Println("Unknown given arguments")
		return true
	}
	return false
}

func main() {
	end := checkArgs()
	if end {
		return
	}
	rand.Seed(time.Now().UTC().UnixNano())

	err := config.Parse("config.json")
	if err != nil {
		log.Printf("[INFO] Failure loading config file config.json (%s), using defaults\n", err)
	} else {
		log.Println("[INFO] Config file (config.json) loaded")
	}

	err = eventconf.Parse(config.CFile.EventConfigFilename)
	if err != nil {
		if !config.CFile.SilenceEventConfigErrors {
			log.Printf("[INFO] Failure loading event config file %s (%s), using defaults\n", config.CFile.EventConfigFilename, err)
		}
	} else {
		log.Printf("[INFO] Event config file (%s) loaded\n", config.CFile.EventConfigFilename)
	}

	err = infoconf.Parse(config.CFile.InfoConfigFilename)
	if err != nil {
		if !config.CFile.SilenceInfoConfigErrors {
			log.Printf("[INFO] Failure loading info config file %s (%s), using defaults\n", config.CFile.InfoConfigFilename, err)
		}
	} else {
		log.Printf("[INFO] Info config file (%s) loaded\n", config.CFile.InfoConfigFilename)
	}

	err = gameconf.Parse(config.CFile.GameConfigFilename)
	if err != nil {
		if !config.CFile.SilenceGameConfigErrors {
			log.Printf("[INFO] Failure loading game config file %s (%s), using defaults\n", config.CFile.GameConfigFilename, err)
		}
	} else {
		log.Printf("[INFO] Game config file (%s) loaded\n", config.CFile.GameConfigFilename)
	}

	err = campaignconf.Parse(config.CFile.CampaignConfigFilename)
	if err != nil {
		if !config.CFile.SilenceCampaignConfigErrors {
			log.Printf("[INFO] Failure loading campaign config file %s (%s), using defaults\n", config.CFile.CampaignConfigFilename, err)
		}
	} else {
		log.Printf("[INFO] Campaign config file (%s) loaded\n", config.CFile.CampaignConfigFilename)
	}

	if config.CFile.EnableRPC {
		orpc.Start()
	}

	h := muxobj.Handle
	router := mux.NewRouter()
	router.StrictSlash(true)
	LogExecutionTime = config.CFile.DoTimeLogging
	LogUnknownRequests = config.CFile.LogUnknownRequests
	prefix := config.CFile.EndpointPrefix
	// == Login ==
	router.HandleFunc(prefix+"/Login/login/", h(muxhandlers.Login, LogExecutionTime))
	//router.HandleFunc(prefix+"/Login/reLogin/", h(muxhandlers.ReLogin, LogExecutionTime))
	router.HandleFunc(prefix+"/Login/getVariousParameter/", h(muxhandlers.GetVariousParameter, LogExecutionTime))
	router.HandleFunc(prefix+"/Game/getDailyChalData/", h(muxhandlers.GetDailyChallengeData, LogExecutionTime))
	router.HandleFunc(prefix+"/Game/getCostList/", h(muxhandlers.GetCostList, LogExecutionTime))
	router.HandleFunc(prefix+"/Game/getCampaignList/", h(muxhandlers.GetCampaignList, LogExecutionTime))
	router.HandleFunc(prefix+"/login/getInformation/", h(muxhandlers.GetInformation, LogExecutionTime))
	router.HandleFunc(prefix+"/login/getTicker/", h(muxhandlers.GetTicker, LogExecutionTime))
	router.HandleFunc(prefix+"/Login/getCountry/", h(muxhandlers.GetCountry, LogExecutionTime))

	// == Login Bonus ==
	router.HandleFunc(prefix+"/Login/loginBonus/", h(muxhandlers.LoginBonus, LogExecutionTime))
	router.HandleFunc(prefix+"/Login/loginBonusSelect/", h(muxhandlers.LoginBonusSelect, LogExecutionTime))

	// == Events ==
	router.HandleFunc(prefix+"/Event/getEventList/", h(muxhandlers.GetEventList, LogExecutionTime))
	router.HandleFunc(prefix+"/Event/getEventReward/", h(muxhandlers.GetEventReward, LogExecutionTime))
	router.HandleFunc(prefix+"/Event/getEventState/", h(muxhandlers.GetEventState, LogExecutionTime))

	// == Raid boss ==
	router.HandleFunc(prefix+"/Event/getEventUserRaidboss/", h(muxhandlers.GetEventUserRaidbossState, LogExecutionTime))
	router.HandleFunc(prefix+"/Event/getEventUserRaidbossList/", h(muxhandlers.GetEventUserRaidbossList, LogExecutionTime))
	//router.HandleFunc(prefix+"/Event/getEventRaidbossDesiredList/", h(muxhandlers.GetEventRaidbossDesiredList, LogExecutionTime))
	//router.HandleFunc(prefix+"/Event/getEventRaidbossUserList/", h(muxhandlers.GetEventRaidbossUserList, LogExecutionTime))
	router.HandleFunc(prefix+"/Event/eventActStart/", h(muxhandlers.EventActStart, LogExecutionTime))
	router.HandleFunc(prefix+"/Event/eventPostGameResults/", h(muxhandlers.EventPostGameResults, LogExecutionTime))
	router.HandleFunc(prefix+"/Event/eventUpdateGameResults/", h(muxhandlers.EventUpdateGameResults, LogExecutionTime))
	router.HandleFunc(prefix+"/Game/drawRaidboss/", h(muxhandlers.DrawRaidBoss, LogExecutionTime))

	// == Player data ==
	router.HandleFunc(prefix+"/Player/getPlayerState/", h(muxhandlers.GetPlayerState, LogExecutionTime))
	router.HandleFunc(prefix+"/Player/getCharacterState/", h(muxhandlers.GetCharacterState, LogExecutionTime))
	router.HandleFunc(prefix+"/Player/getChaoState/", h(muxhandlers.GetChaoState, LogExecutionTime))
	router.HandleFunc(prefix+"/Option/userResult/", h(muxhandlers.GetOptionUserResult, LogExecutionTime))
	router.HandleFunc(prefix+"/Player/setUserName/", h(muxhandlers.SetUsername, LogExecutionTime))
	//router.HandleFunc(prefix+"/Player/setBirthday/", h(muxhandlers.SetBirthday, LogExecutionTime))
	router.HandleFunc(prefix+"/Character/changeCharacter/", h(muxhandlers.ChangeCharacter, LogExecutionTime))
	router.HandleFunc(prefix+"/Chao/equipChao/", h(muxhandlers.EquipChao, LogExecutionTime))
	//router.HandleFunc(prefix+"/Character/useSubCharacter/", h(muxhandlers.UseSubCharacter, LogExecutionTime))
	//router.HandleFunc(prefix+"/Game/getMenuData/", h(muxhandlers.GetMenuData, LogExecutionTime))

	// == Timed mode ==
	router.HandleFunc(prefix+"/Game/quickActStart/", h(muxhandlers.QuickActStart, LogExecutionTime))
	router.HandleFunc(prefix+"/Game/quickPostGameResults/", h(muxhandlers.QuickPostGameResults, LogExecutionTime))

	// == Story mode ==
	router.HandleFunc(prefix+"/Game/actStart/", h(muxhandlers.ActStart, LogExecutionTime))
	router.HandleFunc(prefix+"/Game/getMileageData/", h(muxhandlers.GetMileageData, LogExecutionTime))
	router.HandleFunc(prefix+"/Game/getMileageReward/", h(muxhandlers.GetMileageReward, LogExecutionTime))
	router.HandleFunc(prefix+"/Game/postGameResults/", h(muxhandlers.PostGameResults, LogExecutionTime))

	// == Continue ==
	router.HandleFunc(prefix+"/Game/actRetry/", h(muxhandlers.ActRetry, LogExecutionTime))
	router.HandleFunc(prefix+"/Game/actRetryFree/", h(muxhandlers.ActRetryFree, LogExecutionTime))

	// == Gameplay ==
	router.HandleFunc(prefix+"/Game/getFreeItemList/", h(muxhandlers.GetFreeItemList, LogExecutionTime))

	// == Shop ==
	router.HandleFunc(prefix+"/Store/getRedstarExchangeList/", h(muxhandlers.GetRedStarExchangeList, LogExecutionTime))
	router.HandleFunc(prefix+"/Store/redstarExchange/", h(muxhandlers.RedStarExchange, LogExecutionTime))
	//router.HandleFunc(prefix+"/Store/getRingExchangeList/", h(muxhandlers.GetRingExchangeList, LogExecutionTime))
	//router.HandleFunc(prefix+"/Store/ringExchange/", h(muxhandlers.RingExchange, LogExecutionTime))
	//router.HandleFunc(prefix+"/Store/buyAndroid/", h(muxhandlers.IAPConfirmAndroid, LogExecutionTime))
	//router.HandleFunc(prefix+"/Store/buyIos/", h(muxhandlers.IAPConfirmApple, LogExecutionTime))
	//router.HandleFunc(prefix+"/Store/preparePurchase/", h(muxhandlers.PreparePurchase, LogExecutionTime))
	//router.HandleFunc(prefix+"/Store/purchase/", h(muxhandlers.Purchase, LogExecutionTime))

	// == Facebook Integration/Friends (Required for iOS?) ==
	router.HandleFunc(prefix+"/Friend/getFacebookIncentive/", h(muxhandlers.GetFacebookIncentive, LogExecutionTime))
	//router.HandleFunc(prefix+"/Friend/setInviteCode/", h(muxhandlers.SetInviteCode, LogExecutionTime))
	//router.HandleFunc(prefix+"/Friend/setInviteHistory/", h(muxhandlers.SetInviteHistory, LogExecutionTime))
	//router.HandleFunc(prefix+"/Friend/setFacebookScopedId/", h(muxhandlers.SetFacebookScopedID, LogExecutionTime))
	//router.HandleFunc(prefix+"/Friend/getFriendUserIdList/", h(muxhandlers.GetFriendUserIDList, LogExecutionTime))
	//router.HandleFunc(prefix+"/Friend/requestEnergy/", h(muxhandlers.RequestEnergy, LogExecutionTime))

	// == Roulette ==
	router.HandleFunc(prefix+"/Spin/getWheelOptions/", h(muxhandlers.GetWheelOptions, LogExecutionTime))
	router.HandleFunc(prefix+"/Chao/getChaoWheelOptions/", h(muxhandlers.GetChaoWheelOptions, LogExecutionTime))
	router.HandleFunc(prefix+"/RaidbossSpin/getRaidbossWheelOptions/", h(muxhandlers.GetRaidbossWheelOptions, LogExecutionTime))
	router.HandleFunc(prefix+"/Chao/getPrizeChaoWheelSpin/", h(muxhandlers.GetPrizeChaoWheelSpin, LogExecutionTime))
	router.HandleFunc(prefix+"/RaidbossSpin/getPrizeRaidbossWheelSpin/", h(muxhandlers.GetPrizeRaidbossWheelSpin, LogExecutionTime))
	router.HandleFunc(prefix+"/RaidbossSpin/getItemStockNum/", h(muxhandlers.GetItemStockNum, LogExecutionTime))
	router.HandleFunc(prefix+"/Spin/commitWheelSpin/", h(muxhandlers.CommitWheelSpin, LogExecutionTime))
	router.HandleFunc(prefix+"/Chao/getFirstLaunchChao/", h(muxhandlers.GetFirstLaunchChao, LogExecutionTime))
	router.HandleFunc(prefix+"/Chao/commitChaoWheelSpin/", h(muxhandlers.CommitChaoWheelSpin, LogExecutionTime))
	//router.HandleFunc(prefix+"/RaidbossSpin/commitRaidbossWheelSpin/", h(muxhandlers.CommitRaidbossWheelSpin, LogExecutionTime))
	router.HandleFunc(prefix+"/Spin/getWheelSpinInfo/", h(muxhandlers.GetWheelSpinInfo, LogExecutionTime))

	// == Battle ==
	router.HandleFunc(prefix+"/Battle/getDailyBattleData/", h(muxhandlers.GetDailyBattleData, LogExecutionTime))
	router.HandleFunc(prefix+"/Battle/updateDailyBattleStatus/", h(muxhandlers.UpdateDailyBattleStatus, LogExecutionTime))
	router.HandleFunc(prefix+"/Battle/resetDailyBattleMatching/", h(muxhandlers.ResetDailyBattleMatching, LogExecutionTime))
	router.HandleFunc(prefix+"/Battle/getDailyBattleDataHistory/", h(muxhandlers.GetDailyBattleHistory, LogExecutionTime))
	router.HandleFunc(prefix+"/Battle/getDailyBattleStatus/", h(muxhandlers.GetDailyBattleStatus, LogExecutionTime))
	router.HandleFunc(prefix+"/Battle/postDailyBattleResult/", h(muxhandlers.PostDailyBattleResult, LogExecutionTime))
	router.HandleFunc(prefix+"/Battle/getPrizeDailyBattle/", h(muxhandlers.GetPrizeDailyBattle, LogExecutionTime))

	// == Gift Box ==
	router.HandleFunc(prefix+"/Message/getMessageList/", h(muxhandlers.GetMessageList, LogExecutionTime))
	router.HandleFunc(prefix+"/Message/getMessage/", h(muxhandlers.GetMessage, LogExecutionTime))
	//router.HandleFunc(prefix+"/Message/sendEnergy/", h(muxhandlers.SendEnergy, LogExecutionTime))

	// == Leaderboards & League ==
	router.HandleFunc(prefix+"/Leaderboard/getWeeklyLeaderboardOptions/", h(muxhandlers.GetWeeklyLeaderboardOptions, LogExecutionTime))
	router.HandleFunc(prefix+"/Leaderboard/getWeeklyLeaderboardEntries/", h(muxhandlers.GetWeeklyLeaderboardEntries, LogExecutionTime))
	router.HandleFunc(prefix+"/Leaderboard/getLeagueData/", h(muxhandlers.GetLeagueData, LogExecutionTime))
	//router.HandleFunc(prefix+"/Leaderboard/getLeagueOperatorData/", h(muxhandlers.GetLeagueOperatorData, LogExecutionTime))

	// == Character transactions ==
	router.HandleFunc(prefix+"/Character/unlockedCharacter/", h(muxhandlers.UnlockedCharacter, LogExecutionTime))
	router.HandleFunc(prefix+"/Character/upgradeCharacter/", h(muxhandlers.UpgradeCharacter, LogExecutionTime))

	// == Migration ==
	router.HandleFunc(prefix+"/Login/getMigrationPassword/", h(muxhandlers.GetMigrationPassword, LogExecutionTime))
	router.HandleFunc(prefix+"/Login/migration/", h(muxhandlers.Migration, LogExecutionTime))

	// == Debug ==
	//router.HandleFunc(prefix+"/Debug/addMessage/", h(muxhandlers.DebugSendMessage, LogExecutionTime))
	//router.HandleFunc(prefix+"/Debug/addOpeMessage/", h(muxhandlers.DebugSendOperatorMessage, LogExecutionTime))
	//router.HandleFunc(prefix+"/Debug/deleteUserData/", h(muxhandlers.DebugDeleteUserData, LogExecutionTime))
	//router.HandleFunc(prefix+"/Debug/forceDrawRaidboss/", h(muxhandlers.DebugForceDrawRaidboss, LogExecutionTime))
	//router.HandleFunc(prefix+"/Debug/getSpecialItem/", h(muxhandlers.DebugGetSpecialItem, LogExecutionTime))

	// == Misc. ==
	router.HandleFunc(prefix+"/Sgn/sendApollo/", h(muxhandlers.SendApollo, LogExecutionTime))
	router.HandleFunc(prefix+"/Sgn/setNoahId/", h(muxhandlers.SetNoahID, LogExecutionTime))
	//router.HandleFunc(prefix+"/Sgn/setSerialCode/", h(muxhandlers.SetSerialCode, LogExecutionTime))

	// Server information
	if config.CFile.EnablePublicStats {
		router.HandleFunc("/outrunInfo/stats", inforeporters.Stats)
	}

	// Noah-related files
	router.HandleFunc("/ds.txt", HandleDSFile)
	router.HandleFunc("/pp-ads.txt", HandlePPAdsFile)

	router.PathPrefix("/").HandlerFunc(HandleUnknownRequest)

	go bgtasks.TouchAnalyticsDB()

	port := config.CFile.Port
	log.Printf("Starting server on port %s\n", port)
	panic(http.ListenAndServe(":"+port, removePrependingSlashes(router)))
}
