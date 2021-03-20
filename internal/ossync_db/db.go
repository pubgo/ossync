package ossync_db

import (
	"github.com/pubgo/golug/db"
	"xorm.io/xorm"
)

func GetDb(names ...string) *xorm.Engine { return db.Get(names...) }
