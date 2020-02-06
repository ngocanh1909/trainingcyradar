package main

import (
	"fmt"
	"github.com/ngocanh1909/trainingcyradar/config"
	"github.com/ngocanh1909/trainingcyradar/crawl"
	"github.com/ngocanh1909/trainingcyradar/save"
	"gopkg.in/mgo.v2"
	"time"
)

const (
	hosts    = "localhost:27017"
	database = "crawl"
	username = "admin1"
	password = "admin1"
)

type Post struct {
	Tile    string
	Content string
}

func main() {
	// Tạo phiên kết nối với MongDB
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
	var hashData []config.MalshareData
	hashData, err = crawl.DumpData(&config.WaitGroup{})
	fmt.Println("1")
	save.SaveFile(session.DB(database), hashData)
}
