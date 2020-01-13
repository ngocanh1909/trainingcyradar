package savemgo

import (
	"fmt"
	"github.com/ngocanh1909/request"
	"gopkg.in/mgo.v2"
	"regexp"
	"strings"
)

const URL = "https://malshare.com/daily/"
const LIMIT = 100000

type MalshareData struct {
	Date   string `json:"date" bson:"date"`
	Md5    string `json:"md5" bson:"md5"`
	Sha1   string `json:"sha1" bson:"sha1"`
	Sha256 string `json:"sha256" bson:"sha256"`
}

const (
	hosts      = "localhost:27017"
	database   = "crawl"
	username   = "admin1"
	password   = "admin1"
	collection = "post"
)

func getHash(date string) MalshareData {
	var result MalshareData
	url := fmt.Sprintf("https://malshare.com/daily/%s/malshare_fileList.%s.all.txt", date, date)
	dataStr, err := request.Request(url)
	if err != nil {
		return result
	}
	var md5 = ""
	var sha1 = ""
	var sha256 = ""
	for {
		if len(dataStr) == 0 {
			break
		}
		md5Str := dataStr[0:32]
		sha1Str := dataStr[33:73]
		sha256Str := dataStr[74:138]
		i := strings.Index(dataStr, "\n")
		dataStr = dataStr[i+1 : len(dataStr)]

		md5 = fmt.Sprintf("%s %s", md5, md5Str)
		sha1 = fmt.Sprintf("%s %s", sha1, sha1Str)
		sha256 = fmt.Sprintf("%s %s", sha256, sha256Str)
	}
	result.Date = date
	result.Md5 = md5
	result.Sha1 = sha1
	result.Sha256 = sha256
	return result
}

func worker(id int, jobs <-chan string, results chan<- MalshareData) {
	for j := range jobs {
		fmt.Printf("worker %d start jobs http://malshare.com/daily/%s/malshare_fileList.%s.all.txt \n", id, j, j)
		results <- getHash(j)
		fmt.Printf("worker %d finished jobs http://malshare.com/daily/%s/malshare_fileList.%s.all.txt \n", id, j, j)
	}
}

var session *mgo.Session;

func DumpData(s *mgo.Session) {
	session = s
	bodyStr, err := request.Request(URL)
	if err != nil {
		fmt.Println(err)
		return
	}
	jobs := make(chan string, LIMIT)
	results := make(chan MalshareData, LIMIT)
	magic := regexp.MustCompile(`\"\[DIR\]\"></[a-z]{2}><[a-z]{2}><a\s[a-z]{4}=\"`)
	magicStr := string(magic.Find([]byte(bodyStr)))
	end := regexp.MustCompile("_[a-z]{8}/")
	outEnd := string(end.Find([]byte(bodyStr)))
	for w := 1; w < 101; w++ {
		go worker(w, jobs, results)
	}
	for {
		i := strings.Index(bodyStr, magicStr)
		re := regexp.MustCompile("=\"\\d{4}-\\d{2}-\\d{2}")
		out := re.Find([]byte(bodyStr))
		if len(out) < 3 {
			break
		}
		dateStr := string(out)[2:]
		bodyStr = bodyStr[i+len(magicStr)+1 : len(bodyStr)]
		if dateStr == outEnd {
			break
		}
		jobs <- dateStr
	}
	close(jobs)
	for a := 1; a <= LIMIT; a++ {
		SaveFile(<-results)
	}
}

func SaveFile(data MalshareData) {
	col := session.DB(database).C(collection)
	col.Insert(data)
}
