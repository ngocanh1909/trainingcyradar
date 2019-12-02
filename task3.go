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
	Ipv6Prefixes [] ipv6Prefixes
}
type Prefixes struct {
	IPPrefix string `json:"ip_prefix"`
	Region   string `json:"region"`
	Service  string `json:"service"`
}
type ipv6Prefixes struct {
	Ipv6Prefix string `json:"ipv6_prefix"`
	Region     string `json:"region"`
	Service    string `json:"service"`
}
func timeTrack(start time.Time, name string){
	elapsed := time.Since(start)
	fmt.Printf("%s took %s\n", name, elapsed)
}
func offBit( val uint32 ,k_th int ) uint32{
	var result uint32
	var operator uint32
	operator = 1 << k_th
	operator = ^operator
	result = val  & operator
	return result
}
func handleIPRange(Prefixes string) (uint32,int){
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
	if err != nil{
		fmt.Printf("%s",err)
	}
	return val32,nbit
}
func handleIP(ip string,ipRange uint32,nbit int) uint32{
	ipp := net.ParseIP(ip)
	val32 := uint32(ipp[12])
	val32 = val32*256+uint32(ipp[13])
	val32 = val32*256+uint32(ipp[14])
	val32 = val32*256+uint32(ipp[15])
	//fmt.Printf("%d %b\n ",nbit,bang_tat_bit[32-nbit])
	val32=val32 & tableOffBit[32-nbit]
	return val32 ^  ipRange
}
var tableOffBit[] uint32
func main() {
	var parent IPRange
	defer timeTrack(time.Now(), fmt.Sprintf("Compare IP"))
	read, err := ioutil.ReadFile("ip-ranges.json")
	if err != nil {
		return
	}
	myJson := string(read)
	var bitOff uint32
	bitOff = 0xFFFFFFFF //11111111111111111111111111
	tableOffBit=append(tableOffBit,0xFFFFFFFF)
	for i:=0 ; i<32;i++{
		bitOff=offBit(bitOff,i)
		tableOffBit=append(tableOffBit,bitOff)
	}
	tableOffBit=append(tableOffBit,0)
	json.Unmarshal([]byte(myJson), &parent)
	var line[] string
	for i:=0; i<len(parent.Prefixes); i++{
		line = append(line, parent.Prefixes[i].IPPrefix)
	}
	line = append(line, "10.0.0.0/8")
	line = append(line, "172.16.0.0/12")
	line = append(line, "192.168.0.0/16")
	var iprangeArr [] uint32
	var iprangeNBit [] int
	for i:=0; i<len(line); i++{
		iprange,nbit:= handleIPRange(line[i])
		iprangeArr = append(iprangeArr, iprange)
		iprangeNBit = append(iprangeNBit,nbit)
	}
	file, err := os.Open("all.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var count=0
		ip := scanner.Text()
		for i:=0; i<len(iprangeArr); i++{
			if(handleIP(ip, iprangeArr[i],iprangeNBit[i])==0){
				count++
			}
		}
		if (count==0){
			fmt.Printf("%s\n",ip)
			break
		}

	}
}