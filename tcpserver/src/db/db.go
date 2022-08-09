package db

import (
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"xorm.io/xorm"
	// "tcpserver/src/serverfunc"
)

type User struct {
	Id					int64
	Name 				string		`xorm:"varchar(64) not null unique 'usr_name'"`
	Nickname 			string 		`xorm:"varchar(64) 'usr_nickname'"`
	Password			string		`xorm:"varchar(64) 'usr_password'"`
	Profile_pic_url 	string		`xorm:"varchar(1024) "`
	Ctime				time.Time	`xorm:"created"`
	Mtime				time.Time	`xorm:"updated"`
}

type DBclient struct {
	Engine *xorm.Engine
}

var Client *DBclient

func NewDBclient(dst string, dbMaxOpenConns int) (error) {
	engine, err := xorm.NewEngine("mysql", dst)
	engine.SetMaxOpenConns(dbMaxOpenConns)
	if err != nil {
		return err
	}
	Client = &DBclient{Engine: engine}
	return nil
}

func init(){
	db_dst := fmt.Sprintf("root:%s@tcp(127.0.0.1:3306)/%s?charset=utf8", Password, DBname)
	err := NewDBclient(db_dst, DBMaxOpenConns)
	if err != nil {
		log.Println(err)
	}
}

func DBConnect() (*xorm.Engine, error)  {
	err := Client.Engine.Sync(new(User))
	if err != nil {
		return nil, err
	}

	return Client.Engine, nil
}