package models

import (
	"github.com/eatmoreapple/openwechat"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"time"
)

// 城市编码
type Encode struct {
	Name     string `gorm:"column:name"`      // 城市名字
	AdCode   string `gorm:"column:ad_code"`   // 城市编码
	CityCode string `gorm:"column:city_code"` // 区号
}

var DB *gorm.DB

type MessageTime struct {
	gorm.Model
	Msg          string
	MsgId        string
	FromUserName string
}

// 初始化 MYSQL
func MYSQL() {

	err := InitializeMYSQL()
	if err != nil {
		log.Fatalln("MYSQL 初始化失败！" + err.Error())
	}

	log.Println("MYSQL 初始化成功！")
}

// 定义一个方法程序开始时，
// 初始化数据并且连接数据库拿到数据存放在redis中
func InitializeMYSQL() (err error) {
	dns := viper.GetString("mysql.dns")

	// 创建 ORM 连接池
	DB, err = gorm.Open(mysql.Open(dns), &gorm.Config{
		// 设置输出条件
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		return err
	}

	// 拿到底层的 MySql 连接池 DB
	sqldb, err := DB.DB()
	if err != nil {
		return err
	}

	// 测试是否连接上数据
	err = sqldb.Ping()
	if err != nil {
		return err
	}

	// 连接上了设置连接池的设置(默认的设置就可以应用于平常的需求)
	// 设置连接池的数量时间等
	sqldb.SetMaxOpenConns(10)
	sqldb.SetMaxIdleConns(10)
	sqldb.SetConnMaxIdleTime(time.Hour)
	sqldb.SetConnMaxLifetime(time.Hour)

	//DB.AutoMigrate(&MessageTime{})

	return
}

// 将获取的城市查询对应的编码
func FincCode(addrs string) (*Encode, string) {
	usercode := &Encode{}
	err := DB.Model(&Encode{}).Where("name like ?", "%"+addrs+"%").First(usercode).Error
	if err != nil {
		log.Println(addrs + "城市编码查询错误！ err: " + err.Error())
		return nil, addrs
	}

	return usercode, ""
}

// 记录每一个人发送的消息
func StorageMsg(msg *openwechat.Message) {
	content := &MessageTime{
		Msg:          msg.Content,
		FromUserName: msg.FromUserName,
		MsgId:        msg.MsgId,
	}

	err := DB.Create(content).Error
	if err != nil {
		log.Println(("消息存储错误！" + err.Error()))
	}
}

//func StorageMsg(msg *openwechat.Message) {
//	content := &MessageTime{
//		Msg:          msg.Content,
//		FromUserName: msg.FromUserName,
//		MsgId:        msg.MsgId,
//	}
//
//	err := DB.Create(content).Error
//	if err != nil {
//		log.Println("消息存储错误！")
//	}
//
//}
