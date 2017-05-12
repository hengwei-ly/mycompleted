package controllers

import (
	"dbmgr/app/models"
	"dbmgr/app/routes"

	"github.com/revel/revel"
)

type DevTables struct {
	App
}

type TabDatas struct {
	Datas   map[string]models.Tables
	Tabname string
}

// 列出所有记录
func (c DevTables) Index(page string) revel.Result {
	if page == "" {
		page = "alerts"
	}
	tables, err := models.GetAllTables()
	if err != nil {
		c.Flash.Error(err.Error())
		c.FlashParams()
		return c.Render()
	}

	tab := TabDatas{tables, page}
	return c.Render(tab)
}

func (c DevTables) DeleteByName(name, page string) revel.Result {
	err := models.DeleteByName(name)
	if err != nil {
		c.Flash.Error(err.Error(), "删除失败")
		c.FlashParams()
		c.Redirect(routes.DevTables.Index(page))
	}
	return c.Redirect(routes.DevTables.Index(page))
}
func (c DevTables) DeleteByNames(names []string, page string) revel.Result {
	revel.ERROR.Println(names)
	if len(names) <= 0 {
		c.Flash.Error("请选择至少一条记录")
		c.FlashParams()
		c.Redirect(routes.DevTables.Index(page))
	}
	for _, name := range names {
		err := models.DeleteByName(name)
		if err != nil {
			c.Flash.Error("删除 " + name + " 失败")
			c.FlashParams()
			c.Redirect(routes.DevTables.Index(page))
		}
	}
	return c.Redirect(routes.DevTables.Index(page))
}
