package controllers

import (
	"dbmgr/app/libs"
	"dbmgr/app/models"
	"dbmgr/app/routes"

	"github.com/revel/revel"
	"github.com/runner-mei/orm"
)

// AuthAccounts - 控制器
type AuthAccounts struct {
	App
}

// 列出所有记录
func (c AuthAccounts) Index(pageIndex int, pageSize int) revel.Result {
	//var exprs []db.Expr
	//if "" != name {
	//  exprs = append(exprs, models.AuthAccounts.C.NAME.LIKE("%"+name+"%"))
	//}

	total, err := c.Lifecycle.DB.AuthAccounts().Where().Count()
	if err != nil {
		c.Flash.Error(err.Error())
		c.FlashParams()
		return c.Render(err)
	}

	if pageSize <= 0 {
		pageSize = libs.DEFAULT_SIZE_PER_PAGE
	}

	var authAccounts []models.AuthAccount
	err = c.Lifecycle.DB.AuthAccounts().Where().
		Offset(pageIndex * pageSize).
		Limit(pageSize).
		All(&authAccounts)
	if err != nil {
		c.Flash.Error(err.Error())
		c.FlashParams()
		return c.Render()
	}
	paginator := libs.NewPaginator(c.Request.Request, pageSize, total)
	return c.Render(authAccounts, paginator)
}

// 编辑新建记录
func (c AuthAccounts) New() revel.Result {
	return c.Render()
}

// 创建记录
func (c AuthAccounts) Create(authAccount *models.AuthAccount) revel.Result {
	if authAccount.Validate(c.Validation) {
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect(routes.AuthAccounts.New())
	}

	_, err := c.Lifecycle.DB.AuthAccounts().Insert(authAccount)
	if err != nil {
		if oerr, ok := err.(*orm.Error); ok {
			for _, validation := range oerr.Validations {
				c.Validation.Error(validation.Message).Key(models.KeyForAuthAccounts(validation.Key))
			}
			c.Validation.Keep()
		}
		c.Flash.Error(err.Error())
		c.FlashParams()
		return c.Redirect(routes.AuthAccounts.New())
	}

	c.Flash.Success(revel.Message(c.Request.Locale, "insert.success"))
	return c.Redirect(routes.AuthAccounts.Index(0, 0))
}

// 编辑指定 id 的记录
func (c AuthAccounts) Edit(id int64) revel.Result {
	var authAccount models.AuthAccount
	err := c.Lifecycle.DB.AuthAccounts().Id(id).Get(&authAccount)
	if err != nil {
		if err == orm.ErrNotFound {
			c.Flash.Error(revel.Message(c.Request.Locale, "update.record_not_found"))
		} else {
			c.Flash.Error(err.Error())
		}
		c.FlashParams()
		return c.Redirect(routes.AuthAccounts.Index(0, 0))
	}
	return c.Render(authAccount)
}

// 按 id 更新记录
func (c AuthAccounts) Update(authAccount *models.AuthAccount) revel.Result {
	if authAccount.Validate(c.Validation) {
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect(routes.AuthAccounts.Edit(int64(authAccount.ID)))
	}

	err := c.Lifecycle.DB.AuthAccounts().Id(authAccount.ID).Update(authAccount)
	if err != nil {
		if err == orm.ErrNotFound {
			c.Flash.Error(revel.Message(c.Request.Locale, "update.record_not_found"))
		} else {
			if oerr, ok := err.(*orm.Error); ok {
				for _, validation := range oerr.Validations {
					c.Validation.Error(validation.Message).Key(models.KeyForAuthAccounts(validation.Key))
				}
				c.Validation.Keep()
			}
			c.Flash.Error(err.Error())
		}
		c.FlashParams()
		return c.Redirect(routes.AuthAccounts.Edit(int64(authAccount.ID)))
	}
	c.Flash.Success(revel.Message(c.Request.Locale, "update.success"))
	return c.Redirect(routes.AuthAccounts.Index(0, 0))
}

// 按 id 删除记录
func (c AuthAccounts) Delete(id int64) revel.Result {
	err := c.Lifecycle.DB.AuthAccounts().Id(id).Delete()
	if nil != err {
		if err == orm.ErrNotFound {
			c.Flash.Error(revel.Message(c.Request.Locale, "delete.record_not_found"))
		} else {
			c.Flash.Error(err.Error())
		}
	} else {
		c.Flash.Success(revel.Message(c.Request.Locale, "delete.success"))
	}
	return c.Redirect(AuthAccounts.Index)
}

// 按 id 列表删除记录
func (c AuthAccounts) DeleteByIDs(id_list []int64) revel.Result {
	if len(id_list) == 0 {
		c.Flash.Error("请至少选择一条记录！")
		return c.Redirect(AuthAccounts.Index)
	}
	_, err := c.Lifecycle.DB.AuthAccounts().Where().And(orm.Cond{"id IN": id_list}).Delete()
	if nil != err {
		c.Flash.Error(err.Error())
	} else {
		c.Flash.Success(revel.Message(c.Request.Locale, "delete.success"))
	}
	return c.Redirect(AuthAccounts.Index)
}
