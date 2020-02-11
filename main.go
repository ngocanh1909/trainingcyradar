package main

import (
	"flag"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/gin-gonic/gin"
	"github.com/ngocanh1909/trainingcyradar/config"
	"github.com/ngocanh1909/trainingcyradar/crawl"
	"github.com/ngocanh1909/trainingcyradar/models"
	"github.com/ngocanh1909/trainingcyradar/save"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"net/http"
	"time"
)

type MalshareDAO struct {
	db *mgo.Database
}

func (mal *MalshareDAO) processGET(c *gin.Context) {
	var config config.Config
	if _, err := toml.DecodeFile("config.toml", &config); err != nil {
		fmt.Println(err)
		return
	}

	malshareData := models.MalshareData{}
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
	err = mal.db.C(config.Database.Collection).Find(query).One(&malshareData)
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
	var config config.Config
	if _, err := toml.DecodeFile("config.toml", &config); err != nil {
		log.Fatal(err)
	}
	info := &mgo.DialInfo{
		Addrs:    []string{config.Database.Server},
		Database: config.Database.Database,
		Username: config.Database.User,
		Password: config.Database.Password,
	}
	session, err := mgo.DialWithInfo(info)
	if err != nil {
		log.Fatal(err)
	}
	var hashData [] models.MalshareData
	hashData, err = crawl.DumpData(&models.WaitGroup{})
	wordPtr := flag.String("command", "file", "a string")
	flag.Parse()
	if (*wordPtr == "file") {
		for i := 0; i < len(hashData); i++ {
			save.SaveFile(hashData[i])
		}
	}
	if (*wordPtr == "mgo") {
		save.SaveMgo(session.DB(config.Database.Database), hashData)

	}
	if (*wordPtr == "api") {
		d := MalshareDAO{session.DB(config.Database.Database)}
		r := d.setupRouter()
		r.Run(":8080")
	}
}
