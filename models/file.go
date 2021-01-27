package models

import (
	"os"

	"github.com/pubgo/ossync/internal/ossync_db"
	"github.com/pubgo/xerror"
	"xorm.io/xorm"
)

type SyncFile struct {
	Id        int64       `json:"id"`
	CreatedAt JsonTime    `json:"created" xorm:"created"`
	UpdatedAt JsonTime    `json:"updated" xorm:"updated"`
	Crc64ecma string      `json:"crc64_ecma" xorm:"varchar(100) notnull index 'crc64_ecma'"`
	Md5       string      `json:"md5"`
	Name      string      `json:"name"`
	Path      string      `json:"path" xorm:"varchar(100) notnull index 'path'"`
	PathHash  string      `json:"path_hash" xorm:"varchar(100) notnull unique 'path_hash'"`
	Changed   bool        `json:"changed"`
	Synced    bool        `json:"synced"`
	Size      int64       `json:"size"`
	Mode      os.FileMode `json:"mode"`
	ModTime   int64       `json:"mod_time"`
	IsDir     bool        `json:"-"`
}

func SyncFileCreate(sf *SyncFile) {
	_, err := ossync_db.GetDb().InsertOne(sf)
	xerror.Panic(err)
}

func SyncFileFindOne(where string, a ...interface{}) *SyncFile {
	tb := ossync_db.GetDb().Table(&SyncFile{})
	var sf = new(SyncFile)
	_, err := tb.Where(where, a...).Get(sf)
	if err == xorm.ErrNotExist {
		return nil
	}

	xerror.Panic(err)
	return sf
}

func SyncFileUpdate(sf *SyncFile, where string, a ...interface{}) {
	tb := ossync_db.GetDb().Table(&SyncFile{})
	_, err := tb.Where(where, a...).Update(sf)
	if err == xorm.ErrNotExist {
		return
	}

	xerror.Panic(err)
	return
}

func SyncFileUpdateMap(sf map[string]interface{}, where string, a ...interface{}) {
	tb := ossync_db.GetDb().Table(&SyncFile{})
	_, err := tb.Where(where, a...).Update(sf)
	if err == xorm.ErrNotExist {
		return
	}

	xerror.Panic(err)
	return
}

func SyncFileDelete(where string, a ...interface{}) {
	tb := ossync_db.GetDb().Table(&SyncFile{})
	_, err := tb.Where(where, a...).Delete(&SyncFile{})
	if err == xorm.ErrNotExist {
		return
	}
	xerror.Panic(err)
	return
}

func SyncFileRange(page, perPage int, where string, a ...interface{}) ([]SyncFile, int64, error) {
	tb := ossync_db.GetDb().Table(&SyncFile{})

	var sfList []SyncFile

	count, err := Range(tb, &sfList, page, perPage, where, a...)
	return sfList, count, xerror.Wrap(err)
}

func SyncFileEach(fn func(sf SyncFile)) {
	tb := ossync_db.GetDb().Table(&SyncFile{})

	id := int64(0)
	for i := 1; ; i++ {
		var sfList []SyncFile
		_, perPage, start := pagination(i, 20)
		xerror.Panic(tb.Where("id>=?", id).Limit(perPage, start).Find(&sfList))
		if len(sfList) == 0 {
			break
		}

		for i := range sfList {
			sf := sfList[i]
			id = sf.Id
			fn(sf)
		}
	}
}
