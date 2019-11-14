package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)
func request(url string) string{
	resp, err := http.Get(url)
	if err != nil{
		log.Fatalln(err) // neu get loi thi luu vet r exit
	}
	defer resp.Body.Close() // tri hoan dong body
	body, err := ioutil.ReadAll(resp.Body) // doc du lieu request body su dung
	if err != nil{
		log.Fatalln(err) //neu doc du lieu loi thi...

	}
	return string(body)
}
//ham lay md5, sha1, sha256 trong file
func get_hash(url string) map[string] []string{
	data_str := request(url) //gan data_str = du lieu request body
	var result map[string] []string
	result = make(map[string] []string)
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
	result["md5"] = md5_array //tra ve map mang md5
	result["sha1"] = sha1_array
	result["sha256"] = sha256_array
	return result
}

func dump_data() map[string] map[string] []string {
	var malshare_map map[string] map[string] []string
	malshare_map = make(map[string] map[string] []string)
	body_str := request("https://malshare.com/daily/") // gan body_str = du lieu request body url
	host_str := "https://malshare.com/daily/"
	magic_str := "alt=\"[DIR]\"></td><td><a href=\"" //gan magic_str bang doan string ngay trc vi tri ngay thang nam
	end_str:="_disabled/" //string diem dung cua ngay thang nam
	//fmt.Println(get_hash("https://malshare.com/daily/2019-11-08/malshare_fileList.2019-11-08.all.txt"))
	//fmt.Println(body_str)
	//dem:=0
	for {
		i := strings.Index(body_str, magic_str) // tim kiem doan magic trong chuoi string body url
		//lay yyyy-mm-dd
		//cat vi tri ngay thang nam
		// gan date_str = do dai tu cuoi vi tri tim kiem (i+len) den het vi tri yyyy-mm-dd
		date_str := body_str[i+len(magic_str) : i+len(magic_str)+10]
		//fmt.Println(i)
		//lay yyyy-mm-dd tiep theo
		//cat tu sau yyyy-mm-dd dau tien den het body string
		//gan body_str = do dai tu cuoi vi tri ngay thang nam dc tim thay dau tien den het body string
		body_str = body_str[i+len(magic_str)+1: len(body_str)]
		//fmt.Println(body_str)
		if date_str == end_str{
			break //neu date_str toi diem dung thi beak
		}
		//yyyy := date_str[0:4]
		//mm := date_str[5:7]
		//dd := date_str[8:10]
		url_str := host_str + date_str + "/malshare_fileList." + date_str + ".all.txt"
		//url_str:=strings.Join([]string{str1, body_str[i+len(magic_str) : i+len(magic_str)+10]},"")
		fmt.Println(url_str)
		hash_map := get_hash(url_str) // gan hash_map bang cac doan hash ham get_hash...
		malshare_map[date_str] = hash_map // luu cac doan hash trong mang hash_map vao mang malshare_map ngay thang...
		//dem=dem+1
		//if(dem>10){
		//	break
		//}
	}
	return malshare_map
}
func orginize_data(malshare_map map[string] map[string] []string){
	for k := range malshare_map { // dung for duyet mang malshare_map de lay cac gia tri date_str
		//fmt.Println(k)
		yyyy := k[0:4] //cat gia tri yyyy
		mm := k[5:7] //...
		dd := k[8:10] //...
		//tao folder yyyy, neu chua ton tai
		if _, err := os.Stat("./" + yyyy); os.IsNotExist(err) {
			os.Mkdir("./" + yyyy,777)
		}
		//tao folder mm trong yyyy
		if _, err := os.Stat("./" + yyyy +"/"+ mm); os.IsNotExist(err) {
			os.Mkdir("./" + yyyy +"/"+ mm,777)
		}
		//tao folder dd trong yyyy/mm/
		if _, err := os.Stat("./" + yyyy +"/"+ mm + "/" + dd); os.IsNotExist(err) {
			os.Mkdir("./" + yyyy +"/"+ mm + "/" + dd,777)
		}
		//ghi file md5.txt trong folder yyyy/mm/dd
		file, err := os.Create("./" + yyyy +"/"+ mm + "/" + dd + "/md5.txt")
		if err != nil{
			continue //neu co loi thi ghi sang ngay tiep theo
		}
		//duyet tung phan tu md5 trong mang malshare_map de ghi vao file
		for i:=0; i<len(malshare_map[k]["md5"]); i++{ // lay tu 0 den < tong so luong phan tu trong mang
			file.WriteString(malshare_map[k]["md5"][i]+"\n")
		}
		file.Close() //dong file sau khi ghi
		file, err = os.Create("./" + yyyy +"/"+ mm + "/" + dd + "/sha1.txt")
		if err != nil{
			continue
		}
		for i:=0; i<len(malshare_map[k]["sha1"]); i++{
			file.WriteString(malshare_map[k]["sha1"][i]+"\n")
		}
		file.Close()
		file, err = os.Create("./" + yyyy +"/"+ mm + "/" + dd + "/sha256.txt")
		if err != nil{
			continue
		}
		for i:=0; i<len(malshare_map[k]["sha256"]); i++{
			file.WriteString(malshare_map[k]["sha256"][i]+"\n")
		}
		file.Close()
	}
}
func main()  {
	var data = dump_data();
	orginize_data(data)
}
