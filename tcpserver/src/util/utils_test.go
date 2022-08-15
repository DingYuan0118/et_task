package util

import (
	"fmt"
	"log"
	"testing"
)

func TestHashAndSalt(t *testing.T) {
	hashed_pwd, err := HashAndSalt("dingyuan")
	if err != nil {
		t.Error(err)
	}
	log.Println("dingyuan hashed password: ", hashed_pwd)
}

func TestComparePasswords(t *testing.T) {
	var tests = []struct{
		password string
		hashedpwd string
		want bool
	}{
		{"dingyuan", "$2a$04$EMdYwzi3AQH9LpVbI8wg2O9IUfSute3aJVEygGRkyWEN/FXuscz/u", true},
		{"dingyuan1", "$2a$04$EMdYwzi3AQH9LpVbI8wg2O9IUfSute3aJVEygGRkyWEN/FXuscz/u", false},
	}

	for _, test := range tests {
		got, _ := ComparePasswords(test.hashedpwd, test.password)
		if got != test.want {
			t.Errorf("%s and %s result is %v, want %v", test.password, test.hashedpwd, got, test.want)
		}
	}
}

func TestMD5Encode(t *testing.T) {
	password := "dingyuan"
	md5password := MD5Encode(password)
	fmt.Println(md5password)
}