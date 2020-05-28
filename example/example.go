package main

import (
	"fmt"
	"git.corp.kuaishou.com/yuandunsheng/filters"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"net/http"
)

var DB *gorm.DB

type User struct {
	Name  string `filter:"name:name;order;search;match"`
	Age   int    `filter:"name:age;order"`
	Email string `filter:"name:email;search"`
}

func (*User) TableName() string {
	return "user"
}

func init() {
	ptn := "%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=%s"
	sqlUrl := fmt.Sprintf(ptn, "root", "", "localhost", 3306, "test", "Asia%2FShanghai")
	DB, _ = gorm.Open("mysql", sqlUrl)
	DB.LogMode(true)
}

func migrate() {
	d := DB.Set("gorm:table_options", "DEFAULT CHARSET=utf8mb4")
	d.AutoMigrate(&User{})
	d.AutoMigrate(&User{})
	d.AutoMigrate(&User{})
}

func main() {
	//migrate()
	StartAPIServer()
}

func StartAPIServer() {
	r := gin.Default()
	r.GET("/api/list/", ListHandler)
	r.Run(":80")
}

func ListHandler(c *gin.Context) {
	filter := filters.InitModelFilter(c, User{})

	//filter.Limit(10)


	// 计数
	var cnt int
	e := filter.Count(DB, &cnt).Error
	fmt.Println(e, cnt)

	// 手动指定返回字段
	//filter.Select("name,age")

	//手动指定匹配字段
	//filter.Match("name", "tom")

	// 手动指定复杂查询语句
	//filter.Where("name = ? AND age > ?", "tom", 12)
	//var cnt int
	//e := filter.GetQuerySet(DB).Count(&cnt).Error
	//fmt.Println(e)

	var users []User
	filter.GetQuerySet(DB).Find(&users)
	c.JSON(http.StatusOK, users)
}
