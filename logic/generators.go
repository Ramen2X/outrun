package logic

import (
	"crypto/md5"
	"encoding/hex"

	"github.com/fluofoxxo/outrun/netobj"
)

// Generates the correct login password
func GenerateLoginPassword(player netobj.Player) string {
	data := []byte(player.Key + ":dho5v5yy7n2uswa5iblb:" + player.ID + ":" + player.Password)
	sum := md5.Sum(data)
	return hex.EncodeToString(sum[:])
}
