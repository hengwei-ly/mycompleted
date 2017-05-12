package models

import (
	"time"

	"github.com/revel/revel"
)

type Employee struct {
	ID        int       `json:"id" xorm:"id pk autoincr"`
	Name      string    `json:"name" xorm:"name unique"`
	Age       int       `json:"age" xorm:"age"`
	Job       string    `json:"job" xorm:"job not null"`
	Email     string    `json:"email" xorm:"email"`
	CreatedAt time.Time `json:"created_at,omitempty" xorm:"created_at created"`
}

func (employee *Employee) TableName() string {
	return "tpt_employees"
}

func (employee *Employee) Validate(validation *revel.Validation) bool {

	validation.Required(employee.Name).Message("姓名不能为空").Key("employee.Name")

	validation.Range(employee.Age, 1, 100).Message("年龄必须为1-100的数字").Key("employee.Age")

	validation.Required(employee.Job).Message("职位不能为空").Key("employee.Job")

	validation.Email(employee.Email).Message("邮箱格式不正确").Key("employee.Email")
	return validation.HasErrors()
}

func KeyForEmployees(key string) string {
	switch key {
	case "id":
		return "employee.ID"
	case "name":
		return "employee.Name"
	case "Age":
		return "employee.Age"
	case "Job":
		return "employee.Job"
	case "Email":
		return "employee.Email"
	case "created_at":
		return "employee.CreatedAt"
	}
	return key
}
