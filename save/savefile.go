package save

import (
	"fmt"
	"github.com/ngocanh1909/trainingcyradar/models"
	"os"
)

func SaveFile(data models.MalshareData) (error) {
	yyyy, mm, dd := data.Date.Date()
	yyyy_path := fmt.Sprintf("./%d", yyyy)
	mm_path := fmt.Sprintf("%s/%d", yyyy_path, mm)
	dd_path := fmt.Sprintf("%s/%d", mm_path, dd)
	if _, err := os.Stat(yyyy_path); os.IsNotExist(err) {
		err := os.Mkdir(yyyy_path, 0744)
		if err != nil {
			return err
		}
	}
	if _, err := os.Stat(mm_path); os.IsNotExist(err) {
		err := os.Mkdir(mm_path, 0744)
		if err != nil {
			return err
		}
	}
	if _, err := os.Stat(dd_path); os.IsNotExist(err) {
		err := os.Mkdir(dd_path, 0744)
		if err != nil {
			return err
		}
	}
	file, err := os.Create(fmt.Sprintf("%s/md5.txt", dd_path))
	if err != nil {
		return err
	}
	_, err = file.WriteString(data.Md5)
	if err != nil {
		return err
	}
	file, err = os.Create(fmt.Sprintf("%s/sha1.txt", dd_path))
	if err != nil {
		return err
	}
	_, err = file.WriteString(data.Sha1)
	if err != nil {
		return err
	}
	file.Close()
	file, err = os.Create(fmt.Sprintf("%s/sha256.txt", dd_path))
	if err != nil {
		return err
	}
	_, err = file.WriteString(data.Sha256)
	if err != nil {
		return err
	}
	file.Close()
	return nil
}
