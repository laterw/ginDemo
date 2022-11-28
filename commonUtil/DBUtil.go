package commonUtil

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"net/url"
)

var DB *gorm.DB

// InitDB  初始化 数据信息
func InitDB() *gorm.DB {
	host := "127.0.0.1"    //viper.GetString("datasource.host")
	port := "3306"         //viper.GetString("datasource.port")
	database := "test"     //viper.GetString("datasource.database")
	username := "root"     //viper.GetString("datasource.username")
	password := "root"     // viper.GetString("datasource.password")
	charset := "utf8"      //viper.GetString("datasource.charset")
	loc := "Asia/Shanghai" //viper.GetString("datasource.loc")
	args := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=true&loc=%s",
		username,
		password,
		host,
		port,
		database,
		charset,
		url.QueryEscape(loc))
	fmt.Printf(args)
	db, err := gorm.Open(mysql.Open(args), &gorm.Config{})
	if err != nil {
		panic("fail to connect database, err: " + err.Error())
	}

	DB = db
	return db
}
func GetDB() *gorm.DB {
	return DB
}
