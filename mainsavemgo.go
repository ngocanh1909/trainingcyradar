package main

import (
	"github.com/ngocanh1909/client"
	"gopkg.in/mgo.v2"
	"time"
)

const (
	hosts      = "localhost:27017"
	database   = "crawl"
	username   = "admin1"
	password   = "admin1"
	collection = "post"
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
	client.DumpData(session)
}