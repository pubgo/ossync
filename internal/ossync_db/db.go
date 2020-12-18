package ossync_db

import (
	"github.com/pubgo/golug/golug_consts"
	"github.com/pubgo/golug/golug_db"
	"xorm.io/xorm"
)

var db *xorm.Engine

func GetDb() *xorm.Engine { return db }
func InitDb(name string) {
	if name == "" {
		name = golug_consts.Default
	}
	db = golug_db.GetClient(name)
}
