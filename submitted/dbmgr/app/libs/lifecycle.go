package libs

import (
	"cn/com/hengwei/commons"
	"cn/com/hengwei/commons/types"
	"database/sql"

	"dbmgr/app/models"

	"github.com/Masterminds/squirrel"
	"github.com/go-xorm/xorm"
	_ "github.com/lib/pq"
	"github.com/three-plus-three/sso/client/revel_sso"
)

// DbRunner wraps sql.DB to implement Runner.
type DbRunner struct {
	*sql.DB
}

// QueryRow wraps QueryRow to implement RowScanner.
func (r DbRunner) QueryRow(query string, args ...interface{}) squirrel.RowScanner {
	return r.DB.QueryRow(query, args...)
}

// TxRunner wraps sql.Tx to implement Runner.
type TxRunner struct {
	*sql.Tx
}

// QueryRow wraps QueryRow to implement RowScanner.
func (r TxRunner) QueryRow(query string, args ...interface{}) squirrel.RowScanner {
	return r.Tx.QueryRow(query, args...)
}

// Lifecycle 表示一个运行周期，它包含了所有业务相关的对象
type Lifecycle struct {
	commons.Base
	Env         *commons.Environment
	Definitions *types.TableDefinitions
	DB          models.DB
	DataDB      models.DB
	//MoCache     MoCache
	Variables map[string]interface{}
	UrlPrefix string
	CheckUser revel_sso.CheckFunc

	//timer atomic.Value
}

// // Close 结束一个生命周期
// func (lifecycle *Lifecycle) Close() error {
// 	return lifecycle.CloseWith(func() error {
// 		if o := lifecycle.timer.Load(); o != nil {
// 			if timer, ok := o.(*time.Timer); ok {
// 				timer.Stop()
// 			}
// 		}
// 		return nil
// 	})
// }

// func (lifecycle *Lifecycle) runLoop() {
//  defer lifecycle.CatchThrow("runLoop:", nil)
//
//  tick := time.NewTicker(60 * time.Second)
//  defer tick.Stop()
//
//  for {
//    select {
//    case <-lifecycle.S:
//      return
//    case <-tick.C:
//      if err := lifecycle.MoCache.Refresh(); err != nil {
//        log.Println(err)
//      }
//    }
//  }
// }
//
// func (lifecycle *Lifecycle) refresh() {
// 	defer lifecycle.CatchThrow("Refresh:", nil)
// 	if lifecycle.IsClosed() {
// 		return
// 	}
//
// 	defer lifecycle.timer.Store(time.AfterFunc(60*time.Second, lifecycle.refresh))
//
// 	if err := lifecycle.MoCache.Refresh(); err != nil {
// 		log.Println(err)
// 	}
// }

// NewLifecycle 创建一个生命周期
func NewLifecycle(env *commons.Environment,
	definitions *types.TableDefinitions) (*Lifecycle, error) {

	dbDrv, dbURL := env.Db.Models.Url()
	engine, err := xorm.NewEngine(dbDrv, dbURL)
	if err != nil {
		return nil, err
	}

	lifecycle := &Lifecycle{
		Env:         env,
		Definitions: definitions,
		DB:          models.DB{Engine: engine},
	}
	/*
		lifecycle.MoCache.lifecycle = lifecycle

		lifecycle.MoCache.managedElement = definitions.FindByUnderscoreName("managed_element")
		if nil == lifecycle.MoCache.managedElement {
			return nil, errors.New("class 'ManagedElement' is not found")
		}

		lifecycle.MoCache.managedObject = definitions.FindByUnderscoreName("managed_object")
		if nil == lifecycle.MoCache.managedObject {
			return nil, errors.New("type 'ManagedObject' isn't found")
		}

		lifecycle.MoCache.networkDevice = definitions.FindByUnderscoreName("network_device")
		if nil == lifecycle.MoCache.networkDevice {
			return nil, errors.New("type 'NetworkDevice' isn't found")
		}
		lifecycle.MoCache.managedLink = definitions.FindByUnderscoreName("network_link")
		if nil == lifecycle.MoCache.managedLink {
			return nil, errors.New("type 'NetworkLink' isn't found")
		}

		lifecycle.timer.Store(time.AfterFunc(60*time.Second, lifecycle.refresh))
	*/
	//lifecycle.RunItInGoroutine(lifecycle.runLoop)
	return lifecycle, nil
}
