package util

import (
	"log"
	conf "et-config/src/statusconfig"

	"golang.org/x/crypto/bcrypt"
)

// 第三方包错误统一返回
func ThirdPackageError(err error) (retcode int32, msg string) {
	log.Println(err)
	retcode = int32(conf.StatusThirdPackageErr)
	msg = err.Error()
	return
}

// 数据库密码加密
func HashAndSalt(password string) (string, error) {
	bytepassword := []byte(password)
	hash, err := bcrypt.GenerateFromPassword(bytepassword, bcrypt.MinCost)
	if err != nil {
		return "" , err
	}
	return string(hash), nil
}

// 数据库密码比较
func ComparePasswords(hashedPassword string, password string) (bool, error){
	byteHash := []byte(hashedPassword)
	bytepassword := []byte(password)
	err := bcrypt.CompareHashAndPassword(byteHash, bytepassword)
	if err != nil {
		return false, err
	}
	return true, nil
	
}