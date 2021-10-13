package dao

import (
	"encoding/json"
	"fmt"
	model "github.com/rayln/ops/entity"
	"reflect"
	"strings"
)

type BaseDao struct {
}

func (that *BaseDao) Save(userInfo interface{}, base *model.BaseEntity) {
	_, err := base.Save.Insert(userInfo)
	if err != nil {
		panic(err)
	}
}
func (that *BaseDao) UpdateTable(tablename string, where string, data interface{}, base *model.BaseEntity) {
	m := make(map[string]interface{})
	j, _ := json.Marshal(data)
	json.Unmarshal(j, &m)
	verson, ok := m["version"]
	fmt.Println("=====verson:====", verson, ok, m)
	values, index, versionWhere := make([]string, len(m)), 0, "1=1"
	for k, v := range m {
		if k != "version" {
			fmt.Println("===reflect.TypeOf(v).String()===", reflect.TypeOf(v).String())
			if reflect.TypeOf(v).String() == "string" {
				values[index] = fmt.Sprintf("%s='%s'", k, v)
			} else {
				values[index] = fmt.Sprintf("%s=%v", k, v)
			}
		} else {
			values[index], versionWhere = "version=version+1", fmt.Sprintf("version=%v", v)
		}
		index++
	}
	dataSet := strings.Join(values, ",")
	fmt.Println("===values===", dataSet, versionWhere)
	//base.Save.Table(table).Id(id).Update(m)
	base.Save.Exec(fmt.Sprintf("update %s set %s where %s and %s", tablename, dataSet, where, versionWhere))
}
func (that *BaseDao) Update(id interface{}, userInfo interface{}, base *model.BaseEntity) {
	_, err := base.Save.ID(id).AllCols().Update(userInfo)
	if err != nil {
		panic(err)
	}
	//if count == 0 {
	//	panic("Update count is 0!")
	//}
}
func (that *BaseDao) Delete(id interface{}, userInfo interface{}, base *model.BaseEntity) {
	count, err := base.Save.ID(id).Delete(userInfo)
	if err != nil {
		panic(err)
	}
	if count == 0 {
		panic("Delete count is 0!")
	}
	base.Engine.ClearCache(userInfo)
}
func (that *BaseDao) DeleteWhere(where string, userInfo interface{}, base *model.BaseEntity) {
	_, err := base.Save.Where(where).Delete(userInfo)
	if err != nil {
		panic(err)
	}
	base.Engine.ClearCache(userInfo)
}

/**
内部事务提交
*/
//func (that *BaseDao) TransitCommit(transitFunc func(*xorm.Session, *model.BaseEntity), base *model.BaseEntity) {
//	temp := base.Save.Clone()
//	temp.Begin()
//	transitFunc(temp, base)
//	temp.Commit()
//}
