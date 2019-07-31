package database

import (
	"errors"
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

var (
	// ErrNotFound is returned when a resource cannot be found
	// in the database
	ErrNotFound  = errors.New("models: resource not found")
	ErrInvalidID = errors.New("models: ID provided was invalid")
)

type Config struct {
	DB     DBInfo `yaml:"db"`
	DBTest DBInfo `yaml:"dbtest"`
}

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

func (c *Config) loadConf(path string) *Config {

	file, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(file, c)
	if err != nil {
		panic(err)
	}

	return c
}

func NewDBConn(env string) (*DB, error) {
	var c Config
	c.loadConf(confPath)

	var d DBInfo
	if env == "dev" {
		d = c.DB
	} else if env == "test" {
		d = c.DBTest
	}

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

// HandleDBError will
func HandleDBError(db *gorm.DB) error {
	err := db.Error
	if err == gorm.ErrRecordNotFound {
		return ErrNotFound
	}

	return err
}
