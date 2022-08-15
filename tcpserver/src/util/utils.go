package util

import (
	conf "et-config/src/statusconfig"
	"fmt"

	"golang.org/x/crypto/bcrypt"
	"crypto/md5"
)

// 第三方包错误统一返回
func ThirdPackageError(err error) (int32, string) {
	retcode := int32(conf.StatusServerError)
	msg := err.Error()
	return retcode, msg
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

// MD5 加密
func MD5Encode(str string) string {
	data := []byte(str) //切片
	has := md5.Sum(data)
	md5str := fmt.Sprintf("%x", has) //将[]byte转成16进制
	return md5str
}

// MD5 哈希密码解密
func MD5ComparePasswords(MD5Password string, password string) (bool, error) {
	md5_password := MD5Encode(password)
	if MD5Password == md5_password {
		return true, nil
	}
	return false, nil
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