package save

import (
	"github.com/ngocanh1909/config"
	"gopkg.in/mgo.v2"
)

func SaveFile(db *mgo.Database, hashData [] config.MalshareData) (error) {
	col := db.C("post")
	for i := 0; i < len(hashData); i++ {
		err := col.Insert(hashData[i])
		if err != nil {
			return err
		}
	}
	return nil
}
