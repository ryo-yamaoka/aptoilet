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
	apiEndpoint = "https://pb52o89yja.execute-api.us-east-1.amazonaws.com/Prod/pir"
	timeLayout  = "2006/01/02 15:04:05"
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
	b, err := getToiletInfo()
	if err != nil {
		fmt.Printf("Failed to toilet API request: %s", err.Error())
		os.Exit(1)
	}

	var t []toiletInfo
	if err := json.Unmarshal(b, &t); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if len(t) == 0 {
		fmt.Println("Missing toilet information")
		os.Exit(1)
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

func getToiletInfo() ([]byte, error) {
	resp, err := http.Get(apiEndpoint)
	if err != nil {
		return nil, fmt.Errorf("Failed to API request: %s", err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("Unexpected response code: %d >= 400", resp.StatusCode)
	}

	return ioutil.ReadAll(resp.Body)
}
