package mysql

import (
	"fmt"
	"time"

	"gorm.io/gorm/logger"

	"github.com/IfanTsai/go-lib/config"
	"github.com/pkg/errors"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type DB struct {
	Reader *gorm.DB
	Writer *gorm.DB
}

type DBConfig struct {
	User                string
	Password            string
	Addr                string
	DBName              string
	MaxOpenConnections  int
	MaxIdleConnections  int
	MaxLifetimeDuration time.Duration
	LogFilename         string
	Log                 bool
}

func NewDB(readerConfig, writerConfig *DBConfig) (*DB, error) {
	dbReader, err := dbConnect(readerConfig)
	if err != nil {
		return nil, err
	}

	dbWriter, err := dbConnect(writerConfig)
	if err != nil {
		return nil, err
	}

	return &DB{
		Reader: dbReader,
		Writer: dbWriter,
	}, nil
}

func NewDBInConfig(logFilename string) (*DB, error) {
	readConfig := &DBConfig{
		User:                config.GetConfig().MySQL.Read.User,
		Password:            config.GetConfig().MySQL.Read.Password,
		Addr:                config.GetConfig().MySQL.Read.Addr,
		DBName:              config.GetConfig().MySQL.Read.Name,
		MaxOpenConnections:  config.GetConfig().MySQL.MaxOpenConnections,
		MaxIdleConnections:  config.GetConfig().MySQL.MaxIdleConnections,
		MaxLifetimeDuration: config.GetConfig().MySQL.MaxLifetimeDuration,
		Log:                 config.GetConfig().MySQL.Log,
		LogFilename:         logFilename,
	}

	writeConfig := &DBConfig{
		User:                config.GetConfig().MySQL.Write.User,
		Password:            config.GetConfig().MySQL.Write.Password,
		Addr:                config.GetConfig().MySQL.Write.Addr,
		DBName:              config.GetConfig().MySQL.Write.Name,
		MaxOpenConnections:  config.GetConfig().MySQL.MaxOpenConnections,
		MaxIdleConnections:  config.GetConfig().MySQL.MaxIdleConnections,
		MaxLifetimeDuration: config.GetConfig().MySQL.MaxLifetimeDuration,
		Log:                 config.GetConfig().MySQL.Log,
		LogFilename:         logFilename,
	}

	return NewDB(readConfig, writeConfig)
}

func (db *DB) Close() error {
	err := db.CloseReader()
	if err != nil {
		return err
	}

	return db.CloseWriter()
}

func (db *DB) CloseReader() error {
	sqlDB, err := db.Reader.DB()
	if err != nil {
		return errors.Wrap(err, "failed to get sql db")
	}

	return sqlDB.Close()
}

func (db *DB) CloseWriter() error {
	sqlDB, err := db.Writer.DB()
	if err != nil {
		return errors.Wrap(err, "failed to get sql db")
	}

	return sqlDB.Close()
}

func (db *DB) EnableLog(enable bool) {
	if enable {
		db.Reader.Logger.LogMode(logger.Info)
		db.Writer.Logger.LogMode(logger.Info)
	} else {
		db.Reader.Logger.LogMode(logger.Silent)
		db.Writer.Logger.LogMode(logger.Silent)
	}
}

func (db *DB) AutoMigrate(dst ...interface{}) error {
	if err := db.Reader.AutoMigrate(dst...); err != nil {
		return errors.Wrap(err, "failed to auto migrate reader")
	}

	if err := db.Writer.AutoMigrate(dst...); err != nil {
		return errors.Wrap(err, "failed to auto migrate writer")
	}

	return nil
}

func dbConnect(config *DBConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.User,
		config.Password,
		config.Addr,
		config.DBName,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		Logger: newLogger(config.LogFilename),
	})

	if err != nil {
		return nil, errors.Wrapf(err, "[db connection failed] Database name: %s", config.DBName)
	}

	if config.Log {
		db.Logger.LogMode(logger.Info)
	} else {
		db.Logger.LogMode(logger.Silent)
	}

	db.Set("gorm:table_options", "CHARSET=utf8mb4")

	sqlDB, err := db.DB()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get sql db")
	}

	if config.MaxOpenConnections != 0 {
		sqlDB.SetMaxOpenConns(config.MaxOpenConnections)
	}

	if config.MaxIdleConnections != 0 {
		sqlDB.SetMaxIdleConns(config.MaxIdleConnections)
	}

	if config.MaxLifetimeDuration != 0 {
		sqlDB.SetConnMaxLifetime(config.MaxLifetimeDuration)
	}

	return db, nil
}
