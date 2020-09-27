package entity

type FarmTestInfo struct {
	Id       int64  `xorm:"pk autoincr comment('主键ID')"`
	Name     string `xorm:"not null varchar(24) comment('名字')" json:"name"`
	Content  string `xorm:"not null text comment('内容')" json:"content"`
	LastUser string `xorm:"varchar(10) default 'system'  comment('最后更新人')"`
	LastTime string `xorm:"timestamp updated comment('最后更新时间')"`
	Version  int    `xorm:"version"`
}
