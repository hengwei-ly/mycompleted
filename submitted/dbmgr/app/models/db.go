package models

import (
	"strings"

	"github.com/go-xorm/xorm"
	"github.com/runner-mei/orm"
)

type DB struct {
	Engine *xorm.Engine
}

func (db *DB) OnlineUsers() *orm.Collection {
	return orm.New(func() interface{} {
		return &OnlineUser{}
	})(db.Engine)
}

func (db *DB) AuthAccounts() *orm.Collection {
	return orm.New(func() interface{} {
		return &AuthAccount{}
	})(db.Engine)
}

func (db *DB) Employees() *orm.Collection {
	return orm.New(func() interface{} {
		return &Employee{}
	})(db.Engine)
}

func InitTables(engine *xorm.Engine) error {
	//添加结构体
	beans := []interface{}{
		&OnlineUser{},
		&AuthAccount{},
		&Employee{},
	}

	for _, bean := range beans {
		//创建表
		if err := engine.CreateTables(bean); err != nil {
			if !strings.Contains(err.Error(), "already exists") {
				return err
			}
		}
		//创建索引
		if err := engine.CreateIndexes(bean); err != nil {
			if !strings.Contains(err.Error(), "already exists") {
				return err
			}
		}
		//创建唯一性约束
		if err := engine.CreateUniques(bean); err != nil {
			if !strings.Contains(err.Error(), "already exists") {
				return err
			}
		}
	}
	return nil
}


