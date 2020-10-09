package model

type UserInfo struct {
	UserInfoID string `json:"userInfoId" db:"user_info_id"`
	UserID     string `json:"userId" db:"user_id"`
	StuID      string `json:"stuId" db:"stu_id"`
	RealName   string `json:"realName" db:"real_name"`
	Sex        string `json:"sex" db:"sex"`
	Major      string `json:"major" db:"major_id"`
	ClassID    string `json:"classId" db:"class_id"`
	Grade      string `json:"grade" db:"grade"`
	// EnrollDate 数据库中为datetime
	EnrollDate int64                  `json:"enrollDate" db:"enroll_date"`
	ExtInfo    map[string]interface{} `json:"extInfo" db:"ext_info"`
}

type RoleInfo []string

type JobInfo []string

type AvatarUrl struct {
	Url string
}

type User struct {
	UserInfo  UserInfo
	RoleInfo  RoleInfo
	JobInfo   JobInfo
	AvatarUrl AvatarUrl
}
