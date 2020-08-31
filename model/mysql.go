package model

import (
	"database/sql"
)

type UserInfoDO struct {
	UserInfoID string
	UserID     string
	StuID      sql.NullString
	RealName   sql.NullString
	Sex        sql.NullString
	Major      sql.NullString
	ClassID    sql.NullString
	Grade      sql.NullString
	// EnrollDate 数据库中为datetime
	EnrollDate sql.NullTime
	ExtInfo    sql.NullString
}

func (u *UserInfoDO) GetStuID() string {
	if u.StuID.Valid {
		return u.StuID.String
	}
	return ""
}

func (u *UserInfoDO) GetRealName() string {
	if u.RealName.Valid {
		return u.RealName.String
	}
	return ""
}

func (u *UserInfoDO) GetSex() string {
	if u.Sex.Valid {
		return u.Sex.String
	}
	return ""
}

func (u *UserInfoDO) GetMajor() string {
	if u.Major.Valid {
		return u.Major.String
	}
	return ""
}

func (u *UserInfoDO) GetClassID() string {
	if u.ClassID.Valid {
		return u.ClassID.String
	}
	return ""
}

func (u *UserInfoDO) GetGrade() string {
	if u.Grade.Valid {
		return u.Grade.String
	}
	return ""
}

func (u *UserInfoDO) GetEnrollDate() int64 {
	if u.EnrollDate.Valid {
		return u.EnrollDate.Time.UnixNano()
	}
	return 0
}

func (u *UserInfoDO) GetExtInfo() string {
	if u.ExtInfo.Valid {
		return u.ExtInfo.String
	}
	return ""
}
