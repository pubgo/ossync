package models

import (
	"time"

	"github.com/pubgo/xerror"
	"xorm.io/xorm"
)

type JsonTime time.Time

func (j JsonTime) MarshalJSON() ([]byte, error) {
	return []byte(`"` + time.Time(j).Format("2006-01-02 15:04:05") + `"`), nil
}

func NextPage(page, perPage, total int64) (int64, int64) {
	if total > perPage {
		return page + 1, total/perPage + 1
	}
	return page, total/perPage + 1
}

func Pagination(page, perPage int) (int, int) {
	if perPage < 1 {
		perPage = 20
	}

	if perPage > 100 {
		perPage = 20
	}

	if page < 2 {
		page = 1
	}

	return page, perPage
}

func pagination(page, perPage int) (int, int, int) {
	if perPage < 1 {
		perPage = 20
	}

	if perPage > 100 {
		perPage = 20
	}

	if page < 2 {
		page = 1
	}

	return page, perPage, (page - 1) * perPage
}

func Random(db *xorm.Engine, n int, table string, data interface{}) (err error) {
	defer xerror.RespErr(&err)
	return xerror.Wrap(db.SQL(
		"select * from ? where id>=(select floor(rand() * (select max(id) from ?))) order by id limit ?", n, table, table,
	).Find(data))
}
