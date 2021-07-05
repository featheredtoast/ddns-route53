package iplookup

import (
	"io/ioutil"
	"net/http"
)

type IpGetter struct{
	Server   string
}
func (i IpGetter) GetIp() (string, error) {
	resp, err := http.Get(i.Server)
	defer resp.Body.Close()
	if err != nil {
		return "", err
	}
	body, err := ioutil.ReadAll(resp.Body)
	return string(body), err
}
