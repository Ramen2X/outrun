package netobj

import "time"

type ServerHighScoreEntry struct {
	HighScore int64  `json:"highScore"`
	UserID    string `json:"userId"`
	Timestamp int64  `json:"timestamp"`
}

func NewServerHighScoreEntry(highScore int64, uid string) ServerHighScoreEntry {
	currentTime := time.Now().UTC().Unix()
	return ServerHighScoreEntry{
		highScore,
		uid,
		currentTime,
	}
}
