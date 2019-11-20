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
		return "", err// neu get loi thi luu vet r exit
	}
	defer resp.Body.Close() // tri hoan dong body
	body, err := ioutil.ReadAll(resp.Body) // doc du lieu request body su dung
	if err != nil{
		return "", err //neu doc du lieu loi thi...

	}
	return string(body), err
}
//struct
type malshare_data struct {
	md5 []string
	sha1 []string
	sha256 []string
}
//ham lay md5, sha1, sha256 trong file
func get_hash(url string) malshare_data{
	var result malshare_data
	data_str, err := request(url) //gan data_str = du lieu request body
	if err != nil{
		return result
	}
	var md5_array [] string
	var sha1_array [] string
	var sha256_array [] string
	for{
		if (len(data_str)==0){ // doc du lieu request body = 0 thi dung
			break
		}
		md5_str := data_str[0:31] // do dai cua md5
		sha1_str := data_str[33:73] // do dai...
		sha256_str := data_str[74:138] // do dai...
		i := strings.Index(data_str, "\n") // du lieu co 4 cot tren 1 hang, tim moi cot tren tung hang
		data_str = data_str[i+1 : len(data_str) ] // cat tu vi tri i+1 den het do dai cua data_str
		md5_array = append(md5_array, md5_str) // them tung phan tu md5 vao mang
		sha1_array = append(sha1_array, sha1_str) // them...
		sha256_array = append(sha256_array, sha256_str) //them...

	}
	result.md5 = md5_array //tra ve map mang md5
	result.sha1 = sha1_array
	result.sha256 = sha256_array
	return result
}

func dump_data() map[string] malshare_data {
	malshare_map := make(map[string] malshare_data)
	body_str, err := request("https://malshare.com/daily/") // gan body_str = du lieu request body url
	//host_str := "https://malshare.com/daily/"
	if err != nil{
		return malshare_map
	}
	magic := regexp.MustCompile(`\"\[DIR\]\"></[a-z]{2}><[a-z]{2}><a\s[a-z]{4}=\"`)
	magic_str := string(magic.Find([]byte(body_str)))
	end := regexp.MustCompile("_[a-z]{8}/")
	out_end := string(end.Find([]byte(body_str)))
	dem := 0
	for {
		i := strings.Index(body_str, magic_str) // tim kiem doan magic trong chuoi string body url
		re := regexp.MustCompile("=\"[0-9]{4}-\\d{2}-\\d{2}")
		out := re.Find([]byte (body_str))
		date_str := string(out)[2:]
		body_str = body_str[i+len(magic_str)+1: len(body_str)]
		if date_str == out_end{
			break
		}

		url_str := fmt.Sprintf("https://malshare.com/daily/%s/malshare_fileList.%s.all.txt", date_str, date_str)
		//url_str:=strings.Join([]string{str1, body_str[i+len(magic_str) : i+len(magic_str)+10]},"")
		fmt.Println(url_str)
		hash_map := get_hash(url_str) // gan hash_map bang cac doan hash ham get_hash...
		malshare_map[date_str] = hash_map // luu cac doan hash trong mang hash_map vao mang malshare_map ngay thang...
		dem=dem+1
		if(dem>100){
			break
		}
	}
	return malshare_map
}
func orginize_data(malshare_map map[string] malshare_data){
	for k := range malshare_map { // dung for duyet mang malshare_map de lay cac gia tri date_str
		//fmt.Println(k)
		yyyy := k[0:4] //cat gia tri yyyy
		mm := k[5:7] //...
		dd := k[8:10] //...
		yyyy_path:=fmt.Sprintf("./%s",yyyy)
		mm_path:=fmt.Sprintf("%s/%s",yyyy_path,mm)
		dd_path:=fmt.Sprintf("%s/%s",mm_path,dd)
		//tao folder yyyy, neu chua ton tai
		if _, err := os.Stat(yyyy_path); os.IsNotExist(err) {
			os.Mkdir(yyyy_path,777)
		}
		//tao folder mm trong yyyy
		if _, err := os.Stat(mm_path); os.IsNotExist(err) {
			os.Mkdir(mm_path,777)
		}
		//tao folder dd trong yyyy/mm/
		if _, err := os.Stat(dd_path); os.IsNotExist(err) {
			os.Mkdir(dd_path,777)
		}
		//ghi file md5.txt trong folder yyyy/mm/dd
		file, err := os.Create(fmt.Sprintf("%s/md5.txt", dd_path))
		if err != nil{
			continue //neu co loi thi ghi sang ngay tiep theo
		}
		//duyet tung phan tu md5 trong mang malshare_map de ghi vao file
		for i:=0; i<len(malshare_map[k].md5); i++{ // lay tu 0 den < tong so luong phan tu trong mang
			file.WriteString(malshare_map[k].md5[i]+"\n")
		}
		file.Close() //dong file sau khi ghi
		file, err = os.Create(fmt.Sprintf("%s/sha1.txt", dd_path))
		if err != nil{
			continue
		}
		for i:=0; i<len(malshare_map[k].sha1); i++{
			file.WriteString(malshare_map[k].sha1[i]+"\n")
		}
		file.Close()
		file, err = os.Create(fmt.Sprintf("%s/sha256.txt", dd_path))
		if err != nil{
			continue
		}
		for i:=0; i<len(malshare_map[k].sha256); i++{
			file.WriteString(malshare_map[k].sha256[i]+"\n")
		}
		file.Close()
	}
}
func main()  {
	var data = dump_data();
	orginize_data(data)
}