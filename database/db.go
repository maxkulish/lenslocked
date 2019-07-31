package database

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

const (
	confPath = "/Users/mk/Code/lenslocked/config.yaml"
	dialect  = "postgres"
)

type DBInfo struct {
	Host     string `yaml:"host"`
	Database string `yaml:"database"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

type DB struct {
	Conn *gorm.DB
}

func (d *DB) Close() {
	d.Conn.Close()
}

func (d *DBInfo) loadConf(path string) *DBInfo {

	file, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(file, d)
	if err != nil {
		panic(err)
	}

	return d
}

func NewDBConn() (*DB, error) {
	var d DBInfo
	d.loadConf(confPath)

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		d.Host, d.Port, d.User, d.Password, d.Database)

	var db DB
	var err error
	db.Conn, err = gorm.Open(dialect, psqlInfo)
	if err != nil {
		return nil, err
	}

	db.Conn.LogMode(true)

	return &db, nil

}
