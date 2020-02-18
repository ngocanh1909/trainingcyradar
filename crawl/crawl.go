package crawl

import (
	"fmt"
	"github.com/ngocanh1909/trainingcyradar/models"
	"github.com/ngocanh1909/trainingcyradar/request"
	"regexp"
	"strings"
	"sync"
	"time"
)

const URL = "https://malshare.com/daily/"
const LIMIT = 100000

func getHash(id int, date time.Time) models.Malshare {
	var result models.Malshare
	url := fmt.Sprintf("https://malshare.com/daily/%s/malshare_fileList.%s.all.txt", date.Format("2006-01-02"), date.Format("2006-01-02"))
	dataStr, err := request.Request(url)
	if err != nil {
		result.Err = err
		return result
	}
	var md5 = ""
	var sha1 = ""
	var sha256 = ""
	for {
		if len(dataStr) < 138 {
			break
		}
		md5Str := dataStr[0:32]
		sha1Str := dataStr[33:73]
		sha256Str := dataStr[74:138]
		i := strings.Index(dataStr, "\n")
		if i+1 > len(dataStr) {
			break
		}
		dataStr = dataStr[i+1:]
		md5 = fmt.Sprintf("%s %s", md5, md5Str)
		sha1 = fmt.Sprintf("%s %s", sha1, sha1Str)
		sha256 = fmt.Sprintf("%s %s", sha256, sha256Str)
	}
	result.Mal.Date = date
	result.Mal.Md5 = md5
	result.Mal.Sha1 = sha1
	result.Mal.Sha256 = sha256
	return result
}

func worker(id int, jobs <-chan time.Time, results chan<- models.Malshare, wg *sync.WaitGroup) {
	defer wg.Done()
	for j := range jobs {
		fmt.Printf("worker %d start jobs http://malshare.com/daily/%s/malshare_fileList.%s.all.txt \n", id, j.Format("2006-01-02"), j.Format("2006-01-02"))
		results <- getHash(id, j)
		fmt.Printf("worker %d finished jobs http://malshare.com/daily/%s/malshare_fileList.%s.all.txt \n", id, j.Format("2006-01-02"), j.Format("2006-01-02"))
	}
}

func DumpData(wg *sync.WaitGroup) ([]models.MalshareData, error) {
	var HashArray []models.MalshareData
	bodyStr, err := request.Request(URL)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	jobs := make(chan time.Time, LIMIT)
	results := make(chan models.Malshare, LIMIT)
	magic := regexp.MustCompile(`\"\[DIR\]\"></[a-z]{2}><[a-z]{2}><a\s[a-z]{4}=\"`)
	magicStr := string(magic.Find([]byte(bodyStr)))
	end := regexp.MustCompile("_[a-z]{8}/")
	outEnd := string(end.Find([]byte(bodyStr)))
	for w := 1; w < 100; w++ {
		wg.Add(1)
		go worker(w, jobs, results, wg)
	}
	for {
		i := strings.Index(bodyStr, magicStr)
		re := regexp.MustCompile("=\"\\d{4}-\\d{2}-\\d{2}")
		out := re.Find([]byte(bodyStr))
		if len(out) < 3 {
			break
		}
		dateStr := string(out)[2:]
		if i+len(magicStr)+1 > len(bodyStr) {
			break
		}
		bodyStr = bodyStr[i+len(magicStr)+1:]
		if dateStr == outEnd {
			break
		}
		t, _ := time.Parse("2006-01-02", dateStr)
		jobs <- t
	}
	close(jobs)
	wg.Wait()
	close(results)
	for i := range results {
		if i.Err != nil {
			fmt.Println(i.Err)
			continue
		}
		HashArray = append(HashArray, i.Mal)
	}
	return HashArray, nil
}
