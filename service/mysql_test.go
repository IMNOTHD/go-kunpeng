package service

import (
	"log"
	"testing"
)

const (
	_queryTest = `select major_id, real_name from common_user_info`
)

func TestCreateMysqlWorker(t *testing.T) {
	var err error
	db, err := CreateMysqlWorker()
	if err != nil {
		log.Fatal(err)
		return
	}

	c, err := CreateRedisClient()
	if err != nil {
		log.Fatal(err.Error())
		return
	}

	_ = db.Close()
	_ = c.Close()
}
