package config

import (
	"sync"
	"time"
)

type MalshareData struct {
	Date   time.Time `json:"date" bson:"date"`
	Md5    string    `json:"md5" bson:"md5"`
	Sha1   string    `json:"sha1" bson:"sha1"`
	Sha256 string    `json:"sha256" bson:"sha256"`
}
type WaitGroup struct {
	Wait sync.WaitGroup
}
