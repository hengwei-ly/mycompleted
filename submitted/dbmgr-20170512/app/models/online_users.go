package models

import (
	"time"

	"github.com/revel/revel"
)

type OnlineUser struct {
	ID            int       `json:"id" xorm:"id pk autoincr"`
	Name          string    `json:"name" xorm:"name unique"`
	AuthAccountID int64     `json:"auth_account_id" xorm:"auth_account_id"`
	Ipaddress     string    `json:"ipaddress,omitempty" xorm:"ipaddress"`
	Macaddress    string    `json:"macaddress,omitempty" xorm:"macaddress"`
	CreatedAt     time.Time `json:"created_at,omitempty" xorm:"created_at created"`
}

func (onlineUser *OnlineUser) TableName() string {
	return "tpt_online_users"
}

func (onlineUser *OnlineUser) Validate(validation *revel.Validation) bool {

	validation.Required(onlineUser.Name).Key("onlineUser.Name")

	validation.Required(onlineUser.AuthAccountID).Key("onlineUser.AuthAccountID")

	validation.MaxSize(onlineUser.Ipaddress, 200).Key("onlineUser.Ipaddress")

	validation.MaxSize(onlineUser.Macaddress, 200).Key("onlineUser.Macaddress")

	return validation.HasErrors()
}

func KeyForOnlineUsers(key string) string {
	switch key {
	case "id":
		return "onlineUser.ID"
	case "name":
		return "onlineUser.Name"
	case "auth_account_id":
		return "onlineUser.AuthAccountID"
	case "ipaddress":
		return "onlineUser.Ipaddress"
	case "macaddress":
		return "onlineUser.Macaddress"
	case "created_at":
		return "onlineUser.CreatedAt"
	}
	return key
}
