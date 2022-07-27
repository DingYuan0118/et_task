package db

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"xorm.io/xorm"
)

var engine *xorm.Engine

type User struct {
	Id					int64
	Name 				string		`xorm:"varchar(64) not null unique 'usr_name'"`
	Nickname 			string 		`xorm:"varchar(64) 'usr_nickname'"`
	Password			string		`xorm:"varchar(32) 'usr_password'"`
	Profile_pic_url 	string		`xorm:"varchar(1024) "`
	Ctime				time.Time	`xorm:"created "`
	Mtime				time.Time	`xorm:"datetime "`
}

const Password string = "dingyuan841218"
const DB_name string = "entry_task"

func db_create(name string) error {
	var db_dst string
	db_dst = fmt.Sprintf("root:%s@tcp(127.0.0.1:3306)/", Password)

	db, err := sql.Open("mysql", db_dst)
	if err != nil {
		return err
	}
	defer db.Close()
 
	_,err = db.Exec("CREATE DATABASE "+ name)
	if err != nil {
		return err
	}
	return nil
}

func table_create() error {
	db_dst := fmt.Sprintf("root:%s@tcp(127.0.0.1:3306)/%s?charset=utf8", Password, DB_name)
	engine, err := xorm.NewEngine("mysql", db_dst)
	if err != nil {
		log.Println(err.Error())
	}

	err = engine.Sync2(new(User))
	if err != nil {
		log.Println(err.Error())
		engine.Close()
		db_create("entry_task")
		engine, err = xorm.NewEngine("mysql", db_dst)
		err = engine.Sync2(new(User))
		if err != nil {
			log.Println(err.Error())
		}
	}
	defer engine.Close()
	return nil
}