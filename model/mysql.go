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

type JobInfoDO struct {
	OrganizationName  string
	MemberDescription sql.NullString
}

type AvatarUrlDO struct {
	Url sql.NullString
}

type ActivityRecordDO struct {
	ActivityRecordId sql.NullString
	ActivityId       string
	UserId           string
	ScannerUserId    string
	Time             sql.NullInt32
	Type             sql.NullString
	Status           sql.NullString
	Term             sql.NullString
	Grades           sql.NullString
	ExtInfo          sql.NullString
	// GmtCreate 数据库中为datetime
	GmtCreate sql.NullTime
}

type ActivityDO struct {
	ActivityName        sql.NullString
	OrganizationMessage sql.NullString
	Location            sql.NullString
	Start               sql.NullTime
	End                 sql.NullTime
	Score               sql.NullInt64
}

// -----------------------------------

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

// -----------------------------------

func (j *JobInfoDO) GetMemberDescription() string {
	if j.MemberDescription.Valid {
		return j.MemberDescription.String
	}
	return ""
}

// -----------------------------------

func (a AvatarUrlDO) GetAvatarUrl() string {
	if a.Url.Valid {
		return a.Url.String
	}
	return ""
}

// -----------------------------------

func (a *ActivityRecordDO) GetActivityRecordId() string {
	if a.ActivityRecordId.Valid {
		return a.ActivityRecordId.String
	}
	return ""
}

func (a *ActivityRecordDO) GetTime() int32 {
	if a.Time.Valid {
		return a.Time.Int32
	}
	return 0
}

func (a *ActivityRecordDO) GetType() string {
	if a.Type.Valid {
		return a.Type.String
	}
	return ""
}

func (a *ActivityRecordDO) GetStatus() string {
	if a.Status.Valid {
		return a.Status.String
	}
	return ""
}

func (a *ActivityRecordDO) GetTerm() string {
	if a.Term.Valid {
		return a.Term.String
	}
	return ""
}

func (a *ActivityRecordDO) GetGrades() string {
	if a.Grades.Valid {
		return a.Grades.String
	}
	return ""
}

func (a *ActivityRecordDO) GetExtInfo() string {
	if a.ExtInfo.Valid {
		return a.ExtInfo.String
	}
	return ""
}

func (a *ActivityRecordDO) GetGmtCreate() int64 {
	if a.GmtCreate.Valid {
		return a.GmtCreate.Time.UnixNano()
	}
	return 0
}

// -----------------------------------

func (a *ActivityDO) GetActivityName() string {
	if a.ActivityName.Valid {
		return a.ActivityName.String
	}
	return ""
}

func (a *ActivityDO) GetOrganizationMessage() string {
	if a.OrganizationMessage.Valid {
		return a.OrganizationMessage.String
	}
	return ""
}

func (a *ActivityDO) GetLocation() string {
	if a.Location.Valid {
		return a.Location.String
	}
	return ""
}

func (a *ActivityDO) GetStart() int64 {
	if a.Start.Valid {
		return a.Start.Time.UnixNano()
	}
	return 0
}

func (a *ActivityDO) GetEnd() int64 {
	if a.End.Valid {
		return a.End.Time.UnixNano()
	}
	return 0
}

func (a *ActivityDO) GetScore() int64 {
	if a.Score.Valid {
		return a.Score.Int64
	}
	return 0
}

// -----------------------------------
