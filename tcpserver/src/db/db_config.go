package db

import "time"

const Password string = "dingyuan841218"
const DBname string = "entry_task"

const DBMaxOpenConns int = 50
const DBMaxIdleConns int = 50
const DBConnMaxLifetime = time.Second * 5