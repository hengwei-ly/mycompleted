package main

import (
	"net/http"

	"syslogmgr/models"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Static("/", "public")

	e.File("/", "public/html/Index.html")

	e.GET("hengwei/aaa/Syslogs/GetCounts", GetCounts)
	e.GET("hengwei/aaa/Syslogs/QueryJSON", QuerySyslogs)

	e.Logger.Fatal(e.Start(":9000"))
}

func QuerySyslogs(c echo.Context) error {
	//Syslog := new(models.Syslog)
	//syslogs, err := Syslog.GetAllSyslogs()
	//if err != nil {
	//	return nil
	//}
	var syslogs []models.Syslog
	for i := 0; i < 256; i++ {
		syslog := models.Syslog{Address: "1234", Message: "asdhkjahskjh"}
		syslog.Priority = i
		syslogs = append(syslogs, syslog)
	}
	return c.JSON(http.StatusOK, syslogs)
}

func GetCounts(c echo.Context) error {
	return c.String(http.StatusOK, "256")
}
