package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
	"time"
)

const (
	timeLayout = "2006/01/02 15:04:05"
)

type toiletInfo struct {
	CreateAt int64  `json:"CreateAt"`
	UpdateAt int64  `json:"UpdateAt"`
	Light    int    `json:"Light"`
	Pir      int    `json:"Pir"`
	During   int    `json:"During"`
	SensorID string `json:"SensorId"`
}

func main() {
	endpoint := os.Getenv("AP_TOILET_ENDPOINT")
	if endpoint == "" {
		fmt.Println("Required set environment variable AP_TOILET_ENDPOINT")
		os.Exit(1)
	}

	resp, err := http.Get(endpoint)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	var t []toiletInfo
	if err := json.Unmarshal(b, &t); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	if len(t) == 0 {
		fmt.Println("Missing toilet information")
		return
	}

	sort.Slice(t, func(i, j int) bool {
		return t[i].UpdateAt > t[j].UpdateAt
	})

	usingStatus := "空き"
	if t[0].Pir == 1 {
		usingStatus = "使用中"
	}
	sensingTime := time.Unix(t[0].UpdateAt/1000, 0) // UpdateAt はミリ秒3桁が含まれるため下3桁を破棄する
	fmt.Printf("%s %s\n", sensingTime.Format(timeLayout), usingStatus)
}
