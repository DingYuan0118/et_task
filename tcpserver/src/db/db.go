package db

import (
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"xorm.io/xorm"
	// "tcpserver/src/serverfunc"
)

type DBclient struct {
	Engine *xorm.Engine
}

var Client *DBclient

func NewDBclient(dst string, DBMaxOpenConns, DBMaxIdleConns int, DBConnMaxLifetime time.Duration) (error) {
	engine, err := xorm.NewEngine("mysql", dst)
	engine.SetMaxIdleConns(DBMaxIdleConns)
	engine.SetMaxOpenConns(DBMaxOpenConns)
	engine.SetConnMaxLifetime(DBConnMaxLifetime)
	if err != nil {
		return err
	}
	Client = &DBclient{Engine: engine}
	return nil
}

func init(){
	db_dst := fmt.Sprintf("root:%s@tcp(127.0.0.1:3306)/%s?charset=utf8", Password, DBname)
	err := NewDBclient(db_dst, DBMaxOpenConns, DBMaxIdleConns, DBConnMaxLifetime)
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