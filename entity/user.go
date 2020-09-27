package entity

type User struct {
	//TODO
	Id       int64  `xorm:"pk autoincr"`
	Name     string `xorm:"varchar(24)"`
	LastUser string `xorm:"varchar(40)"`
	LastTime string `xorm:"time.Time updated"`
}
