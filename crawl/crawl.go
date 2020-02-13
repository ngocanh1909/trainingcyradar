package crawl

import (
	"fmt"
	"github.com/ngocanh1909/trainingcyradar/models"
	"github.com/ngocanh1909/trainingcyradar/request"
	"regexp"
	"strings"
	"time"
)

const URL = "https://malshare.com/daily/"
const LIMIT = 100000

func getHash(id int, date time.Time) models.Malshare {
	var result models.Malshare
	url := fmt.Sprintf("https://malshare.com/daily/%s/malshare_fileList.%s.all.txt", date.Format("2006-01-02"), date.Format("2006-01-02"))
	dataStr, err := request.Request(url)
	if err != nil {
		err = result.Err
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
	result.Mal.Date = date
	result.Mal.Md5 = md5
	result.Mal.Sha1 = sha1
	result.Mal.Sha256 = sha256
	return result
}

func worker(id int, jobs <-chan time.Time, results chan<- models.Malshare, wg *models.WaitGroup) {
	defer wg.Wait.Done()
	for j := range jobs {
		fmt.Printf("worker %d start jobs http://malshare.com/daily/%s/malshare_fileList.%s.all.txt \n", id, j.Format("2006-01-02"), j.Format("2006-01-02"))
		results <- getHash(id, j)
		fmt.Printf("worker %d finished jobs http://malshare.com/daily/%s/malshare_fileList.%s.all.txt \n", id, j.Format("2006-01-02"), j.Format("2006-01-02"))
	}
}

func DumpData(wg *models.WaitGroup) ([]models.MalshareData, error) {
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
	for w := 1; w < 101; w++ {
		wg.Wait.Add(1)
		go worker(w, jobs, results, wg)
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
		t, _ := time.Parse("2006-01-02", dateStr)
		jobs <- t
		c++
		if c > 5 {
			break
		}
	}
	close(jobs)
	wg.Wait.Wait()
	close(results)
	for i := range results {
		HashArray = append(HashArray, i.Mal)
		if (i.Err != nil) {
			return nil, i.Err
		}
	}
	return HashArray, nil
}
