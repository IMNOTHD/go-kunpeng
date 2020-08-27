package service

import (
	"database/sql"
	"fmt"

	mc "go-kunpeng/config/mysql"

	_ "github.com/go-sql-driver/mysql"
)

func CreateMysqlWorker() (*sql.DB, error) {
	// [username[:password]@][protocol[(address)]]/dbname[?param1=value1&...&paramN=valueN]
	dsn := fmt.Sprintf("%s:%s@%s(%s:%d)/%s?%s", mc.Username, mc.Password, mc.Protocol, mc.Address, mc.Port, mc.Dbname, mc.Addition)

	db, err := sql.Open("mysql", dsn)

	if err != nil {
		return nil, err
	}

	err = db.Ping()

	if err != nil {
		return nil, err
	}

	return db, nil
}
