package service

import (
	"fmt"
	"log"
	"testing"
	"time"
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

	u, err := QueryUserByGrade(db, "2015")

	if err != nil {
		log.Fatal(err)
		return
	}

	c, err := CreateRedisClient()
	if err != nil {
		log.Fatal(err.Error())
		return
	}

	s := 0
	ti := time.Now()
	for _, i := range u {
		tt := time.Now()
		err = SetUserInfoRedis(c, &i)
		if err != nil {
			log.Fatal(err.Error())
		}
		s++
		fmt.Println("Time add one user: ", time.Since(tt))
	}
	elapsed := time.Since(ti)
	fmt.Println("App elapsed: ", elapsed)
	fmt.Println("Success add: ", s)

	_ = db.Close()
	_ = c.Close()
}
