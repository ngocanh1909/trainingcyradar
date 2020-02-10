package main

import (
	"github.com/gin-gonic/gin"
	"github.com/ngocanh1909/trainingcyradar/config"
	"github.com/ngocanh1909/trainingcyradar/crawl"
	"github.com/ngocanh1909/trainingcyradar/save"
	"flag"
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

type MalshareDAO struct {
	*mgo.Database
}

func (mal *MalshareDAO) processGET(c *gin.Context) {
	malshareData := config.MalshareData{}
	date := c.Params.ByName("date")
	dateParse, err := time.Parse("2006-01-02", date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"messages": err.Error(),
		})
		return
	}
	query := bson.M{
		"date": dateParse,
	}
	err = mal.C(collection).Find(query).One(&malshareData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"messages": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"messages": "ok",
		"date":     dateParse,
		"md5":      malshareData.Md5,
		"sha1":     malshareData.Sha1,
		"sha256":   malshareData.Sha256,
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
	var hashData []config.MalshareData
	hashData, err = crawl.DumpData(&config.WaitGroup{})
	wordPtr := flag.String("command", "file", "go run main.go [-comand=<name>]")
	flag.Parse()
	if (*wordPtr == "file") {
		for i := 0; i < len(hashData); i++ {
			save.SaveFile(hashData[i])
		}
	}
	if (*wordPtr == "mgo") {
		save.SaveMgo(session.DB(database), hashData)

	}
	if (*wordPtr == "api") {
		d := MalshareDAO{session.DB(database)}
		r := d.setupRouter()
		r.Run(":8080")
	}
}
