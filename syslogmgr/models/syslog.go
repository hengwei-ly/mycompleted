package models

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/blevesearch/bleve"
)

type Syslog struct {
	Address   string    `json:"address"`
	Facility  int       `json:"facility"`
	Host      string    `json:"host"`
	Message   string    `json:"message"`
	Priority  int       `json:"priority"`
	Severity  int       `json:"severity"`
	Timestamp time.Time `json:"timestamp"`
}

//查询数据，
func (*Syslog) GetAllSyslogs() ([]Syslog, error) {
	query := bleve.NewMatchAllQuery()
	searchRequest := bleve.NewSearchRequest(query)
	searchRequest.Size = 20
	searchRequest.From = 0
	return GetSyslogs(searchRequest)
}

//给定查询条件 返回结果
func GetSyslogs(searchRequest *bleve.SearchRequest) ([]Syslog, error) {
	pjson, json_err := json.Marshal(searchRequest)
	fmt.Println(bytes.NewBuffer(pjson))
	if json_err != nil {
		return nil, json_err
	}

	client := &http.Client{}
	request, req_err := http.NewRequest("POST", "http://127.0.0.1:8880", bytes.NewBuffer(pjson))
	if req_err != nil {
		return nil, req_err
	}

	request.Header.Set("Content-Type", "application/json")
	resp, client_err := client.Do(request)
	if client_err != nil {
		return nil, client_err
	}
	body, resp_err := ioutil.ReadAll(resp.Body)

	if resp_err != nil {
		fmt.Println("Read resp error")
	}
	// fmt.Println(string(body))
	syslogs := []Syslog{}
	json_err1 := json.Unmarshal(body, &syslogs)
	return syslogs, json_err1
}
