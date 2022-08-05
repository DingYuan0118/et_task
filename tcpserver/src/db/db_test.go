package db

import "testing"

func TestTableInsert(t *testing.T) {
	tableCreate()
	// hashed password for "dingyuan"
	err := tableInsert("Ding", "$2a$04$EMdYwzi3AQH9LpVbI8wg2O9IUfSute3aJVEygGRkyWEN/FXuscz/u")
	if err != nil {
		t.Error(err.Error())
	}
}