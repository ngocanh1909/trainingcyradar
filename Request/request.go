package Request

import (
	"io/ioutil"
	"net/http"
)

func Request(url string) (string,error)  {
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