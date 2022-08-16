package db

import (
	"time"
)

type User struct {
	Id              uint64    `xorm:"not null pk autoincr UNSIGNED BIGINT(20)"`
	Name            string    `xorm:"not null 'usr_name' default '' unique(usr_name) VARCHAR(64)"`
	Nickname        string    `xorm:"not null 'usr_nickname' default '' VARCHAR(64)"`
	Password        string    `xorm:"not null 'usr_password' default '' VARCHAR(64)"`
	Profile_pic_url string    `xorm:"not null 'profile_pic_url' default '' VARCHAR(1024)"`
	Ctime           time.Time `xorm:"not null default CURRENT_TIMESTAMP comment('create time') DATETIME"`
	Mtime           time.Time `xorm:"not null default CURRENT_TIMESTAMP comment('modified time') DATETIME"`
}
