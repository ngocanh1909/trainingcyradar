package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

type malshareData struct{
	md5 []string
	sha1 []string
	sha256 []string
}

func request(url string,urlChan chan string) {
	resp, err := http.Get(url)
	if err != nil{
		urlChan <- ""
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil{
		urlChan <- ""
	}
	urlChan <- string(body)

}

func getHash (url string, malshareChan chan malshareData) {
	fmt.Println(url)
	urlChan := make(chan string)
	go request(url,urlChan)
	dataStr:= <- urlChan
	var result malshareData
	var md5Array []string
	var sha1Array []string
	var sha256Array []string
	for{
		if (len(dataStr) == 0){
			break
		}
		md5Str := dataStr[0:31]
		sha1Str := dataStr[33:73]
		sha256Str := dataStr[34:138]
		i := strings.Index(dataStr, "\n")
		dataStr = dataStr[i+1 : len(dataStr)]
		md5Array = append(md5Array, md5Str)
		sha1Array = append(sha1Array, sha1Str)
		sha256Array = append(sha256Array, sha256Str)
	}
	result.md5 = md5Array
	result.sha1 = sha1Array
	result.sha256= sha256Array
	malshareChan <- result
}

func dumpData() (map[string] malshareData, error){
	malshareMap := make(map[string] malshareData)
	resp, err := http.Get("https://malshare.com/daily/")
	if err != nil{
		return nil, err
	}
	defer resp.Body.Close()
	bodyUrl, err := ioutil.ReadAll(resp.Body)
	if err != nil{
		return nil, err
	}
	bodyStr:=string(bodyUrl)
	if len(bodyStr) == 0{
		return nil, errors.New("Can not dump data")
	}
	magic := regexp.MustCompile(`\"\[DIR\]\"></[a-z]{2}><[a-z]{2}><a\s[a-z]{4}=\"`)
	magicStr := string(magic.Find([]byte(bodyStr)))
	end := regexp.MustCompile("_[a-z]{8}/")
	outEnd := string(end.Find([]byte(bodyStr)))
	for{
		i := strings.Index(bodyStr, magicStr)
		if i == -1{
			return nil, errors.New("Can not find magic string")
		}
		re := regexp.MustCompile("=\"[0-9]{4}-\\d{2}-\\d{2}")
		out := re.Find([]byte (bodyStr))
		if len(out) == 0{
			return nil, errors.New("Can not find date string")
		}
		dateStr := string(out)[2:]
		bodyStr = bodyStr[i+len(magicStr)+1 : len(bodyStr)]
		if dateStr == outEnd{
			break
		}
		urlStr := fmt.Sprintf("https://malshare.com/daily/%s/malshare_fileList.%s.all.txt", dateStr, dateStr)
		fmt.Println(urlStr)
		done := make(chan malshareData)
		go getHash(urlStr,done)
		malshareMap[dateStr] = <-done
	}
	return malshareMap, nil
}
func orginizeData(malshareMap map[string] malshareData){
	for k := range malshareMap{
		yyyy := k[0:4]
		mm := k[5:7]
		dd := k[8:10]
		yyyyPath:=fmt.Sprintf("./%s",yyyy)
		mmPath:=fmt.Sprintf("%s/%s",yyyyPath,mm)
		ddPath:=fmt.Sprintf("%s/%s",mmPath,dd)
		if _, err := os.Stat(yyyyPath); os.IsNotExist(err){
			os.Mkdir(yyyyPath, 777)
		}
		if _, err := os.Stat(mmPath); os.IsNotExist(err){
			os.Mkdir(mmPath, 777)
		}
		if _, err := os.Stat(ddPath); os.IsNotExist(err){
			os.Mkdir(ddPath, 777)
		}
		file, err := os.Create(fmt.Sprintf("%s/md5.txt",ddPath))
		if err != nil{
			continue
		}
		for i:=0; i<len(malshareMap[k].md5); i++{
			file.WriteString(malshareMap[k].md5[i] )
			file.WriteString("\n")
		}
		file.Close()
		file, err = os.Create(fmt.Sprintf("%s/sha1.txt",ddPath))
		if err  != nil{
			continue
		}
		for i:=0; i<len(malshareMap[k].sha1) ; i++  {
			file.WriteString(malshareMap[k].sha1[i] )
			file.WriteString("\n")
		}
		file.Close()
		file, err =os.Create(fmt.Sprintf("%s/sha256.txt",ddPath))
		if err != nil{
			continue
		}
		for i:=0; i<len(malshareMap[k].sha256) ; i++  {
			file.WriteString(malshareMap[k].sha256[i] )
			file.WriteString("\n")
		}
		file.Close()
	}
}
func main(){
	go dumpData()
	time.Sleep(time.Hour)
	//orginizeData(data)
}
