package ossync_db

import (
	"github.com/pubgo/golug/golug_db"
	"xorm.io/xorm"
)

var db *xorm.Engine

func GetDb() *xorm.Engine { return db }
func InitDb(name string)  { db = golug_db.GetClient(name) }
