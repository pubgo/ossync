package models

import "os"

type SyncFile struct {
	Id        int64       `json:"id"`
	CreatedAt JsonTime    `json:"created" xorm:"created"`
	UpdatedAt JsonTime    `json:"updated" xorm:"updated"`
	Crc64ecma uint64      `json:"crc64_ecma"`
	Md5       string      `json:"md5"`
	Name      string      `json:"name"`
	Path      string      `json:"path" xorm:"varchar(100) notnull unique 'path'"`
	Changed   bool        `json:"changed"`
	Synced    bool        `json:"synced"`
	Size      int64       `json:"size"`
	Mode      os.FileMode `json:"mode"`
	ModTime   int64       `json:"mod_time"`
	IsDir     bool        `json:"-"`
}
