package db

import (
	"log"
	"os"
	"time"

	"github.com/gofish2020/easysearch/conf"
	"github.com/gofish2020/easysearch/tools"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const SpiderDBName = "spider"
const DictDBName = "dictory"

const baseUrl = "https://www.baidu.com"

var mysqlDB *gorm.DB

func InitDB() {

	// gorm SQL 日志
	file, err := os.Create("gorm-log.txt")
	if err != nil {
		panic(err)
	}

	logLevel := logger.Warn
	if conf.DBConfig.LogLevel == "debug" {
		logLevel = logger.Info
	}

	fileLogger := logger.New(
		log.New(file, "", log.LstdFlags), // io writer（日志输出的目标，前缀和日志包含的内容——译者注）
		logger.Config{
			SlowThreshold:             time.Second * 6, // 慢 SQL 阈值
			LogLevel:                  logLevel,        // 日志级别
			IgnoreRecordNotFoundError: true,            // 忽略ErrRecordNotFound（记录未找到）错误
			Colorful:                  false,           // 禁用彩色打印
		},
	)

	gormConfig := gorm.Config{
		Logger: fileLogger,
	}

	// 2.建立db连接
	_db0, _ := gorm.Open(mysql.Open(conf.DBConfig.ConnectString()), &gormConfig)
	dbdb0, _ := _db0.DB()
	dbdb0.SetMaxIdleConns(1)
	dbdb0.SetMaxOpenConns(20)
	dbdb0.SetConnMaxLifetime(time.Hour)

	mysqlDB = _db0

	//3.创建库 spider
	mysqlDB.Exec("CREATE DATABASE IF NOT EXISTS " + SpiderDBName + " CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci ;")
	//3.创建库 dictory
	mysqlDB.Exec("CREATE DATABASE IF NOT EXISTS " + DictDBName + " CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci ;")

	//4.创建表 t_dict
	mysqlDB.Exec("use " + DictDBName + ";")
	mysqlDB.Exec("CREATE TABLE if not exists `" + dictTable.TableName() + "` (   `id` bigint unsigned NOT NULL AUTO_INCREMENT comment '自增id', `word` varchar(255) DEFAULT '' comment '分词', `positions` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_bin comment '文档id,文档长度,词频,分词偏移-文档id,文档长度,词频,分词偏移',`deleted_at` timestamp  DEFAULT NULL,`created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '插入时间',`updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '插入时间', PRIMARY KEY (`id`), UNIQUE KEY `word` (`word`) )    ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;")
	//5.创建表 t_source
	mysqlDB.Exec("use " + SpiderDBName + ";")
	mysqlDB.Exec("CREATE TABLE if not exists `" + srcTable.TableName() + "` ( `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '自增id',`title` varchar(1000) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '标题',`url` varchar(1024) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '链接',`md5` varchar(32) CHARACTER SET utf8 COLLATE utf8_unicode_ci NOT NULL DEFAULT '' COMMENT 'url对应的md5值',`html` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci COMMENT '页面文字',`host` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '域名',`craw_done` tinyint NOT NULL DEFAULT '0' COMMENT '0 未爬 1 已爬 2 失败',`dict_done` tinyint NOT NULL DEFAULT '0' COMMENT '0 未分词 1 已分词',`deleted_at` timestamp  DEFAULT NULL,`created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '插入时间',`updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '插入时间',PRIMARY KEY (`id`),UNIQUE KEY `md5` (`md5`)) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;")
	mysqlDB.Exec("INSERT IGNORE INTO `"+srcTable.TableName()+"` (url,md5) values (?,?) ", baseUrl, tools.GetMD5Hash(baseUrl))

	//6.创建黑名单-分词
	mysqlDB.Exec("CREATE TABLE if not exists `" + blackList.TableName() + "` (   `id` int unsigned NOT NULL AUTO_INCREMENT,   `word` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,`deleted_at` timestamp  DEFAULT NULL,`created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '插入时间',`updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',   PRIMARY KEY (`id`) ) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;")
	mysqlDB.Exec("INSERT IGNORE INTO `" + blackList.TableName() + "` (`id`, `word`) VALUES   (1, 'px'),   (2, '20'),   (3, '('),   (4, ')'),   (5, ','),   (6, '.'),   (7, '-'),   (8, '/'),   (9, ':'),   (10, 'var'),   (11, '的'),   (12, 'com'),   (13, ';'),   (14, '['),   (15, ']'),   (16, '{'),   (17, '}'),   (18, \"'\"),   (19, '\"'),   (20, '_'),   (21, '?'),   (22, 'function'),   (23, 'document'),   (24, '|'),   (25, '='),   (26, 'html'),   (27, '内容'),   (28, '0'),   (29, '1'),   (30, '3'),   (31, 'https'),   (32, 'http'),   (33, '2'),   (34, '!'),   (35, 'window'),   (36, 'if'),   (37, '“'),   (38, '”'),   (39, '。'),   (40, 'src'),   (41, '中'),   (42, '了'),   (43, '6'),   (44, '｡'),   (45, '<'),   (46, '>'),   (47, '联系'),   (48, '号'),   (49, 'getElementsByTagName'),   (50, '5'),   (51, '､'),   (52, 'script'),   (53, 'js');")
}

var srcTable = Source{}

// 库spider 表 t_source
type Source struct {
	gorm.Model

	Title    string `gorm:"column:title"`
	Url      string `gorm:"column:url"`
	Md5      string `gorm:"column:md5"`
	Html     string `gorm:"column:html"`
	Host     string `gorm:"column:host"`
	CrawDone string `gorm:"column:craw_done"`
	DictDone string `gorm:"column:dict_done"`
}

func (s Source) TableName() string {
	return "t_source"
}

var dictTable = Dict{}

// 库 dictory 表 t_dict
type Dict struct {
	gorm.Model

	Word      string `gorm:"column:word"`
	Positions string `gorm:"column:positions"`
}

func (d Dict) TableName() string {
	return "t_dict"
}

var blackList = BlackList{}

type BlackList struct {
	gorm.Model

	Word string `gorm:"column:word"`
}

func (b BlackList) TableName() string {
	return "t_blacklist"
}
