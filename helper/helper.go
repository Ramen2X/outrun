package helper

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Ramen2X/outrun/config"
	"github.com/Ramen2X/outrun/cryption"
	"github.com/Ramen2X/outrun/db"
	"github.com/Ramen2X/outrun/emess"
	"github.com/Ramen2X/outrun/netobj"
	"github.com/Ramen2X/outrun/netobj/constnetobjs"
	"github.com/Ramen2X/outrun/requests"
	"github.com/Ramen2X/outrun/responses"
	"github.com/Ramen2X/outrun/responses/responseobjs"
	"github.com/Ramen2X/outrun/status"
)

const (
	PrefixErr            = "ERR"
	PrefixOut            = "OUT"
	PrefixWarn           = "WARN"
	PrefixUncatchableErr = "UNCATCHABLE ERR"
	PrefixDebugOut       = "DEBUG (OUT)"

	LogOutBase = "[%s] (%s) %s\n"
	LogErrBase = "[%s] (%s) %s: %s\n"

	InternalServerError = "Internal server error"
	BadRequest          = "Bad request"

	//DefaultIV = "HotAndSunnyMiami"
	DefaultIV = "FoundDeadInMiami"
)

type Helper struct {
	CallerName string
	RespW      http.ResponseWriter
	Request    *http.Request
}

func MakeHelper(callerName string, r http.ResponseWriter, request *http.Request) *Helper {
	return &Helper{
		callerName,
		r,
		request,
	}
}

func (r *Helper) GetGameRequest() []byte {
	recv, err := cryption.GetReceivedMessage(r.Request)
	if err != nil {
		r.SendCompatibleResponse(responses.NewBaseResponse(r.BaseInfo(emess.OK, status.DecryptionFailure)), false)
	}
	return recv
}
func (r *Helper) SendResponse(i interface{}) error {
	out, err := json.Marshal(i)
	if err != nil {
		return err
	}
	r.Respond(out, true)
	return nil
}
func (r *Helper) SendInsecureResponse(i interface{}) error {
	out, err := json.Marshal(i)
	if err != nil {
		return err
	}
	r.RespondInsecure(out, true)
	return nil
}
func (r *Helper) RespondRaw(out []byte, secureFlag, iv string, sendErrorResponseOnError bool) {
	if config.CFile.LogAllResponses {
		nano := time.Now().UnixNano()
		nanoStr := strconv.Itoa(int(nano))
		filename := r.Request.RequestURI + "--" + nanoStr
		filename = strings.ReplaceAll(filename, ".", "-")
		filename = strings.ReplaceAll(filename, "/", "-") + ".txt"
		filepath := "logging/all_responses/" + filename
		r.Out("DEBUG: Saving response to " + filepath)
		err := ioutil.WriteFile(filepath, out, 0644)
		if err != nil {
			r.Out("DEBUG ERROR: Unable to write file '" + filepath + "'")
		}
	}
	response := map[string]string{}
	if secureFlag != "0" && secureFlag != "1" {
		r.Warn("Improper secureFlag in call to RespondRaw!")
	}
	response["secure"] = secureFlag
	response["key"] = iv
	if secureFlag == "1" {
		encrypted := cryption.Encrypt(out, cryption.EncryptionKey, []byte(iv))
		encryptedBase64 := cryption.B64Encode(encrypted)
		response["param"] = encryptedBase64
	} else {
		response["param"] = string(out)
	}
	toClient, err := json.Marshal(response)
	if err != nil {
		if sendErrorResponseOnError {
			r.SendCompatibleResponse(responses.NewBaseResponse(r.BaseInfo(emess.OK, status.InvalidResponse)), false)
		}
		r.InternalErr("Error marshalling in RespondRaw", err)
		return
	}
	r.RespW.Write(toClient)
}
func (r *Helper) SendCompatibleResponse(out interface{}, sendErrorResponseOnError bool) error {
	response := map[string]interface{}{"secure": "0", "param": out}
	toClient, err := json.Marshal(response)
	if err != nil {
		if sendErrorResponseOnError {
			r.SendCompatibleResponse(responses.NewBaseResponse(r.BaseInfo(emess.OK, status.InvalidResponse)), false)
		}
		r.InternalErr("Error marshalling in SendCompatibleResponse", err)
		return err
	}
	if config.CFile.LogAllResponses {
		nano := time.Now().UnixNano()
		nanoStr := strconv.Itoa(int(nano))
		filename := r.Request.RequestURI + "--" + nanoStr
		filename = strings.ReplaceAll(filename, ".", "-")
		filename = strings.ReplaceAll(filename, "/", "-") + ".txt"
		filepath := "logging/all_responses/" + filename
		r.Out("DEBUG: Saving response to " + filepath)
		err := ioutil.WriteFile(filepath, toClient, 0644)
		if err != nil {
			r.Out("DEBUG ERROR: Unable to write file '" + filepath + "'")
		}
	}
	r.RespW.Write(toClient)
	return nil
}
func (r *Helper) Respond(out []byte, sendErrorResponseOnError bool) {
	r.RespondRaw(out, "1", DefaultIV, sendErrorResponseOnError)
}
func (r *Helper) RespondInsecure(out []byte, sendErrorResponseOnError bool) {
	r.RespondRaw(out, "0", "", sendErrorResponseOnError)
}
func (r *Helper) Out(s string, a ...interface{}) {
	msg := fmt.Sprintf(s, a...)
	log.Printf(LogOutBase, PrefixOut, r.CallerName, msg)
}
func (r *Helper) DebugOut(s string, a ...interface{}) {
	if config.CFile.DebugPrints {
		msg := fmt.Sprintf(s, a...)
		log.Printf(LogOutBase, PrefixDebugOut, r.CallerName, msg)
	}
}
func (r *Helper) Warn(s string, a ...interface{}) {
	msg := fmt.Sprintf(s, a...)
	log.Printf(LogOutBase, PrefixWarn, r.CallerName, msg)
}
func (r *Helper) WarnErr(msg string, err error) {
	log.Printf(LogErrBase, PrefixWarn, r.CallerName, msg, err.Error())
}
func (r *Helper) Uncatchable(msg string) {
	log.Printf(LogOutBase, PrefixOut, r.CallerName, msg)
}
func (r *Helper) InternalErr(msg string, err error) {
	log.Printf(LogErrBase, PrefixErr, r.CallerName, msg, err.Error())
	//	r.RespW.WriteHeader(http.StatusBadRequest)
	//	r.RespW.Write([]byte(BadRequest))
	r.SendCompatibleResponse(responses.NewBaseResponse(r.BaseInfo(emess.OK, status.InternalServerError)), false)
}
func (r *Helper) Err(msg string, err error) {
	log.Printf(LogErrBase, PrefixErr, r.CallerName, msg, err.Error())
	//	r.RespW.WriteHeader(http.StatusBadRequest)
	//	r.RespW.Write([]byte(BadRequest))
	r.SendCompatibleResponse(responses.NewBaseResponse(r.BaseInfo(emess.OK, status.ServerSystemError)), false)
}
func (r *Helper) ErrRespond(msg string, err error, response string) {
	// TODO: remove if never used in stable builds
	log.Printf(LogErrBase, PrefixErr, r.CallerName, msg, err.Error())
	r.RespW.WriteHeader(http.StatusInternalServerError) // ideally include an option for this, but for now it's inconsequential
	r.RespW.Write([]byte(response))
}
func (r *Helper) InternalFatal(msg string, err error) {
	log.Fatalf(LogErrBase, PrefixErr, r.CallerName, msg, err.Error())
	//	r.RespW.WriteHeader(http.StatusBadRequest)
	//	r.RespW.Write([]byte(BadRequest))
	r.SendCompatibleResponse(responses.NewBaseResponse(r.BaseInfo(emess.OK, status.InternalServerError)), false)
}
func (r *Helper) Fatal(msg string, err error) {
	log.Fatalf(LogErrBase, PrefixErr, r.CallerName, msg, err.Error())
	//	r.RespW.WriteHeader(http.StatusBadRequest)
	//	r.RespW.Write([]byte(BadRequest))
	r.SendCompatibleResponse(responses.NewBaseResponse(r.BaseInfo(emess.OK, status.ServerSystemError)), false)
}
func (r *Helper) BaseInfo(em string, statusCode int64) responseobjs.BaseInfo {
	return responseobjs.NewBaseInfo(em, statusCode)
}
func (r *Helper) InvalidRequest() {
	//	r.RespW.WriteHeader(http.StatusBadRequest)
	//	r.RespW.Write([]byte(BadRequest))
	r.SendCompatibleResponse(responses.NewBaseResponse(r.BaseInfo(emess.OK, status.ClientError)), false)
}
func (r *Helper) GetCallingPlayer(sendErrorResponseOnError bool) (netobj.Player, error) {
	// Powerful function to get the player directly from the response
	recv := r.GetGameRequest()
	var request requests.Base
	err := json.Unmarshal(recv, &request)
	if err != nil {
		if sendErrorResponseOnError {
			r.SendCompatibleResponse(responses.NewBaseResponse(r.BaseInfo(emess.OK, status.RequestParamError)), false)
		}
		return constnetobjs.BlankPlayer, err
	}
	sid := request.SessionID
	player, err := db.GetPlayerBySessionID(sid)
	if err != nil {
		if sendErrorResponseOnError {
			r.SendCompatibleResponse(responses.NewBaseResponse(r.BaseInfo(emess.OK, status.ExpiredSession)), false)
		}
		return constnetobjs.BlankPlayer, err
	}
	if config.CFile.PrintPlayerNames {
		r.Out("Player '" + player.Username + "' (" + player.ID + ")")
	}
	return player, nil
}
