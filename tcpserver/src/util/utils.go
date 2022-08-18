package util

import (
	conf "et-config/src/statusconfig"
	"fmt"

	"golang.org/x/crypto/bcrypt"
	"crypto/md5"
)

// ThirdPackageError return the error code and msg casused by third-party package
func ThirdPackageError(err error) (int32, string) {
	retcode := int32(conf.StatusServerError)
	msg := err.Error()
	return retcode, msg
}

// HashAndSalt return the encrypted user password by bcrypt
func HashAndSalt(password string) (string, error) {
	bytepassword := []byte(password)
	hash, err := bcrypt.GenerateFromPassword(bytepassword, bcrypt.MinCost)
	if err != nil {
		return "" , err
	}
	return string(hash), nil
}

// MD5Encode return the encoded password by MD5 method
func MD5Encode(str string) string {
	data := []byte(str) //切片
	has := md5.Sum(data)
	md5str := fmt.Sprintf("%x", has) //将[]byte转成16进制
	return md5str
}

// MD5ComparePasswords return the compare result of md5Password and user password
func MD5ComparePasswords(md5Password string, password string) (bool, error) {
	md5_password := MD5Encode(password)
	if md5Password == md5_password {
		return true, nil
	}
	return false, nil
} 

// ComparaPassword return the compare result of bcypted password and user password
func ComparePasswords(hashedPassword string, password string) (bool, error){
	byteHash := []byte(hashedPassword)
	bytepassword := []byte(password)
	err := bcrypt.CompareHashAndPassword(byteHash, bytepassword)
	if err != nil {
		return false, err
	}
	return true, nil
	
}