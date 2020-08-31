package service

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	mc "go-kunpeng/config/mysql"
	"go-kunpeng/model"

	_ "github.com/go-sql-driver/mysql"
)

const (
	_queryUserInfoByUserId = `select user_info_id, user_id, stu_id, real_name, sex, major_id, class_id, grade, enroll_date, ext_info from common_user_info where user_id = ?`
	_queryUserInfoByGrade  = `select user_info_id, user_id, stu_id, real_name, sex, major_id, class_id, grade, enroll_date, ext_info from common_user_info where grade = ?`
	_queryUserInfoByClass  = `select user_info_id, user_id, stu_id, real_name, sex, major_id, class_id, grade, enroll_date, ext_info from common_user_info where class = ?`
	_queryUserInfoAll      = `select user_info_id, user_id, stu_id, real_name, sex, major_id, class_id, grade, enroll_date, ext_info from common_user_info`
	_queryRoleIdByUserID   = `select role_id from common_user_role_relation where user_id = ?`
	_queryRoleCodeByRoleId = `select role_code from common_role where role_id = ?`
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

func QueryUserInfoByUserId(db *sql.DB, userId string) (*model.UserInfo, *model.RoleInfo, *model.JobInfo, error) {
	var err error

	uStmt, err := db.Prepare(_queryUserInfoByUserId)
	if err != nil {
		return nil, nil, nil, err
	}

	row := uStmt.QueryRow(userId)

	if row == nil {
		return nil, nil, nil, errors.New("no such user")
	}

	var u model.UserInfoDO

	err = row.Scan(&u.UserInfoID, &u.UserID, &u.StuID, &u.RealName, &u.Sex, &u.Major, &u.ClassID, &u.Grade, &u.EnrollDate, &u.ExtInfo)
	if err != nil {
		return nil, nil, nil, err
	}

	e := make(map[string]interface{})
	_ = json.Unmarshal([]byte(u.GetExtInfo()), &e)

	rStmt, err := db.Prepare(_queryRoleIdByUserID)
	if err != nil {
		return nil, nil, nil, err
	}

	rows, err := rStmt.Query(userId)
	if err != nil {
		return nil, nil, nil, err
	}

	r := make(model.RoleInfo, 0)

	if rows != nil {
		var roleId, roleCode string
		cStmt, err := db.Prepare(_queryRoleCodeByRoleId)
		for rows.Next() {
			err = rows.Scan(&roleId)
			if err != nil {
				return nil, nil, nil, err
			}

			row := cStmt.QueryRow(roleId)
			err = row.Scan(&roleCode)
			if err != nil {
				return nil, nil, nil, err
			}

			r = append(r, roleCode)
		}
	}

	j := make(model.JobInfo, 0)

	return &model.UserInfo{
		UserInfoID: u.UserInfoID,
		UserID:     u.UserID,
		StuID:      u.GetStuID(),
		RealName:   u.GetRealName(),
		Sex:        u.GetSex(),
		Major:      u.GetMajor(),
		ClassID:    u.GetClassID(),
		Grade:      u.GetGrade(),
		EnrollDate: u.GetEnrollDate(),
		ExtInfo:    e,
	}, &r, &j, nil
}

func QueryUserInfoByGrade(db *sql.DB, grade string) ([]model.UserInfo, error) {
	var err error

	stmt, err := db.Prepare(_queryUserInfoByGrade)
	if err != nil {
		return nil, err
	}

	rows, err := stmt.Query(grade)

	if rows == nil {
		return nil, errors.New("no such user")
	}
	if err != nil {
		return nil, err
	}

	var u []model.UserInfo

	for rows.Next() {
		var ud model.UserInfoDO

		err = rows.Scan(&ud.UserInfoID, &ud.UserID, &ud.StuID, &ud.RealName, &ud.Sex, &ud.Major, &ud.ClassID, &ud.Grade, &ud.EnrollDate, &ud.ExtInfo)
		if err != nil {
			return nil, err
		}

		e := make(map[string]interface{})
		_ = json.Unmarshal([]byte(ud.GetExtInfo()), &e)

		u = append(u, model.UserInfo{
			UserInfoID: ud.UserInfoID,
			UserID:     ud.UserID,
			StuID:      ud.GetStuID(),
			RealName:   ud.GetRealName(),
			Sex:        ud.GetSex(),
			Major:      ud.GetMajor(),
			ClassID:    ud.GetClassID(),
			Grade:      ud.GetGrade(),
			EnrollDate: ud.GetEnrollDate(),
			ExtInfo:    e,
		})
	}

	return u, nil
}

func QueryUserInfoByClass(db *sql.DB, class string) ([]model.UserInfo, error) {
	var err error

	stmt, err := db.Prepare(_queryUserInfoByClass)
	if err != nil {
		return nil, err
	}

	rows, err := stmt.Query(class)

	if rows == nil {
		return nil, errors.New("no such user")
	}
	if err != nil {
		return nil, err
	}

	var u []model.UserInfo

	for rows.Next() {
		var ud model.UserInfoDO

		err = rows.Scan(&ud.UserInfoID, &ud.UserID, &ud.StuID, &ud.RealName, &ud.Sex, &ud.Major, &ud.ClassID, &ud.Grade, &ud.EnrollDate, &ud.ExtInfo)
		if err != nil {
			return nil, err
		}

		e := make(map[string]interface{})
		_ = json.Unmarshal([]byte(ud.GetExtInfo()), &e)

		u = append(u, model.UserInfo{
			UserInfoID: ud.UserInfoID,
			UserID:     ud.UserID,
			StuID:      ud.GetStuID(),
			RealName:   ud.GetRealName(),
			Sex:        ud.GetSex(),
			Major:      ud.GetMajor(),
			ClassID:    ud.GetClassID(),
			Grade:      ud.GetGrade(),
			EnrollDate: ud.GetEnrollDate(),
			ExtInfo:    e,
		})
	}

	return u, nil
}

func QueryUserInfoAll(db *sql.DB) ([]model.UserInfo, error) {
	var err error

	stmt, err := db.Prepare(_queryUserInfoAll)
	if err != nil {
		return nil, err
	}

	rows, err := stmt.Query()

	if rows == nil {
		return nil, errors.New("no such user")
	}
	if err != nil {
		return nil, err
	}

	var u []model.UserInfo

	for rows.Next() {
		var ud model.UserInfoDO

		err = rows.Scan(&ud.UserInfoID, &ud.UserID, &ud.StuID, &ud.RealName, &ud.Sex, &ud.Major, &ud.ClassID, &ud.Grade, &ud.EnrollDate, &ud.ExtInfo)
		if err != nil {
			return nil, err
		}

		e := make(map[string]interface{})
		_ = json.Unmarshal([]byte(ud.GetExtInfo()), &e)

		u = append(u, model.UserInfo{
			UserInfoID: ud.UserInfoID,
			UserID:     ud.UserID,
			StuID:      ud.GetStuID(),
			RealName:   ud.GetRealName(),
			Sex:        ud.GetSex(),
			Major:      ud.GetMajor(),
			ClassID:    ud.GetClassID(),
			Grade:      ud.GetGrade(),
			EnrollDate: ud.GetEnrollDate(),
			ExtInfo:    e,
		})
	}

	return u, nil
}
