package cryption

import (
	"log"
	"net/http"
	"regexp"
)

var EncryptionIv = []byte("")

//var EncryptionKey = []byte("Ec7bLaTdSuXuf5pW") // encryption key in UnmanagedProcess
var EncryptionKey = []byte("vMdkkY8bfVmUS6qr") // encryption key in Assembly-CSharp

func CleanBytes(b []byte) []byte {
	re := regexp.MustCompile("[^\u0020-\u007f]+")
	return []byte(re.ReplaceAllLiteralString(string(b), ""))
}

func GetReceivedMessage(r *http.Request) ([]byte, error) {
	err := r.ParseForm()
	if err != nil {
		log.Println("[ERR] Error in parsing form: " + err.Error())
	}
	param := r.Form.Get("param")
	iv := r.Form.Get("key")
	secure := r.Form.Get("secure")
	if secure != "1" {
		return []byte(param), nil // message is not encrypted!
	}
	defer func() {
		if rErr := recover(); rErr != nil {
			log.Println("[ERR] Panic during decryption", rErr)
		}
	}()

	EncryptionIv = []byte(iv)
	paramUnB64, err := B64Decode(param)
	if err != nil {
		log.Println("[ERR] Couldn't decode param: " + err.Error())
		return []byte(param), err
	}
	decrypted := Decrypt(paramUnB64, EncryptionKey, EncryptionIv)
	decrypted = CleanBytes(decrypted)
	return decrypted, nil
}
