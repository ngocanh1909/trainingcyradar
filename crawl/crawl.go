package crawl

import (
	"fmt"
	"github.com/ngocanh1909/config"
	"github.com/ngocanh1909/request"
	"regexp"
	"strings"
	"sync"
)

const URL = "https://malshare.com/daily/"
const LIMIT= 100000

func getHash(id int, date string) config.MalshareData {
	var result config.MalshareData
	url := fmt.Sprintf("https://malshare.com/daily/%s/malshare_fileList.%s.all.txt", date, date)
	dataStr, err := request.Request(url)
	if err != nil {
		return result
	}

	var md5 = ""
	var sha1 = ""
	var sha256 = ""
	fmt.Printf("Worker ID %d\n", id)
	for {
		if len(dataStr) <= 0 {
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

var HashData []config.MalshareData
var wg sync.WaitGroup

func worker(id int, jobs <-chan string, results chan<- config.MalshareData, wg *sync.WaitGroup) {
	defer wg.Done()
	for j := range jobs {
		fmt.Printf("worker %d start jobs http://malshare.com/daily/%s/malshare_fileList.%s.all.txt \n", id, j, j)
		results <- getHash(id, j)
		fmt.Printf("worker %d finished jobs http://malshare.com/daily/%s/malshare_fileList.%s.all.txt \n", id, j, j)
	}
}

func DumpData() ([]config.MalshareData, error) {
	bodyStr, err := request.Request(URL)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	jobs := make(chan string, LIMIT)
	results := make(chan config.MalshareData, LIMIT)
	magic := regexp.MustCompile(`\"\[DIR\]\"></[a-z]{2}><[a-z]{2}><a\s[a-z]{4}=\"`)
	magicStr := string(magic.Find([]byte(bodyStr)))
	end := regexp.MustCompile("_[a-z]{8}/")
	outEnd := string(end.Find([]byte(bodyStr)))
	for w := 1; w < 101; w++ {
		wg.Add(1)
		go worker(w, jobs, results, &wg)
	}
	c := 0
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
		if (dateStr == "2019-11-14" || dateStr == "2019-12-01") {
			continue
		}
		jobs <- dateStr
		c++
	}
	close(jobs)
	wg.Wait()
	for a := 1; a <= c; a++ {
		HashArr(<-results)
	}
	return HashData, nil
}

func HashArr(data config.MalshareData) {
	HashData = append(HashData, data)
}
