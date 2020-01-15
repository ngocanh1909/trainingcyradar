package request

import (
	"io/ioutil"
	"net/http"
)

func Request(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	//fmt.Printf("Error %d\n", resp.StatusCode)
	if resp.StatusCode == 404 {
		return "", err
	}
	return string(body), err
}
