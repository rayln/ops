package dao

import (
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/rayln/ops/entity"
)

type TestDao struct {
	BaseDao
}

/**
执行增加一个新记录
*/
func (that *TestDao) Save(baseEntity *entity.BaseEntity) {
	farmTestInfo := entity.FarmTestInfo{Name: "rayln", LastUser: "system", Content: "在结构体中引入tag标签，这样匹配的时候json串对应的字段名需要与tag标签中"}
	//使用baseEntity中的Save来做增删改。（读写分离）
	_, err := baseEntity.Save.Insert(farmTestInfo)
	if err != nil {
		panic(err.Error())
	}
}

/**
查询操作
*/
func (that *TestDao) Query(baseEntity *entity.BaseEntity) *entity.FarmTestInfo {
	var farmTestInfo = new(entity.FarmTestInfo)
	//使用baseEntity.Load对象来做查询（读写分离）
	_, err := baseEntity.Load.ID(268781).Get(farmTestInfo)
	if err != nil {
		panic(err.Error())
	}
	return farmTestInfo
}

/**
redis操作测试
*/
func (that *TestDao) RedisTest(baseEntity *entity.BaseEntity) {
	baseEntity.Redis.Use(func(conn redis.Conn) {
		//存单个数据
		_, err := conn.Do("SET", "mykey", "superWang")
		if err != nil {
			fmt.Println("redis set failed:", err)
		}
		//读取单个数据
		username, err := redis.String(conn.Do("GET", "mykey"))
		if err != nil {
			fmt.Println("redis get failed:", err)
		} else {
			fmt.Printf("Get mykey: %v \n", username)
		}
		//数据是否存在
		is_key_exit, err := redis.Bool(conn.Do("EXISTS", "mykey1"))
		if err != nil {
			fmt.Println("error:", err)
		} else {
			fmt.Printf("exists or not: %v \n", is_key_exit)
		}
		//删除数据
		_, err = conn.Do("DEL", "mykey")
		if err != nil {
			fmt.Println("redis delelte failed:", err)
		}
		//存map数据
		key1 := "profile"
		imap1 := map[string]string{"username": "666", "phonenumber": "888"}
		value1, _ := json.Marshal(imap1)
		n, err := conn.Do("SETNX", key1, value1)
		if err != nil {
			fmt.Println(err)
		}
		if n == int64(1) {
			fmt.Println("success")
		}
		//取map数据
		var imapGet map[string]string
		valueGet, err := redis.Bytes(conn.Do("GET", key1))
		if err != nil {
			fmt.Println(err)
		}
		errShal := json.Unmarshal(valueGet, &imapGet)
		if errShal != nil {
			fmt.Println(err)
		}
		fmt.Println(imapGet["username"])
		fmt.Println(imapGet["phonenumber"])

		//存list
		_, err = conn.Do("lpush", "runoobkey", "redis")
		if err != nil {
			fmt.Println("redis set failed:", err)
		}

		_, err = conn.Do("lpush", "runoobkey", "mongodb")
		if err != nil {
			fmt.Println("redis set failed:", err)
		}
		_, err = conn.Do("lpush", "runoobkey", "mysql")
		if err != nil {
			fmt.Println("redis set failed:", err)
		}
		//取list
		values, _ := redis.Values(conn.Do("lrange", "runoobkey", "0", "100"))
		for _, v := range values {
			fmt.Println(string(v.([]byte)))
		}
		//删除
		conn.Do("DEL", "runoobkey")
	})
}
