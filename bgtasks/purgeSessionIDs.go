package bgtasks

import (
	"time"

	"github.com/Ramen2X/outrun/db"
)

func PurgeSessionIDs() {
	for true {
		time.Sleep(10 * time.Minute)
		db.PurgeAllExpiredSessionIDs()
	}
}
