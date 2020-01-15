package save

import (
	"github.com/ngocanh1909/config"
	"gopkg.in/mgo.v2"
	"log"
)

const (
	database   = "crawl"
	collection = "post"
)

func SaveFile(session *mgo.Session, hashData [] config.MalshareData) {
	col := session.DB(database).C(collection)
	for i := 0; i < len(hashData); i++ {
		err := col.Insert(hashData[i])
		if err != nil {
			log.Fatal(err)
		}
	}
}
