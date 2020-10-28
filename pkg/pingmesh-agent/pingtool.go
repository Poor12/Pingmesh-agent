package pingmesh_agent

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"k8s.io/klog/v2"
	"net"
	"net/url"
)

type PingTool struct {
	storage	*Storage
}

type pinglist struct {
	WorkerName           string
	Patition			 string
	PingList             map[string][]string
}

func NewPingtool(store *Storage) *PingTool {
	return &PingTool{
		storage: 	store,

	}
}

func (c* PingTool) getPinglist(url string) ([]byte,error){
	client := NewHTTPClient()
	req, err := BuildRequest("GET", url, nil, "", "")
	if err != nil {
		return nil, err
	}

	res, err := SendRequest(req, client)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	pingList, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return pingList, nil
}

func (c* PingTool) uploadData(url string,body io.Reader) ([]byte, error) {
	 client := NewHTTPClient()
	 req, err := BuildRequest("POST",url,body,"","")
	 if err != nil{
	 	return nil, err
	 }

	 res, err := SendRequest(req,client)
	 if err != nil{
	 	return nil, err
	 }
	 defer res.Body.Close()

	 isSuccess, err := ioutil.ReadAll(res.Body)
	 if err != nil{
	 	return nil, err
	 }
	 return isSuccess, nil
}

func (c *PingTool) Ping(baseCtx context.Context) error{
	daddr := (&url.URL{
		Scheme: "http",
		Host:   net.JoinHostPort(DefaultPingmeshServiceAddr, DefaultPingmeshServerPort),
	}).String() + DefaultPingmeshDownloadURL
	pl := &pinglist{}
	pingList, err := c.getPinglist(daddr)
	klog.Info("get pinglist: ",pingList)
	if err != nil{
		pl = c.storage.pinglist
		if pl == nil {
			return fmt.Errorf("unable to get pinglist from pingmesh-server")
		}
	} else{
		err = json.Unmarshal(pingList, &pl)
		if err != nil {
			return fmt.Errorf("Umarshal failed:", err)
		}
		c.storage.pinglist = pl
	}
	pros := ProbeICMPs(*pl)
	uaddr := (&url.URL{
		Scheme: "http",
		Host: net.JoinHostPort(DefaultPingmeshServiceAddr,DefaultPingmeshServerPort),
	}).String() + DefaultPingmeshUploadURL
	klog.Info(pros)
	for _, pro := range pros {
		for _, metrics := range pro {
			mjson, _ := json.Marshal(metrics)
			data := bytes.NewReader(mjson)
			isSuccess, err := c.uploadData(uaddr, data)
			if err != nil || string(isSuccess) == "fail" {
				return fmt.Errorf("unable to upload metrics to pingmesh-server")
			}
		}
	}
	return nil
}

