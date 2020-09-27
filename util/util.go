package util

import (
	"crypto/md5"
	"encoding/hex"
	"golang.org/x/crypto/bcrypt"
	"reflect"
)

type Util struct {
}

/**
Struct转Map方法（通过反射实现）
*/
func StructToMap(obj interface{}) map[string]interface{} {
	obj1 := reflect.TypeOf(obj)
	obj2 := reflect.ValueOf(obj)

	var data = make(map[string]interface{})
	for i := 0; i < obj1.NumField(); i++ {
		data[obj1.Field(i).Name] = obj2.Field(i).Interface()
	}

	return data
}

/**
对密码进行HASH加密
*/
func GeneratePassword(pwd string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
}

/**
校验密码是否正确
userPassword 用户输入的密码
dbUserPassword 数据库的用户密码
*/
func IsValidatePassword(userPassword string, dbUserPassword []byte) (bool, error) {
	if err := bcrypt.CompareHashAndPassword(dbUserPassword, []byte(userPassword)); err != nil {
		return false, err
	}
	return true, nil
}

/**
md5加密
*/
func MD5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}
