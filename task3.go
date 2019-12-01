package main
import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)
type IPRange struct {
	SyncToken    string `json:"syncToken"`
	CreateDate   string `json:"createDate"`
	Prefixes     [] Prefixes
	Ipv6Prefixes [] Ipv6Prefixes
}
type Prefixes struct {
	IPPrefix string `json:"ip_prefix"`
	Region   string `json:"region"`
	Service  string `json:"service"`
}
type Ipv6Prefixes struct {
	Ipv6Prefix string `json:"ipv6_prefix"`
	Region     string `json:"region"`
	Service    string `json:"service"`
}
func TimeTrack(start time.Time, name string){
	elapsed := time.Since(start)
	fmt.Printf("%s took %s\n", name, elapsed)
}

func tatBit( val uint32 ,k_th int ) uint32{
	var ketqua uint32
	var toan_tu uint32
	toan_tu = 1 << k_th
	toan_tu = ^toan_tu
	ketqua= val  & toan_tu
	return ketqua
}
func xulidaiIP(Prefixes string) (uint32,int){
	//Doc cai IP trong 1 dai IP
	iprangeIP := Prefixes[0 : strings.Index(Prefixes, "/")]
	ip := net.ParseIP(iprangeIP)
	val32 := uint32(ip[12])
	val32 = val32*256 + uint32(ip[13])
	val32 = val32*256 + uint32(ip[14])
	val32 = val32*256 + uint32(ip[15])
	//Doc so luong bit dung trong  1 dai IP
	nb := Prefixes[ strings.Index(Prefixes, "/")+1: ]
	nbit,err:= strconv.Atoi(nb)
	if (err!=nil){

	}
	return val32,nbit
}
func xuliIP(ip string,ip_range uint32,nbit int) uint32{
	ipp := net.ParseIP(ip)
	val32 := uint32(ipp[12])
	val32 = val32*256+uint32(ipp[13])
	val32 = val32*256+uint32(ipp[14])
	val32 = val32*256+uint32(ipp[15])
	//fmt.Printf("%d %b\n ",nbit,bang_tat_bit[32-nbit])
	val32=val32 & bang_tat_bit[32-nbit]

	return val32 ^  ip_range
}
var bang_tat_bit[] uint32
func main() {
	var parent IPRange
	defer TimeTrack(time.Now(), fmt.Sprintf("Compare IP"))
	read, err := ioutil.ReadFile("ip-ranges.json")
	if err != nil {
		return
	}
	myJson := string(read)
	var bangtatbit uint32
	bangtatbit = 0xFFFFFFFF// 11111111111111111111111111
	bang_tat_bit=append(bang_tat_bit,0xFFFFFFFF)
	for i:=0 ; i<32;i++{
		bangtatbit=tatBit(bangtatbit,i)
		bang_tat_bit=append(bang_tat_bit,bangtatbit)
	}
	bang_tat_bit=append(bang_tat_bit,0)
	json.Unmarshal([]byte(myJson), &parent)
	var line[] string
	for i:=0; i<len(parent.Prefixes); i++{
		line = append(line, parent.Prefixes[i].IPPrefix)
	}
	line = append(line, "10.0.0.0/8")
	line = append(line, "172.16.0.0/12")
	line = append(line, "192.168.0.0/16")

	var iprange_arr [] uint32
	var iprange_nbit [] int
	for i:=0; i<len(line); i++{
		iprange,nbit:= xulidaiIP(line[i])
		iprange_arr = append(iprange_arr, iprange)
		iprange_nbit = append(iprange_nbit,nbit)
	}
	file, err := os.Open("all.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		var d=0
		ip := scanner.Text()
		for i:=0; i<len(iprange_arr); i++{
			if(xuliIP(ip, iprange_arr[i],iprange_nbit[i])==0){
				d++
			}
		}
		if (d==0){
			fmt.Printf("%s\n",ip)
			break
		}

	}
}