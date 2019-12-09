package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
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
func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	fmt.Printf("%s took %s\n", name, elapsed)
}
func offBit(val uint32, k_th int) uint32 {
	var result uint32
	var operator uint32
	operator = 1 << k_th
	operator = ^operator
	result = val & operator
	return result
}
func ipv4ToBit(ip string) uint32 {
	parts := strings.Split(ip, ".")
	var val32 uint32
	if len(parts) < 4 {
		fmt.Printf("This is NOT a valid IPv4 address")
	} else if len(parts) == 4 {
		for i := 0; i < len(parts); i++ {
			if j, err := strconv.Atoi(parts[i]); err == nil {
				if j < 0 || j > 255 {
					fmt.Printf("This is NOT a valid IPv4 address")
				}
				val32 = val32*256 + uint32(j)
			}
		}
	}
	return val32
}
func handleIPRange(prefixes string) (uint32, uint32) {
	iprangeIP := prefixes[0:strings.Index(prefixes, "/")]
	val32 := ipv4ToBit(iprangeIP)
	nb := prefixes[ strings.Index(prefixes, "/")+1:]
	nbit, err := strconv.Atoi(nb)
	if err != nil {
		fmt.Printf("%s", err)
	}
	var numberAnd uint32
	numberAnd = ^numberAnd
	for  i := 0; i<32-nbit; i++{
		numberAnd = offBit(numberAnd,i)
	}
	return val32, numberAnd
}
func handleIP(ip string, ipRange uint32, numberAnd uint32) uint32 {
	val32 := ipv4ToBit(ip)
	val32 = val32 & numberAnd
	return val32 ^ ipRange
}
func main() {
	var parent IPRange
	defer timeTrack(time.Now(), fmt.Sprintf("Compare IP"))
	read, err := ioutil.ReadFile("ip-ranges.json")
	if err != nil {
		return
	}
	myJson := string(read)
	json.Unmarshal([]byte(myJson), &parent)
	var line [] string
	for i := 0; i < len(parent.Prefixes); i++ {
		line = append(line, parent.Prefixes[i].IPPrefix)
	}
	line = append(line, "10.0.0.0/8")
	line = append(line, "172.16.0.0/12")
	line = append(line, "192.168.0.0/16")
	var iprangeArr [] uint32
	var numberAndArr [] uint32
	for i := 0; i < len(line); i++ {
		iprange, so_de_and := handleIPRange(line[i])
		iprangeArr = append(iprangeArr, iprange)
		numberAndArr = append(numberAndArr, so_de_and)
	}
	fmt.Printf("%b ", numberAndArr)
	file, err := os.Open("all.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var count = 0
		ip := scanner.Text()
		for i := 0; i < len(iprangeArr); i++ {
			if (handleIP(ip, iprangeArr[i], numberAndArr[i]) == 0) {
				count++
			}
		}
		if (count == 0) {
			fmt.Printf("%s\n", ip)
			break
		}
	}
}
