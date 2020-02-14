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
	"os"
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
	err = mal.db.C(config.DB.Collection).Find(query).One(&malshareData)
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

func (mal *MalshareDAO) SetupRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/malshare/:date", mal.processGET)
	return r
}

func main() {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(dir)
	var config config.Config
	if _, err := toml.DecodeFile("config.toml", &config); err != nil {
		log.Fatal(err)
	}
	info := &mgo.DialInfo{
		Addrs:    []string{config.DB.Server},
		Timeout:  60 * time.Second,
		Database: config.DB.Database,
		Username: config.DB.Username,
		Password: config.DB.Password,
	}
	session, err := mgo.DialWithInfo(info)
	if err != nil {
		log.Fatal(err)
	}
	var hashData [] models.MalshareData
	hashData, err = crawl.DumpData(&models.WaitGroup{})
	if err != nil {
		log.Fatal(err)
	}
	choose := flag.String("command", "mgo", "-command=<choose>")
	flag.Parse()
	if *choose == "file" {
		for i := 0; i < len(hashData); i++ {
			err := save.SaveFile(hashData[i])
			if err != nil {
				log.Fatal(err)
			}
		}
	}
	if *choose == "mgo" {
		err := save.SaveMgo(session.DB(config.DB.Database), hashData)
		if err != nil {
			log.Fatal(err)
		}
	}
	if *choose == "api" {
		d := MalshareDAO{session.DB(config.DB.Database)}
		r := d.SetupRouter()
		err = r.Run(":8080")
		if err != nil {
			log.Fatal(err)
		}
	}
}
