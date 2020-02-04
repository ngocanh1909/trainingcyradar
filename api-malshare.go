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

var coll *mgo.Collection

type MalshareData struct {
	ID     bson.ObjectId `bson:"_id,omitempty"`
	Date   time.Time     `bson:"date"`
	Md5    string        `bson:"md5"`
	Sha1   string        `bson:"sha1"`
	Sha256 string        `bson:"sha256"`
}

func processGET(c *gin.Context) {

	malshareData := MalshareData{}
	date := c.Params.ByName("date")
	dateParse, err := time.Parse("2006-01-02", date)
	if err != nil {
		errors.New("")
	}
	query := bson.M{
		"date": dateParse,
	}
	err = coll.Find(query).One(&malshareData)
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{"date": dateParse, "status": "no value"})
	}
	if err == nil {
		c.JSON(http.StatusOK, gin.H{
			"date":   dateParse,
			"status": "ok",
			"md5":    malshareData.Md5,
			"sha1":   malshareData.Sha1,
			"sha256": malshareData.Sha256,
		})
	}

}

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/malshare/:date", processGET)
	return r
}

func main() {
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
	coll = session.DB(database).C(collection)
	r := setupRouter()
	r.Run(":8080")
}
