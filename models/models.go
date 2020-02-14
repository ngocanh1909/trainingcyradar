package models

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

type MalshareData struct {
	ID     bson.ObjectId `bson:"_id,omitempty"`
	Date   time.Time     `json:"date" bson:"date"`
	Md5    string        `json:"md5" bson:"md5"`
	Sha1   string        `json:"sha1" bson:"sha1"`
	Sha256 string        `json:"sha256" bson:"sha256"`
}

type Malshare struct {
	Mal MalshareData
	Err error
}

