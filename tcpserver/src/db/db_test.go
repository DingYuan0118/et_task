package db

import "testing"

func TestTableInsert(t *testing.T) {
	tableCreate()
	err := tableInsert()
	if err != nil {
		t.Error(err.Error())
	}
}