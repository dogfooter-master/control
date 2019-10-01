package service

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/mgo.v2"
	"time"
)

var mgoSession *mgo.Session
var mySqlDB *sql.DB

func init() {
	err := initializeMongo()
	if err != nil {
		panic(err)
	}

	// sql.DB 객체 생성
	dataSourceName := mySqlConfig.Username +
		":" + mySqlConfig.Password +
		"@tcp("+mySqlConfig.Hosts+")/"+
		mySqlConfig.Database
	mySqlDB, err = sql.Open("mysql", dataSourceName)
	if err != nil {
		panic(err)
	}
	go func() {
		for range time.Tick(time.Second * 10) {
			ConnectionCheck()
		}
	}()
}

func initializeMongo() (err error) {
	info := &mgo.DialInfo{
		Addrs:    []string{mgoConfig.Hosts},
		Timeout:  60 * time.Second,
		Database: mgoConfig.Database,
		Username: mgoConfig.Username,
		Password: mgoConfig.Password,
	}

	mgoSession, err = mgo.DialWithInfo(info)
	if err != nil {
		err = fmt.Errorf("fail to DialWithInfo(%#v) error - %v", info, err)
		return
	}

	return
}
