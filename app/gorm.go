package app

import (
	"errors"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	. "k3gin/app/config"
	"log"
	"os"
	"strings"
	"time"
)

const (
	R = iota
	W
)

type DB struct {
	RDB *gorm.DB
	WDB *gorm.DB
}

func InitGormDB() (*DB, func(), error) {

	RDB, err := NewGormDB(R)
	if err != nil {
		return nil, nil, err
	}

	WDB, err := NewGormDB(W)
	if err != nil {
		return nil, nil, err
	}

	db := DB{
		RDB: RDB,
		WDB: WDB,
	}

	// 以前Gorm需要关注close链接，现在gorm不需要再关注close db的操作了
	clean := func() {

	}

	// TODO 自动映射数据库表结构, 后面有时间再补充
	if C.Gorm.EnableAutoMigrate {

	}

	return &db, clean, nil
}

func NewGormDB(RW int) (*gorm.DB, error) {
	var dialector gorm.Dialector

	switch strings.ToLower(C.Gorm.DBType) {
	case "mysql":
		if RW == R {
			dialector = mysql.Open(C.RMySQL.DSN())
		} else {
			dialector = mysql.Open(C.WMySQL.DSN())
		}
	case "postgers":
		// TODO 后续可以补充
	case "sqlit":
		// TODO 后续可以补充
	default:
		return nil, errors.New("unknown db")
	}

	db, err := gorm.Open(dialector, &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   C.Gorm.TablePrefix,
			SingularTable: true,
		},
		Logger: logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			Colorful:                  false,       // Disable color
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			LogLevel:                  logger.Info, // sql日志等级
		}),
	})

	if err != nil {
		return nil, err
	}

	if C.Gorm.Debug {
		db = db.Debug()
	}

	sqlDB, err := db.DB()

	if err != nil {
		return nil, err
	}

	// 设置连接池
	sqlDB.SetMaxIdleConns(C.Gorm.MaxIdleConns)
	sqlDB.SetMaxOpenConns(C.Gorm.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(C.Gorm.MaxLifetime) * time.Second)

	return db, nil
}
