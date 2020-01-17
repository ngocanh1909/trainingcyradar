package save

import (
	"fmt"
	"github.com/ngocanh1909/config"
	"gopkg.in/mgo.v2"
)

func SaveFile(db *mgo.Database, session *mgo.Session, hashData [] config.MalshareData){
	col := db.C("post")
	for i := 0; i < len(hashData); i++ {
		err := col.Insert(hashData[i])
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}
