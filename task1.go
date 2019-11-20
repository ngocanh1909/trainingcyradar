package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
)
func request(url string) (string, error){
	resp, err := http.Get(url)
	if err != nil{
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil{
		return "", err

	}
	return string(body), err
}
type malshareData struct {
	md5 []string
	sha1 []string
	sha256 []string
}
func getHash(url string) malshareData{
	var result malshareData
	dataStr, err := request(url)
	if err != nil{
		return result
	}
	var md5Array [] string
	var sha1Array [] string
	var sha256Array [] string
	for{
		if (len(dataStr)==0){
			break
		}
		md5Str := dataStr[0:31]
		sha1Str := dataStr[33:73]
		sha256Str := dataStr[74:138]
		i := strings.Index(dataStr, "\n")
		dataStr = dataStr[i+1 : len(dataStr) ]
		md5Array = append(md5Array, md5Str)
		sha1Array = append(sha1Array, sha1Str)
		sha256Array = append(sha256Array, sha256Str)

	}
	result.md5 = md5Array
	result.sha1 = sha1Array
	result.sha256 = sha256Array
	return result
}

func dumpData() map[string] malshareData {
	malshareMap := make(map[string] malshareData)
	bodyStr, err := request("https://malshare.com/daily/")
	if err != nil{
		return malshareMap
	}
	magic := regexp.MustCompile(`\"\[DIR\]\"></[a-z]{2}><[a-z]{2}><a\s[a-z]{4}=\"`)
	magicStr := string(magic.Find([]byte(bodyStr)))
	end := regexp.MustCompile("_[a-z]{8}/")
	outEnd := string(end.Find([]byte(bodyStr)))
	dem := 0
	for {
		i := strings.Index(bodyStr, magicStr)
		re := regexp.MustCompile("=\"[0-9]{4}-\\d{2}-\\d{2}")
		out := re.Find([]byte (bodyStr))
		dateStr := string(out)[2:]
		bodyStr = bodyStr[i+len(magicStr)+1: len(bodyStr)]
		if dateStr == outEnd{
			break
		}
		url_str := fmt.Sprintf("https://malshare.com/daily/%s/malshare_fileList.%s.all.txt", dateStr, dateStr)
		fmt.Println(url_str)
		hash_map := getHash(url_str)
		malshareMap[dateStr] = hash_map
		dem=dem+1
		if(dem>100){
			break
		}
	}
	return malshareMap
}
func orginizeData(malshareMap map[string] malshareData){
	for k := range malshareMap {
		yyyy := k[0:4]
		mm := k[5:7]
		dd := k[8:10]
		yyyy_path:=fmt.Sprintf("./%s",yyyy)
		mm_path:=fmt.Sprintf("%s/%s",yyyy_path,mm)
		dd_path:=fmt.Sprintf("%s/%s",mm_path,dd)
		if _, err := os.Stat(yyyy_path); os.IsNotExist(err) {
			os.Mkdir(yyyy_path,777)
		}
		if _, err := os.Stat(mm_path); os.IsNotExist(err) {
			os.Mkdir(mm_path,777)
		}
		if _, err := os.Stat(dd_path); os.IsNotExist(err) {
			os.Mkdir(dd_path,777)
		}
		file, err := os.Create(fmt.Sprintf("%s/md5.txt", dd_path))
		if err != nil{
			continue
		}
		for i:=0; i<len(malshareMap[k].md5); i++{
			file.WriteString(malshareMap[k].md5[i]+"\n")
		}
		file.Close()
		file, err = os.Create(fmt.Sprintf("%s/sha1.txt", dd_path))
		if err != nil{
			continue
		}
		for i:=0; i<len(malshareMap[k].sha1); i++{
			file.WriteString(malshareMap[k].sha1[i]+"\n")
		}
		file.Close()
		file, err = os.Create(fmt.Sprintf("%s/sha256.txt", dd_path))
		if err != nil{
			continue
		}
		for i:=0; i<len(malshareMap[k].sha256); i++{
			file.WriteString(malshareMap[k].sha256[i]+"\n")
		}
		file.Close()
	}
}
func main()  {
	var data = dumpData()
	orginizeData(data)
}