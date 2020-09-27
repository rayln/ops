package intercepter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/kataras/iris/v12/context"
	"github.com/rayln/ops/util"
	"sort"
	"strings"
)

type Intercepter struct {
	Key string
}

func (that *Intercepter) Init() *Intercepter {
	that.Key = "&key=ABC"
	return that
}

func (that *Intercepter) Handle(request context.Context) bool {
	var jsonMap map[string]interface{}
	err := request.ReadJSON(&jsonMap)
	if err != nil {
		request.Application().Logger().Error("Intercepter.Handle Error: " + err.Error())
		return false
	}
	dataJson, _ := json.Marshal(jsonMap)
	d := json.NewDecoder(bytes.NewReader(dataJson))
	d.UseNumber()
	d.Decode(&jsonMap)
	fmt.Println("====jsonMap====", jsonMap)
	//d := jsonMap.NewDecoder(bytes.NewReader(dataJson))
	//获取key组成array
	var arrayKey = make([]string, 0, len(jsonMap))
	for key, _ := range jsonMap {
		arrayKey = append(arrayKey, key)
	}
	//排序
	sort.Strings(arrayKey)
	//拼装
	var strList = make([]string, 0, len(arrayKey))
	var webSign string
	for _, key := range arrayKey {
		//拼装移除sign部分
		if key == "sign" {
			webSign = fmt.Sprintf("%v", jsonMap[key])
			continue
		}
		var temp interface{}
		switch jsonMap[key].(type) {
		case float64:
			temp = jsonMap[key].(float64)
		case string:
			temp = jsonMap[key].(string)
		default:
			temp = jsonMap[key]
		}
		value := fmt.Sprintf("%v", temp)
		strList = append(strList, fmt.Sprintf("%s=%s", key, value))
	}
	var str = fmt.Sprintf("%s%s", strings.Join(strList, "&"), that.Key)
	var sign = util.MD5(str)
	request.Application().Logger().Info("path:", request.Path(), "  requestStr:", str, "  webSign:", webSign, "  sign:", sign)
	return webSign == sign
}
