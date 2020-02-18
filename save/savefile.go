package save

import (
	"fmt"
	"github.com/ngocanh1909/trainingcyradar/models"
	"os"
)

func SaveFile(data models.MalshareData) (error) {
	yyyy, mm, dd := data.Date.Date()
	yyyyPath := fmt.Sprintf("./%d", yyyy)
	mmPath := fmt.Sprintf("%s/%d", yyyyPath, mm)
	ddPath := fmt.Sprintf("%s/%d", mmPath, dd)
	if _, err := os.Stat(yyyyPath); os.IsNotExist(err) {
		err := os.Mkdir(yyyyPath, 0744)
		if err != nil {
			return err
		}
	}
	if _, err := os.Stat(mmPath); os.IsNotExist(err) {
		err := os.Mkdir(mmPath, 0744)
		if err != nil {
			return err
		}
	}
	if _, err := os.Stat(ddPath); os.IsNotExist(err) {
		err := os.Mkdir(ddPath, 0744)
		if err != nil {
			return err
		}
	}
	file, err := os.Create(fmt.Sprintf("%s/md5.txt", ddPath))
	if err != nil {
		return err
	}
	_, err = file.WriteString(data.Md5)
	if err != nil {
		return err
	}
	file, err = os.Create(fmt.Sprintf("%s/sha1.txt", ddPath))
	if err != nil {
		return err
	}
	_, err = file.WriteString(data.Sha1)
	if err != nil {
		return err
	}
	file.Close()
	file, err = os.Create(fmt.Sprintf("%s/sha256.txt", ddPath))
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
