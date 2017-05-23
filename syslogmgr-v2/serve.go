package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"time"

	"strconv"

	"syslogmgr/public"

	"net/http/httputil"
	"net/url"

	"github.com/blevesearch/bleve"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type Syslog struct {
	Address   string    `json:"address"`
	Facility  int       `json:"facility"`
	Host      string    `json:"host"`
	Message   string    `json:"message"`
	Priority  int       `json:"priority"`
	Severity  int       `json:"severity"`
	Timestamp time.Time `json:"timestamp"`
	Level     string    `json:"level"`
}
type SyslogImpl struct {
}

var query *filter.Query

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Static("/public", "public")

	//e.GET("hengwei/aaa/Syslogs/GetCount", GetCount)
	//e.GET("hengwei/aaa/Syslogs/QueryJSON", QuerySyslogs)
	//e.GET("hengwei/aaa/SyslogFilters/GetQuerys", query.GetQuerys)
	//e.File("/Filters", "public/views/filters/index.html")
	//e.File("/editFilter", "public/views/filters/edit.html")

	//e.GET("/getQueryById", query.GetQueryById)

	//e.POST("/createFilter", query.Create)
	//e.POST("/updateFilter", query.Update)
	//e.POST("/deleteFilter", query.Delete)

	reverseProxy := httputil.NewSingleHostReverseProxy(&url.URL{Scheme: "http", Host: "127.0.0.1:8880"})
	proxy := echo.WrapHandler(reverseProxy)
	e.HTTPErrorHandler = func(err error, ctx echo.Context) {
		if err == echo.ErrNotFound {
			err = proxy(ctx)
		}
		if err != nil {
			e.DefaultHTTPErrorHandler(err, ctx)
		}
	}
	e.Logger.Fatal(e.Start(":9011"))
}

func GetCount(c echo.Context) error {
	filterName := c.QueryParam("filterName")
	SyslogImpl := new(SyslogImpl)
	searchRequest := SyslogImpl.GetSearchByFilterName(filterName)
	count, err := SyslogImpl.GetSyslogCounts(searchRequest)
	if err != nil {
		return err
	}
	return c.String(http.StatusOK, count)
}

//查询指定信息
func QuerySyslogs(c echo.Context) error {
	page := c.QueryParam("page")
	pagesize := c.QueryParam("pagesize")
	pagesizeInt, err := strconv.Atoi(pagesize)
	if err != nil {
		return err
	}
	if pagesizeInt <= 0 {
		return errors.New("pagesize is not right")
	}
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		return err
	}
	if pageInt <= 0 {
		return errors.New("page is not exists")
	}
	filterName := c.QueryParam("filterName")
	SyslogImpl := new(SyslogImpl)
	searchRequest := SyslogImpl.GetSearchByFilterName(filterName)
	searchRequest.Size = pagesizeInt
	searchRequest.From = (pageInt - 1) * pagesizeInt
	var syslogs []Syslog

	syslogs, err = SyslogImpl.GetSyslogs(searchRequest)
	if err != nil {
		return error(err)
	}
	for i := 0; i < len(syslogs); i++ {
		switch syslogs[i].Severity {
		case 0:
			syslogs[i].Level = "Emergency"
			break
		case 1:
			syslogs[i].Level = "Alert"
			break
		case 2:
			syslogs[i].Level = "Critical"
			break
		case 3:
			syslogs[i].Level = "Error"
			break
		case 4:
			syslogs[i].Level = "Warning"
			break
		case 5:
			syslogs[i].Level = "Notice"
			break
		case 6:
			syslogs[i].Level = "Informational"
			break
		case 7:
			syslogs[i].Level = "Debug"
		}
	}
	return c.JSON(http.StatusOK, syslogs)
}

//给定查询条件 返回结果
func (*SyslogImpl) GetSyslogs(searchRequest *bleve.SearchRequest) ([]Syslog, error) {
	pjson, json_err := json.Marshal(searchRequest)
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
		return nil, resp_err
	}
	syslogs := []Syslog{}
	json_err1 := json.Unmarshal(body, &syslogs)
	return syslogs, json_err1
}

//根据filterName获取记录数
func (*SyslogImpl) GetCountsByFilter(filterName string) (int, error) {
	SyslogImpl := new(SyslogImpl)
	searchRequest := SyslogImpl.GetSearchByFilterName(filterName)
	syslogs, err := SyslogImpl.GetSyslogs(searchRequest)
	return len(syslogs), err
}

//等待重构 传入一个filterName  获取对应Filter  返回一个searchRequrse
func (*SyslogImpl) GetSearchByFilterName(filterName string) *bleve.SearchRequest {
	var searchRequest *bleve.SearchRequest
	if filterName == "filter1" {
		query := bleve.NewMatchQuery("this")
		searchRequest = bleve.NewSearchRequest(query)
	} else if filterName == "filter2" {
		query := bleve.NewMatchQuery("asdasafffffffff")
		searchRequest = bleve.NewSearchRequest(query)
	} else if filterName == "filter3" {
		query := bleve.NewMatchQuery("filter5")
		searchRequest = bleve.NewSearchRequest(query)
	} else if filterName == "filter4" {
		query := bleve.NewMatchQuery("xiao")
		searchRequest = bleve.NewSearchRequest(query)
	} else if filterName == "filter5" {
		query := bleve.NewMatchQuery("da")
		searchRequest = bleve.NewSearchRequest(query)
	} else if filterName == "filter7" {
		query := bleve.NewMatchQuery("filter7")
		searchRequest = bleve.NewSearchRequest(query)
	} else {
		query := bleve.NewMatchAllQuery()
		searchRequest = bleve.NewSearchRequest(query)
	}
	return searchRequest
}

//根据searchRequest获取记录数
func (*SyslogImpl) GetSyslogCounts(searchRequest *bleve.SearchRequest) (string, error) {
	pjson, json_err := json.Marshal(searchRequest)
	if json_err != nil {
		return "0", json_err
	}
	client := &http.Client{}
	request, req_err := http.NewRequest("POST", "http://127.0.0.1:8880/summary", bytes.NewBuffer(pjson))
	if req_err != nil {
		return "0", req_err
	}

	request.Header.Set("Content-Type", "application/json")
	resp, client_err := client.Do(request)
	if client_err != nil {
		return "0", client_err
	}
	body, resp_err := ioutil.ReadAll(resp.Body)
	return string(body), resp_err
}
