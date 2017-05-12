package models

import (
	"strings"

	"github.com/go-xorm/xorm"
)

//自定义对象 对应文件信息
type DevTable struct {
	Name string
	//type 只在histories信息中用作判断月还是天，对应类型为 M D
	Type     string
	DateName string
	Size     string
}

var engine *xorm.Engine

// Employees - 控制器
type Tables []DevTable

func InitDevDB() {
	var err error
	engine, err = xorm.NewEngine("postgres", "user=tpt password=extreme dbname=tpt_data host=localhost port=35432 sslmode=disable")
	if err != nil {
		panic(err.Error())
	}
}

func GetAllTables() (map[string]Tables, error) {
	sql := `SELECT
		table_schema || '.' || table_name AS table_name,
		pg_size_pretty(pg_total_relation_size('"' || table_schema || '"."' || table_name || '"')) AS size
		FROM information_schema.tables
		WHERE
   			table_schema = 'public'
		ORDER BY
    			pg_total_relation_size('"' || table_schema || '"."' || table_name || '"') DESC `
	results, err := engine.Query(sql)
	if err != nil {
		return nil, err
	}

	var listAlertings Tables
	var listHistories Tables
	var listTraps Tables
	var listSyslogs Tables
	for _, value := range results {
		table := DevTable{
			Name: string(value["table_name"]),
			Size: string(value["size"]),
		}
		if strings.Contains(table.Name, "tpt_alert_histories_") {
			table.DateName = strings.TrimPrefix(table.Name, "public.tpt_alert_histories_")
			listAlertings = append(listAlertings, table)
		} else if strings.Contains(table.Name, "tpt_syslog_histories_") {
			table.DateName = strings.TrimPrefix(table.Name, "public.tpt_syslog_histories_")
			listSyslogs = append(listSyslogs, table)
		} else if strings.Contains(table.Name, "tpt_snmptrap_histories_") {
			table.DateName = strings.TrimPrefix(table.Name, "public.tpt_snmptrap_histories_")
			listTraps = append(listTraps, table)
		} else if strings.Contains(table.Name, "tpt_histories_month_") {
			table.DateName = strings.TrimPrefix(table.Name, "public.tpt_histories_month_")
			table.Type = "M"
			listHistories = append(listHistories, table)
		} else if strings.Contains(table.Name, "tpt_histories_raw_") {
			table.DateName = strings.TrimPrefix(table.Name, "public.tpt_histories_raw_")
			table.Type = "D"
			listHistories = append(listHistories, table)
		}
	}
	return map[string]Tables{"alerts": listAlertings,
		"histories": listHistories,
		"traps":     listTraps,
		"syslogs":   listSyslogs}, nil
}

func DeleteByName(name string) error {
	sql := "drop table " + name
	_, err := engine.Exec(sql)
	return err
}
