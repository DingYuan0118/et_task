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
	Mtime				time.Time	`xorm:"created "`
}

func dbCreate(name string) error {
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

func tableCreate() error {
	db_dst := fmt.Sprintf("root:%s@tcp(127.0.0.1:3306)/%s?charset=utf8", Password, DBname)
	engine, err := xorm.NewEngine("mysql", db_dst)
	if err != nil {
		log.Println(err.Error())
	}
	err = engine.Sync2(new(User))
	if err != nil {
		// 创建数据库
		log.Println(err.Error())
		engine.Close()
		dbCreate("entry_task")
		engine, err = xorm.NewEngine("mysql", db_dst)
		if err != nil {
			log.Println(err.Error())
		}
		
		err = engine.Sync2(new(User))
		// 同步失败
		if err != nil {
			log.Println(err.Error())
			return err
		}
	}
	defer engine.Close()
	return nil
}

func tableInsert() error {
	db_dst := fmt.Sprintf("root:%s@tcp(127.0.0.1:3306)/%s?charset=utf8", Password, DBname)
	engine, err := xorm.NewEngine("mysql", db_dst)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	defer engine.Close()
	err = engine.Sync2(new(User))
	if err != nil {
		log.Println(err.Error())
		return err
	}
	item := User{
		Name: "Ding",
		Nickname: "Ding",
		Password: "dingyuan",
		Profile_pic_url: "User/yuan",
	}
	_, err = engine.Insert(&item)
	return err
}

func DBConnect(dbname, password string) (*xorm.Engine, error)  {
	db_dst := fmt.Sprintf("root:%s@tcp(127.0.0.1:3306)/%s?charset=utf8", password, dbname)
	engine, err := xorm.NewEngine("mysql", db_dst)
	if err != nil {
		return nil, err
	}

	err = engine.Sync(new(User))
	if err != nil {
		engine.Close()
		return nil, err
	}

	return engine, nil
}