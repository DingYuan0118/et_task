package db

import (
	"fmt"
	"testing"
)

func TestDBupdate(t *testing.T) {
	engine, err := DBConnect()
	if err != nil {
		t.Fatal(err)
	}
	user := new(User)
	user.Nickname = "Ding1234"

	session := engine.NewSession()
	session.Begin()
	defer session.Close()

	res, err := session.Where("usr_name = ?", "Ding").Cols("usr_nickname").Update(user)
	session.Commit()
	if err != nil {
		t.Fatal(err)
	}
	if res == 0 {
		t.Fatalf("res: %d", res)
	}else{
		fmt.Printf("res: %d", res)
	}
}