package save

import (
	"fmt"
	"github.com/ngocanh1909/trainingcyradar/config"
	"os"
)

func SaveFile(data config.MalshareData) {
	yyyy := data.Date.Format("2006")
	mm := data.Date.Format("01")
	dd := data.Date.Format("02")
	yyyy_path := fmt.Sprintf("./%s", yyyy)
	mm_path := fmt.Sprintf("%s/%s", yyyy_path, mm)
	dd_path := fmt.Sprintf("%s/%s", mm_path, dd)
	if _, err := os.Stat(yyyy_path); os.IsNotExist(err) {
		os.Mkdir(yyyy_path, 0744)
	}
	if _, err := os.Stat(mm_path); os.IsNotExist(err) {
		os.Mkdir(mm_path, 0744)
	}
	if _, err := os.Stat(dd_path); os.IsNotExist(err) {
		os.Mkdir(dd_path, 0744)
	}
	file, err := os.Create(fmt.Sprintf("%s/md5.txt", dd_path))
	if err != nil {
		return;
	}
	for i := 0; i < len(data.Md5); i++ {
		file.WriteString(data.Md5)
	}
	file, err = os.Create(fmt.Sprintf("%s/sha1.txt", dd_path))
	if err != nil {
		return;
	}
	for i := 0; i < len(data.Sha1); i++ {
		file.WriteString(data.Sha1)
	}
	file.Close()
	file, err = os.Create(fmt.Sprintf("%s/sha256.txt", dd_path))
	if err != nil {
		return;
	}
	for i := 0; i < len(data.Sha256); i++ {
		file.WriteString(data.Sha256)
	}
	file.Close()
}
