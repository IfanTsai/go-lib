package mysql

import (
	"fmt"

	"github.com/IfanTsai/go-lib/config"
	"github.com/pkg/errors"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type DB struct {
	Read  *gorm.DB
	Write *gorm.DB
}

type DBConfig struct {
	User     string
	Password string
	Addr     string
	Name     string
}

func NewDB(readConfig, writeConfig *DBConfig) (*DB, error) {
	dbRead, err := dbConnect(readConfig.User, readConfig.Password, readConfig.Password, readConfig.Name)
	if err != nil {
		return nil, err
	}

	dbWrite, err := dbConnect(writeConfig.User, writeConfig.Password, writeConfig.Password, writeConfig.Name)
	if err != nil {
		return nil, err
	}

	return &DB{
		Read:  dbRead,
		Write: dbWrite,
	}, nil
}

func NewDBInConfig() (*DB, error) {
	readConfig := &DBConfig{
		User:     config.GetConfig().MySQL.Read.User,
		Password: config.GetConfig().MySQL.Read.Password,
		Addr:     config.GetConfig().MySQL.Read.Addr,
		Name:     config.GetConfig().MySQL.Read.Name,
	}

	writeConfig := &DBConfig{
		User:     config.GetConfig().MySQL.Write.User,
		Password: config.GetConfig().MySQL.Write.Password,
		Addr:     config.GetConfig().MySQL.Write.Addr,
		Name:     config.GetConfig().MySQL.Write.Name,
	}

	return NewDB(readConfig, writeConfig)
}

func (db *DB) Close() error {
	err := db.DBReadClose()
	if err != nil {
		return err
	}

	return db.DBWriteClose()
}

func (db *DB) DBReadClose() error {
	sqlDB, err := db.Read.DB()
	if err != nil {
		return errors.Wrap(err, "failed to get sql db")
	}

	return sqlDB.Close()
}

func (db *DB) DBWriteClose() error {
	sqlDB, err := db.Write.DB()
	if err != nil {
		return errors.Wrap(err, "failed to get sql db")
	}

	return sqlDB.Close()
}

func dbConnect(user, password, addr, dbName string) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user,
		password,
		addr,
		dbName,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		//Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		return nil, errors.Wrapf(err, "[db connection failed] Database name: %s", dbName)
	}

	db.Set("gorm:table_options", "CHARSET=utf8mb4")

	sqlDB, err := db.DB()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get sql db")
	}

	sqlDB.SetMaxOpenConns(config.GetConfig().MySQL.MaxOpenConnections)
	sqlDB.SetMaxIdleConns(config.GetConfig().MySQL.MaxIdleConnections)
	sqlDB.SetConnMaxLifetime(config.GetConfig().MySQL.MaxLifetimeDuration)

	return db, nil
}
