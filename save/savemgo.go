package save

import (
	"github.com/ngocanh1909/trainingcyradar/models"
	"gopkg.in/mgo.v2"
)

func SaveMgo(db *mgo.Database, hashData []models.MalshareData) (error) {
	col := db.C("post")
	for i := 0; i < len(hashData); i++ {
		err := col.Insert(hashData[i])
		if err != nil {
			return err
		}
	}
	return nil
}
