package filter

import (
	"encoding/json"
	"net/http"

	"io/ioutil"

	"bytes"

	"github.com/labstack/echo"
)

var OpList = []string{
	"Phrase",
	"Prefix",
	"Regexp",
	"Term",
	"Wildcard",
	"DateRange",
	"NumericRange",
	"QueryString",
}

type Filter struct {
	Field  string   `json:"field,omitempty"`
	Op     string   `json:"op"`
	Values []string `json:"values"`
}

type Query struct {
	Id          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description,omitempty"`
	Filters     []Filter `json:"filters,omitempty"`
}

//跳转至首页
func (*Query) Index(c echo.Context) error {
	return c.Redirect(http.StatusMovedPermanently, "Filters")
}

//跳转至编辑页
func (*Query) Edit(c echo.Context) error {
	id := c.QueryParam("id")
	url := "editFilter" + id
	return c.Redirect(http.StatusMovedPermanently, url)
}

//跳转至新建页
func (*Query) New(c echo.Context) error {
	return c.Redirect(http.StatusMovedPermanently, "newFilter")
}

//获取所有query对象
func (*Query) GetQuerys(c echo.Context) error {
	querys, err := getQuerys()
	if err != nil {
		return error(err)
	}
	return c.JSON(http.StatusOK, querys)
}

func (*Query) GetQueryById(c echo.Context) error {
	id := c.QueryParam("id")
	query, err := GetQueryById(id)
	if err != nil {
		//沒有結果返回nil  页面显示为null
		return c.JSON(http.StatusOK, nil)
	}
	return c.JSON(http.StatusOK, query)
}

//创建query对象
func (*Query) Create(c echo.Context) error {
	query := new(Query)
	name := c.FormValue("name")
	description := c.FormValue("description")
	address := c.FormValue("address")
	severity := c.FormValue("severity")
	facility := c.FormValue("facility")
	message := c.FormValue("message")
	var messages []string
	if name == "" {
		messages = append(messages, "过滤名不能为空!")
	} else {
		querys, err := getQuerys()
		if err != nil {
			return error(err)
		}
		for _, q := range querys {
			if q.Name == name {
				messages = append(messages, "名称不能重复")
				break
			}
		}
	}

	query.Name = name
	query.Description = description
	if description == "" {
		messages = append(messages, "描述信息不能为空!")
	}
	if address != "" {
		var filter Filter
		filter.Field = "Address"
		filter.Op = "QueryString"
		filter.Values = append(filter.Values, address)
		query.Filters = append(query.Filters, filter)
	}
	if facility != "" {
		var filter Filter
		filter.Field = "Facility"
		filter.Op = "QueryString"
		filter.Values = append(filter.Values, facility)
		query.Filters = append(query.Filters, filter)
	}
	if severity != "" {
		var filter Filter
		filter.Field = "Severity"
		filter.Op = "QueryString"
		filter.Values = append(filter.Values, severity)
		query.Filters = append(query.Filters, filter)
	}
	if message != "" {
		var filter Filter
		filter.Field = "Message"
		filter.Op = "Phrase"
		filter.Values = append(filter.Values, message)
		query.Filters = append(query.Filters, filter)
	}
	if len(query.Filters) < 1 {
		messages = append(messages, "至少选择一个过滤条件")
	}
	if len(messages) != 0 {
		return c.JSON(http.StatusOK, messages)
	}

	q_err := createFilter(query)
	if q_err != nil {
		return error(q_err)
	}
	return c.JSON(http.StatusOK, messages)
}

//修改query对象
func (*Query) Update(c echo.Context) error {
	query := new(Query)
	id := c.FormValue("id")
	name := c.FormValue("name")
	description := c.FormValue("description")
	address := c.FormValue("address")
	severity := c.FormValue("severity")
	facility := c.FormValue("facility")
	message := c.FormValue("message")
	if id == "" {
		return nil
	}
	query.Id = id
	if name == "" {
		return c.Redirect(http.StatusMovedPermanently, "/UpdateFilter")
	}
	query.Name = name
	query.Description = description
	if address != "" {
		var filter Filter
		filter.Field = "Address"
		filter.Op = "QueryString"
		filter.Values = append(filter.Values, address)
		query.Filters = append(query.Filters, filter)
	}
	if facility != "" {
		var filter Filter
		filter.Field = "Facility"
		filter.Op = "QueryString"
		filter.Values = append(filter.Values, facility)
		query.Filters = append(query.Filters, filter)
	}
	if severity != "" {
		var filter Filter
		filter.Field = "Severity"
		filter.Op = "QueryString"
		filter.Values = append(filter.Values, severity)
		query.Filters = append(query.Filters, filter)
	}
	if message != "" {
		var filter Filter
		filter.Field = "Message"
		filter.Op = "Phrase"
		filter.Values = append(filter.Values, message)
		query.Filters = append(query.Filters, filter)
	}
	q_err := updateFilter(query)
	if q_err != nil {
		return error(q_err)
	}
	return c.Redirect(http.StatusMovedPermanently, "/Filters")
}

//删除query对象
func (*Query) Delete(c echo.Context) error {
	id := c.FormValue("id")
	err := deleteById(id)
	if err != nil {
		return error(err)
	}
	return c.Redirect(http.StatusMovedPermanently, "/Filters")
}

func getQuerys() ([]Query, error) {
	querys := []Query{}
	client := &http.Client{}
	request, req_err := http.NewRequest("GET", "http://127.0.0.1:8880/filters", nil)
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
		return nil, resp_err
	}
	j_err := json.Unmarshal(body, &querys)
	return querys, j_err
}

func GetQueryById(id string) (*Query, error) {
	query := new(Query)
	client := &http.Client{}
	request, req_err := http.NewRequest("GET", "http://127.0.0.1:8880/filters/"+id, nil)
	if req_err != nil {
		return query, req_err
	}

	request.Header.Set("Content-Type", "application/json")
	resp, client_err := client.Do(request)
	if client_err != nil {
		return query, client_err
	}
	body, resp_err := ioutil.ReadAll(resp.Body)
	if resp_err != nil {
		return query, resp_err
	}
	j_err := json.Unmarshal(body, &query)
	return query, j_err
}

func createFilter(query *Query) error {
	queryjson, json_err := json.Marshal(query)
	if json_err != nil {
		return json_err
	}
	client := &http.Client{}
	request, req_err := http.NewRequest("POST", "http://127.0.0.1:8880/filters", bytes.NewBuffer(queryjson))
	if req_err != nil {
		return req_err
	}
	request.Header.Set("Content-Type", "application/json")
	resp, client_err := client.Do(request)
	if client_err != nil {
		return client_err
	}
	_, resp_err := ioutil.ReadAll(resp.Body)
	return resp_err
}

func updateFilter(query *Query) error {
	url := "http://127.0.0.1:8880/filters/" + query.Id
	queryjson, json_err := json.Marshal(query)
	if json_err != nil {
		return json_err
	}
	client := &http.Client{}

	request, req_err := http.NewRequest("PUT", url, bytes.NewBuffer(queryjson))
	if req_err != nil {
		return req_err
	}
	request.Header.Set("Content-Type", "application/json")
	resp, client_err := client.Do(request)
	if client_err != nil {
		return client_err
	}
	_, resp_err := ioutil.ReadAll(resp.Body)
	return resp_err
}

func deleteById(id string) error {
	url := "http://127.0.0.1:8880/filters/" + id
	client := &http.Client{}

	request, req_err := http.NewRequest("DELETE", url, nil)
	if req_err != nil {
		return req_err
	}
	request.Header.Set("Content-Type", "application/json")
	resp, client_err := client.Do(request)
	if client_err != nil {
		return client_err
	}
	_, resp_err := ioutil.ReadAll(resp.Body)
	return resp_err
}

func asdasf() {
	//	var results = map[string]map[string]uint64{}
	//增加severity分组
	//var facilities = map[string]uint64{}
	//var severities = map[string]uint64{}
	//for i := 0; i < 8; i++ {
	//	query := bleve.NewMatchQuery("severity:" + string(i))
	//	searchRequest := bleve.NewSearchRequest(query)
	//	searchResult, _ := bleve.Index.Search(searchRequest)
	//	severities[string(i)] = searchResult.Total
	//}
	//	results["severities"] = severities
	//	//增加facility分组
	//	for i := 1; i <= 15; i++ {
	//		query := bleve.NewMatchQuery("facility:" + string(i))
	//		searchRequest := bleve.NewSearchRequest(query)
	//		searchResult, _ := bleve.Index.Search(searchRequest)
	//		facilities[string(i)] = searchResult.Total
	//	}
	//	query := bleve.NewNumericRangeQuery(&float64(16), &float64(24))
	//	searchRequest := bleve.NewSearchRequest(query)
	//	searchRequest.Fields = new["facility"]
	//	searchResult, _ := bleve.Index.Search(searchRequest)
	//	facilities["自定义"] = searchResult.Total
	//	results["facilities"] = facilities
}
