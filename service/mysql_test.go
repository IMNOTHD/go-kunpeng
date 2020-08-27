package service

import (
	"database/sql"
	"fmt"
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

	stmt, err := db.Prepare(_queryTest)
	if err != nil {
		log.Fatal(err)
		return
	}

	rows, err := stmt.Query()

	var majorId, realName *sql.NullString

	for rows.Next() {
		err = rows.Scan(&majorId, &realName)

		if err != nil {
			log.Fatal(err)
			return
		}

		fmt.Println(realName, majorId)
	}
}
