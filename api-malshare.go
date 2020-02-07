package main

import (
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
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

type MalshareDAO struct {
	*mgo.Database
}

func (mal *MalshareDAO) processGET(c *gin.Context) {
	malshareData := MalshareData{}
	date := c.Params.ByName("date")
	dateParse, err := time.Parse("2006-01-02", date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"messages" : "Date Parse Fail",
			"date":err.Error(),
			"md5": err.Error(),
			"sha1": err.Error(),
			"sha256": err.Error(),
		})
		return
	}
	query := bson.M{
		"date": dateParse,
	}
	err = mal.C(collection).Find(query).One(&malshareData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"messages" : "Query Fail",
			"date":err.Error(),
			"md5": err.Error(),
			"sha1": err.Error(),
			"sha256": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"messages": "ok",
		"date":   dateParse,
		"md5":    malshareData.Md5,
		"sha1":   malshareData.Sha1,
		"sha256": malshareData.Sha256,
	})
}

func (mal *MalshareDAO) setupRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/malshare/:date", mal.processGET)
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
		log.Fatal(err)
	}
	d := MalshareDAO{session.DB(database)}
	r := d.setupRouter()
	r.Run(":8080")
}
