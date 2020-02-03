package main

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"time"
)

const (
	hosts      = "localhost:27017"
	database   = "crawl"
	username   = "admin1"
	password   = "admin1"
	collection = "post"
)

type MalshareData struct {
	ID     bson.ObjectId `bson:"_id,omitempty"`
	Date   time.Time     `bson:"date"`
	Md5    string        `bson:"md5"`
	Sha1   string        `bson:"sha1"`
	Sha256 string        `bson:"sha256"`
}

func DBConn() (session *mgo.Session) {
	info := &mgo.DialInfo{
		Addrs:    []string{hosts},
		Timeout:  60 * time.Second,
		Database: database,
		Username: username,
		Password: password,
	}
	session, err := mgo.DialWithInfo(info)
	if err != nil {
		errors.New("Connect fail")
		return
	}
	return session
}

func processGET(c *gin.Context) {
	s := DBConn()
	collection := s.DB(database).C(collection)
	malshareData := MalshareData{}
	date := c.Params.ByName("date")
	dateParse, _ := time.Parse("2006-01-02", date)
	query := bson.M{
		"date": dateParse,
	}
	err := collection.Find(query).One(&malshareData)
	if err == nil {
		c.JSON(http.StatusOK, gin.H{
			"date":   dateParse,
			"status": "ok",
			"md5":    malshareData.Md5,
			"sha1":   malshareData.Sha1,
			"sha256": malshareData.Sha256,
		})
	}
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"date": dateParse, "status": "no value"})
	}
}

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/malshare/:date", processGET)
	return r
}

func main() {
	r := setupRouter()
	r.Run(":8080")
}
