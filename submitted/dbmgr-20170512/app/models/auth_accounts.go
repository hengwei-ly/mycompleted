package models

import (
	"time"

	"github.com/revel/revel"
)

type AuthAccount struct {
	ID          int       `json:"id" xorm:"id pk autoincr"`
	Name        string    `json:"name" xorm:"name unique"`
	Password    string    `json:"password,omitempty" xorm:"password"`
	Description string    `json:"description,omitempty" xorm:"description"`
	CreatedAt   time.Time `json:"created_at,omitempty" xorm:"created_at created"`
	UpdatedAt   time.Time `json:"updated_at,omitempty" xorm:"updated_at updated"`
}

func (authAccount *AuthAccount) TableName() string {
	return "tpt_auth_accounts"
}

func (authAccount *AuthAccount) Validate(validation *revel.Validation) bool {

	validation.Required(authAccount.Name).Key("authAccount.Name")

	validation.MinSize(authAccount.Password, 8).Key("authAccount.Password")

	validation.MaxSize(authAccount.Password, 250).Key("authAccount.Password")

	validation.MaxSize(authAccount.Description, 2000).Key("authAccount.Description")

	return validation.HasErrors()
}

func KeyForAuthAccounts(key string) string {
	switch key {
	case "id":
		return "authAccount.ID"
	case "name":
		return "authAccount.Name"
	case "password":
		return "authAccount.Password"
	case "description":
		return "authAccount.Description"
	case "created_at":
		return "authAccount.CreatedAt"
	case "updated_at":
		return "authAccount.UpdatedAt"
	}
	return key
}
