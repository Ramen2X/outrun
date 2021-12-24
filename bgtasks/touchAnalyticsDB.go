package bgtasks

import (
	"log"

	"github.com/Ramen2X/outrun/consts"
	"github.com/Ramen2X/outrun/db/dbaccess"
)

func TouchAnalyticsDB() {
	err := dbaccess.Set(consts.DBBucketAnalytics, "touch", []byte{})
	if err != nil {
		log.Println("[ERR] Unable to touch " + consts.DBBucketAnalytics + ": " + err.Error())
	}
}
