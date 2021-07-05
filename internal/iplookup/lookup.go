package iplookup

import (
	"io/ioutil"
	"net"
	"net/http"
)

type IpGetter struct {
	Server    string
	Record    string
}

func (i *IpGetter) IpChanged() (bool, string, error) {
	desiredIp, err := i.GetIp()
	if err != nil {
		return false, "", err
	}
	ips, err := net.LookupIP(i.Record)
	if err != nil {
		return false, "", err
	}
	noMatch := true
	for _, ip := range ips {
		if desiredIp == ip.String() {
			noMatch = false
		}
	}
	return noMatch, desiredIp, nil
}

func (i *IpGetter) GetIp() (string, error) {
	resp, err := http.Get(i.Server)
	defer resp.Body.Close()
	if err != nil {
		return "", err
	}
	body, err := ioutil.ReadAll(resp.Body)
	return string(body), err
}
