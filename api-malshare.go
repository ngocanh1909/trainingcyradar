package main

import (
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
	Date   string        `bson:"date"`
	Md5    string        `bson:"md5"`
	Sha1   string        `bson:"sha1"`
	Sha256 string        `bson:"sha256"`
}

var db = make(map[string]string)

func processGET(c *gin.Context) {
	info := &mgo.DialInfo{
		Addrs:    []string{hosts},
		Timeout:  60 * time.Second,
		Database: database,
		Username: username,
		Password: password,
	}
	session, err := mgo.DialWithInfo(info)
	if err != nil {
		panic(err)
	}
	collection := session.DB(database).C(collection)
	date := c.Params.ByName("date")
	query := bson.M{
		"date": date,
	}
	malshareData := MalshareData{}
	err = collection.Find(query).One(&malshareData)
	//fmt.Print(malshareData)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"date": date, "status": "no value"})
	} else {
		c.JSON(http.StatusOK, gin.H{"date": date, "status": "ok", "md5": malshareData.Md5})
	}
}

func setupRouter() *gin.Engine {
	// Disable Console Color
	// gin.DisableConsoleColor()
	r := gin.Default()
	// Get user value
	r.GET("/malshare/:date", processGET)
	return r
}

func main() {
	r := setupRouter()
	// Listen and Server in 0.0.0.0:8080
	r.Run(":8080")
}
