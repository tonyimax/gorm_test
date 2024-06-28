package main

import (
	"database/sql"
	"fmt"
	"gorm.io/gorm"
)                              //数据库ORM
import "gorm.io/driver/sqlite" //sqlite数据库驱动
import "gorm.io/driver/mysql"  //mysql数据库驱动

// Product 数据库表字段映射
type Product struct {
	gorm.Model    `json:"gorm_._model"` //数据库记录主键，创建日期，更新日期，删除索引字段
	Code          string                `json:"code,omitempty"`                           //代码
	Price         uint                  `gorm:"not null" json:"price,omitempty"`          //价格
	ProductDetail ProductDetail         `gorm:"embedded" json:"product_detail,omitempty"` //明细
}

// ProductDetail 结构内嵌
type ProductDetail struct {
	Name    sql.NullString `json:"name"`    //名称
	Address sql.NullString `json:"address"` //地址
}

// ProductTidb TiDB使用结构
type ProductTidb struct {
	ID    uint   `gorm:"primaryKey;default:auto_random()"` //表主键
	Code  string `gorm:"type:varchar(100);unique"`         //Code字段
	Price uint   `gorm:"not null"`                         //Price字段
}

func main() {
	//打开sqlite数据库
	db, err := gorm.Open(sqlite.Open("test.db"), //Dialector结构
		&gorm.Config{}) //可选配置
	if err != nil {
		panic("数据库连接失败")
	}
	var product Product                                                  //表字段结构
	db.AutoMigrate(&Product{})                                           //同步结构字段到数据库表，表名为结构名小写的复数形式
	db.Create(&Product{Code: "SN-12345678", Price: 100})                 //向表中添加一条记录
	db.First(&product, "code = ?", "SN-12345678")                        //查询，匹配表products中code字段为SN-12345678
	db.First(&product, 1)                                                //查询，根据主键查询，主键为1的记录
	db.Model(&product).Update("Price", 200)                              //更新，更新主键为1的记录Price字段为200
	db.Model(&product).Updates(Product{Code: "SN-88888888", Price: 500}) //更新，更新主键为1的记录Code字段为SN-88888888，Price字段为500
	db.Delete(&product, 1)                                               //删除主键为1的记录,会向DeletedAt字段插入一个操作日期

	//打开mysql
	dsn := "root:123456@tcp(127.0.0.1:3306)/world?charset=utf8mb4&parseTime=True&loc=Local"
	dbmysql, mysqlerr := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if mysqlerr != nil {
		fmt.Println(mysqlerr)
	}
	dbmysql.AutoMigrate(&Product{})                           //同步表到数据库
	dbmysql.Create(&Product{Code: "SN-12345678", Price: 100}) //添加记录
	dbmysql.Model(&product).Update("Name", "Apple")           //更新字段
	dbmysql.Model(&product).Update("Address", "GZ-CN")        //更新字段

	//打开TiDB
	dsnTidb := "root:123456@tcp(127.0.0.1:4000)/test"
	Tidb, tidbErr := gorm.Open(mysql.Open(dsnTidb), &gorm.Config{})
	if tidbErr != nil {
		fmt.Println(tidbErr)
	}
	Tidb.AutoMigrate(&ProductTidb{})                           //同步表到数据库
	Tidb.Create(&ProductTidb{Code: "SN-12345678", Price: 100}) //添加记录

}
