package controllers

import (
	"dbmgr/app/libs"
	"dbmgr/app/models"
	"dbmgr/app/routes"

	"github.com/revel/revel"
	"github.com/runner-mei/orm"
)

// OnlineUsers - 控制器
type OnlineUsers struct {
	App
}

// 列出所有记录
func (c OnlineUsers) Index(pageIndex int, pageSize int) revel.Result {
	//var exprs []db.Expr
	//if "" != name {
	//  exprs = append(exprs, models.OnlineUsers.C.NAME.LIKE("%"+name+"%"))
	//}

	total, err := c.Lifecycle.DB.OnlineUsers().Where().Count()
	if err != nil {
		c.Flash.Error(err.Error())
		c.FlashParams()
		return c.Render(err)
	}

	if pageSize <= 0 {
		pageSize = libs.DEFAULT_SIZE_PER_PAGE
	}

	var onlineUsers []models.OnlineUser
	err = c.Lifecycle.DB.OnlineUsers().Where().
		Offset(pageIndex * pageSize).
		Limit(pageSize).
		All(&onlineUsers)
	if err != nil {
		c.Flash.Error(err.Error())
		c.FlashParams()
		return c.Render()
	}
	paginator := libs.NewPaginator(c.Request.Request, pageSize, total)
	return c.Render(onlineUsers, paginator)
}

// 编辑新建记录
func (c OnlineUsers) New() revel.Result {
	return c.Render()
}

// 创建记录
func (c OnlineUsers) Create(onlineUser *models.OnlineUser) revel.Result {
	if onlineUser.Validate(c.Validation) {
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect(routes.OnlineUsers.New())
	}

	_, err := c.Lifecycle.DB.OnlineUsers().Insert(onlineUser)
	if err != nil {
		if oerr, ok := err.(*orm.Error); ok {
			for _, validation := range oerr.Validations {
				c.Validation.Error(validation.Message).Key(models.KeyForOnlineUsers(validation.Key))
			}
			c.Validation.Keep()
		}
		c.Flash.Error(err.Error())
		c.FlashParams()
		return c.Redirect(routes.OnlineUsers.New())
	}

	c.Flash.Success(revel.Message(c.Request.Locale, "insert.success"))
	return c.Redirect(routes.OnlineUsers.Index(0, 0))
}

// 编辑指定 id 的记录
func (c OnlineUsers) Edit(id int64) revel.Result {
	var onlineUser models.OnlineUser
	err := c.Lifecycle.DB.OnlineUsers().Id(id).Get(&onlineUser)
	if err != nil {
		if err == orm.ErrNotFound {
			c.Flash.Error(revel.Message(c.Request.Locale, "update.record_not_found"))
		} else {
			c.Flash.Error(err.Error())
		}
		c.FlashParams()
		return c.Redirect(routes.OnlineUsers.Index(0, 0))
	}
	return c.Render(onlineUser)
}

// 按 id 更新记录
func (c OnlineUsers) Update(onlineUser *models.OnlineUser) revel.Result {
	if onlineUser.Validate(c.Validation) {
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect(routes.OnlineUsers.Edit(int64(onlineUser.ID)))
	}

	err := c.Lifecycle.DB.OnlineUsers().Id(onlineUser.ID).Update(onlineUser)
	if err != nil {
		if err == orm.ErrNotFound {
			c.Flash.Error(revel.Message(c.Request.Locale, "update.record_not_found"))
		} else {
			if oerr, ok := err.(*orm.Error); ok {
				for _, validation := range oerr.Validations {
					c.Validation.Error(validation.Message).Key(models.KeyForOnlineUsers(validation.Key))
				}
				c.Validation.Keep()
			}
			c.Flash.Error(err.Error())
		}
		c.FlashParams()
		return c.Redirect(routes.OnlineUsers.Edit(int64(onlineUser.ID)))
	}
	c.Flash.Success(revel.Message(c.Request.Locale, "update.success"))
	return c.Redirect(routes.OnlineUsers.Index(0, 0))
}

// 按 id 删除记录
func (c OnlineUsers) Delete(id int64) revel.Result {
	err := c.Lifecycle.DB.OnlineUsers().Id(id).Delete()
	if nil != err {
		if err == orm.ErrNotFound {
			c.Flash.Error(revel.Message(c.Request.Locale, "delete.record_not_found"))
		} else {
			c.Flash.Error(err.Error())
		}
	} else {
		c.Flash.Success(revel.Message(c.Request.Locale, "delete.success"))
	}
	return c.Redirect(OnlineUsers.Index)
}

// 按 id 列表删除记录
func (c OnlineUsers) DeleteByIDs(id_list []int64) revel.Result {
	if len(id_list) == 0 {
		c.Flash.Error("请至少选择一条记录！")
		return c.Redirect(OnlineUsers.Index)
	}
	_, err := c.Lifecycle.DB.OnlineUsers().Where().And(orm.Cond{"id IN": id_list}).Delete()
	if nil != err {
		c.Flash.Error(err.Error())
	} else {
		c.Flash.Success(revel.Message(c.Request.Locale, "delete.success"))
	}
	return c.Redirect(OnlineUsers.Index)
}
