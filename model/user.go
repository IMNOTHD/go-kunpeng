package model

type UserInfo struct {
	UserInfoID string `json:"userInfoId"`
	UserID     string `json:"userId"`
	StuID      string `json:"stuId"`
	RealName   string `json:"realName"`
	Sex        string `json:"sex"`
	Major      string `json:"major"`
	ClassID    string `json:"classId"`
	Grade      string `json:"grade"`
	// EnrollDate 数据库中为datetime
	EnrollDate int64                  `json:"enrollDate"`
	ExtInfo    map[string]interface{} `json:"extInfo"`
}

type RoleInfo []string

type JobInfo []string
