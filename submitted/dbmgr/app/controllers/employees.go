package controllers

import (
	"dbmgr/app/libs"
	"dbmgr/app/models"
	"dbmgr/app/routes"

	"github.com/revel/revel"
	"github.com/runner-mei/orm"
)

// Employees - 控制器
type Employees struct {
	App
}

// 列出所有记录
func (c Employees) Index(pageIndex int, pageSize int) revel.Result {
	//var exprs []db.Expr
	//if "" != name {
	//  exprs = append(exprs, models.Employees.C.NAME.LIKE("%"+name+"%"))
	//}

	total, err := c.Lifecycle.DB.Employees().Where().Count()
	if err != nil {
		c.Flash.Error(err.Error())
		c.FlashParams()
		return c.Render(err)
	}

	if pageSize <= 0 {
		pageSize = libs.DEFAULT_SIZE_PER_PAGE
	}

	var employees []models.Employee
	err = c.Lifecycle.DB.Employees().Where().
		Offset(pageIndex * pageSize).
		Limit(pageSize).
		All(&employees)
	if err != nil {
		c.Flash.Error(err.Error())
		c.FlashParams()
		return c.Render()
	}
	paginator := libs.NewPaginator(c.Request.Request, pageSize, total)
	return c.Render(employees, paginator)
}

// 编辑新建记录
func (c Employees) New() revel.Result {
	return c.Render()
}

// 创建记录
func (c Employees) Create(employee *models.Employee) revel.Result {
	if employee.Validate(c.Validation) {
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect(routes.Employees.New())
	}

	_, err := c.Lifecycle.DB.Employees().Insert(employee)
	if err != nil {
		if oerr, ok := err.(*orm.Error); ok {
			for _, validation := range oerr.Validations {
				c.Validation.Error("该员工已存在，请检查姓名后再试").Key(models.KeyForEmployees(validation.Key))
			}
			c.Validation.Keep()
		}
		c.Flash.Error(err.Error())
		c.FlashParams()
		return c.Redirect(routes.Employees.New())
	}

	c.Flash.Success(revel.Message(c.Request.Locale, "员工"+employee.Name+"已成功添加"))
	return c.Redirect(routes.Employees.Index(0, 0))
}

// 编辑指定 id 的记录
func (c Employees) Edit(id int64) revel.Result {
	var employee models.Employee
	err := c.Lifecycle.DB.Employees().Id(id).Get(&employee)
	if err != nil {
		if err == orm.ErrNotFound {
			c.Flash.Error(revel.Message(c.Request.Locale, "对应员工未找到"))
		} else {
			c.Flash.Error(err.Error())
		}
		c.FlashParams()
		return c.Redirect(routes.Employees.Index(0, 0))
	}
	return c.Render(employee)
}

// 按 id 更新记录
func (c Employees) Update(employee *models.Employee) revel.Result {
	if employee.Validate(c.Validation) {
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect(routes.Employees.Edit(int64(employee.ID)))
	}

	err := c.Lifecycle.DB.Employees().Id(employee.ID).Update(employee)
	if err != nil {
		if err == orm.ErrNotFound {
			c.Flash.Error(revel.Message(c.Request.Locale, "对应员工未找到"))
		} else {
			if oerr, ok := err.(*orm.Error); ok {
				for _, validation := range oerr.Validations {
					c.Validation.Error(validation.Message).Key(models.KeyForEmployees(validation.Key))
				}
				c.Validation.Keep()
			}
			c.Flash.Error(err.Error())
		}
		c.FlashParams()
		return c.Redirect(routes.Employees.Edit(int64(employee.ID)))
	}
	c.Flash.Success(revel.Message(c.Request.Locale, "修改成功"))
	return c.Redirect(routes.Employees.Index(0, 0))
}

// 按 id 删除记录
func (c Employees) Delete(id int64) revel.Result {
	err := c.Lifecycle.DB.Employees().Id(id).Delete()
	if nil != err {
		if err == orm.ErrNotFound {
			c.Flash.Error(revel.Message(c.Request.Locale, "对应员工未找到"))
		} else {
			c.Flash.Error(err.Error())
		}
	} else {
		c.Flash.Success(revel.Message(c.Request.Locale, "删除成功"))
	}
	return c.Redirect(Employees.Index)
}

// 按 id 列表删除记录
func (c Employees) DeleteByIDs(id_list []int64) revel.Result {
	if len(id_list) == 0 {
		c.Flash.Error("请至少选择一条记录！")
		return c.Redirect(Employees.Index)
	}
	_, err := c.Lifecycle.DB.Employees().Where().And(orm.Cond{"id IN": id_list}).Delete()
	if nil != err {
		c.Flash.Error(err.Error())
	} else {
		c.Flash.Success(revel.Message(c.Request.Locale, "批量删除成功"))
	}
	return c.Redirect(Employees.Index)
}
