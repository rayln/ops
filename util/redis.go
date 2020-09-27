package util

import (
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"time"
)

type Redis struct {
	Pool *redis.Pool
}

/**
初始化redis
*/
func (that *Redis) Init(ip string, pwd string) *Redis {
	var redisPool = &redis.Pool{
		MaxIdle:     100,
		MaxActive:   1000,
		IdleTimeout: 30 * time.Second,
		Dial: func() (conn redis.Conn, err error) {
			return redis.Dial("tcp", fmt.Sprintf("%s:6379", ip), redis.DialPassword(fmt.Sprintf("%s", pwd)), redis.DialDatabase(0))
		},
	}
	that.Pool = redisPool
	return that
}

/**
使用Redis
*/
func (that *Redis) Use(action func(redis.Conn)) {
	conn := that.Pool.Get()
	defer conn.Close()
	action(conn)
}

/**
设置单个key，value的值
*/
func (that *Redis) SetValue(key string, value interface{}, conn redis.Conn) {
	_, err := conn.Do("SET", key, value)
	if err != nil {
		panic(err)
	}
}

/**
获取单个value值，根据key
*/
func (that *Redis) GetValueString(key string, conn redis.Conn) string {
	result, err := redis.String(conn.Do("GET", key))
	if err != nil {
		panic(err)
	}
	return result
}

/**
获取单个value值，根据key
*/
func (that *Redis) GetValueInt(key string, value interface{}, conn redis.Conn) int {
	result, err := redis.Int(conn.Do("GET", key))
	if err != nil {
		panic(err)
	}
	return result
}

/**
获取单个value值，根据key
*/
func (that *Redis) GetValueInterface(key string, value interface{}, conn redis.Conn) interface{} {
	result, err := conn.Do("GET", key)
	if err != nil {
		panic(err)
	}
	return result
}

/**
Set Map方法，会先删除Redis中Key值的值。再进行插入操作。
*/
func (that *Redis) SetMap(key string, value interface{}, conn redis.Conn) {
	valueStr, _ := json.Marshal(value)
	conn.Do("DEL", key)
	n, err := conn.Do("SETNX", key, valueStr)
	if err != nil {
		panic(err)
	}
	if n == int64(1) {
		fmt.Println("map save success!!")
	} else {
		panic("map save fail!!!")
	}
}

func (that *Redis) GetMap(key string, result interface{}, conn redis.Conn) {
	valueGet, err := redis.Bytes(conn.Do("GET", key))
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(valueGet, &result)
	if err != nil {
		panic(err)
	}
}

func (that *Redis) Delete(key string, conn redis.Conn) {
	conn.Do("DEL", key)
}

/**
放入List
注意一点，List是用插入最前的方式，也就是最新的数放在最前面
*/
func (that *Redis) SetList(key string, value string, conn redis.Conn) {
	_, err := conn.Do("lpush", key, value)
	if err != nil {
		panic(err)
	}
}

/**
取得List
start=0 开始
end=100 结束第100
*/

func (that *Redis) GetList(key string, start string, end string, conn redis.Conn) []interface{} {
	values, _ := redis.Values(conn.Do("lrange", key, start, end))
	/**
	//value的取值方法
	for _, v := range testList {
		fmt.Println(string(v.([]byte)))
	}
	*/
	return values
}

/**
测试用例
*/
func (that *Redis) Test() {
	conn := that.Pool.Get()
	defer conn.Close()
	_, err := conn.Do("SET", "mykey", "superWang")
	if err != nil {
		panic(err.Error())
	}
	username, err := redis.String(conn.Do("GET", "mykey"))
	if err != nil {
		panic(err.Error())
	} else {
		fmt.Printf("Get mykey: %v \n", username)
	}

}
