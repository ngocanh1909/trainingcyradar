package Client

import (
	"fmt"
	"github.com/ngocanh1909/request"
	"os"
	"regexp"
	"strings"
)

const URL  = "https://malshare.com/daily/"
const LIMIT = 100000

type malshareData struct {
	date string
	md5[] string
	sha1[] string
	sha256[] string
}

func getHash(date string) malshareData{
	var result malshareData
	url := fmt.Sprintf("https://malshare.com/daily/%s/malshare_fileList.%s.all.txt", date, date)
	dataStr, err := request.Request(url)
	if err != nil{
		return result
	}
	var md5Array[] string
	var sha1Aarray[] string
	var sha256Array[] string
	for{
		if len(dataStr) == 0{
			break
		}
		md5Str := dataStr[0:31]
		sha1Str := dataStr[33:73]
		sha256Str := dataStr[74:138]
		i := strings.Index(dataStr, "\n")
		dataStr = dataStr[i+1 : len(dataStr)]
		md5Array = append(md5Array, md5Str)
		sha1Aarray = append(sha1Aarray, sha1Str)
		sha256Array = append(sha256Array, sha256Str)
	}
	result.date = date
	result.md5 = md5Array
	result.sha1 = sha1Aarray
	result.sha256 = sha256Array
	return result
}

func worker(id int, jobs <- chan string, results chan <- malshareData){
	for j:= range jobs{
		fmt.Printf("worker %d start jobs http://malshare.com/daily/%s/malshare_fileList.%s.all.txt \n", id, j, j)
		results <- getHash(j)
		fmt.Printf("worker %d finished jobs http://malshare.com/daily/%s/malshare_fileList.%s.all.txt \n", id, j, j)
	}
}

func DumpData(){
	bodyStr, err := request.Request(URL)
	if err != nil{
		fmt.Println(err)
		return
	}
	jobs := make(chan string, LIMIT)
	results := make(chan malshareData, LIMIT)
	magic := regexp.MustCompile(`\"\[DIR\]\"></[a-z]{2}><[a-z]{2}><a\s[a-z]{4}=\"`)
	magicStr := string(magic.Find([]byte(bodyStr)))
	end := regexp.MustCompile("_[a-z]{8}/")
	outEnd := string(end.Find([]byte(bodyStr)))
	for w := 1 ; w < 101; w++{
		go worker(w, jobs, results)
	}
	for {
		i := strings.Index(bodyStr, magicStr)
		re := regexp.MustCompile("=\"\\d{4}-\\d{2}-\\d{2}")
		out := re.Find([]byte(bodyStr))
		if len(out) < 3{
			break
		}
		dateStr := string(out)[2:]
		bodyStr = bodyStr[i+len(magicStr)+1 : len(bodyStr)]
		if dateStr == outEnd{
			break
		}
		jobs <- dateStr
	}
	close(jobs)
	for a:=1; a <= LIMIT; a++{
		saveFile(<-results)
	}
}

func saveFile(data malshareData) {
	yyyy := data.date[0:4]
	mm := data.date[5:7]
	dd := data.date[8:10]
	yyyy_path := fmt.Sprintf("./%s", yyyy)
	mm_path := fmt.Sprintf("%s/%s", yyyy_path, mm)
	dd_path := fmt.Sprintf("%s/%s", mm_path, dd)
	if _, err := os.Stat(yyyy_path); os.IsNotExist(err) {
		os.Mkdir(yyyy_path, 777)
	}
	if _, err := os.Stat(mm_path); os.IsNotExist(err) {
		os.Mkdir(mm_path, 777)
	}
	if _, err := os.Stat(dd_path); os.IsNotExist(err) {
		os.Mkdir(dd_path, 777)
	}
	file, err := os.Create(fmt.Sprintf("%s/md5.txt", dd_path))
	if err != nil {
		return;
	}
	for i := 0; i < len(data.md5); i++ {
		file.WriteString(data.md5[i] + "\n")
	}
	file, err = os.Create(fmt.Sprintf("%s/sha1.txt", dd_path))
	if err != nil {
		return;
	}
	for i := 0; i < len(data.sha1); i++ {
		file.WriteString(data.sha1[i] + "\n")
	}
	file.Close()
	file, err = os.Create(fmt.Sprintf("%s/sha256.txt", dd_path))
	if err != nil {
		return;
	}
	for i := 0; i < len(data.sha256); i++ {
		file.WriteString(data.sha256[i] + "\n")
	}
	file.Close()
}

