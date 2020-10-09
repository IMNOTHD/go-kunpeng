package model

type ActivityRecord struct {
	ActivityRecordID string `json:"activityRecordId" db:"activity_record_id"`
	ActivityID       string `json:"activityId" db:"activity_id"`
	UserID           string `json:"userId" db:"user_id"`
	ScannerUserID    string `json:"scannerUserId" db:"scanner_user_id"`
	Time             int    `json:"time" db:"time"`
	Type             string `json:"type" db:"type"`
	Status           string `json:"status" db:"status"`
	Term             string `json:"term" db:"term"`
	Grades           string `json:"grades" db:"grades"`
	// ExtInfo 数据库中存储的是json字符串
	ExtInfo    map[string]interface{} `json:"extInfo" db:"ext_info"`
	CreateTime int64                  `json:"createTime" db:"gmt_create"`
	// ActivityTime 只是Time的格式化
	ActivityTime string `json:"activityTime"`
	ScannerName  string `json:"scannerName"`
}

type Activity struct {
	ActivityName        string `json:"activityName" db:"activity_name"`
	OrganizationMessage string `json:"organizationMessage" db:"organization_message"`
	Location            string `json:"location" db:"location"`
	StartTime           int64  `json:"startTime" db:"start"`
	EndTime             int64  `json:"endTime" db:"end"`
	Score               int64  `json:"score" db:"score"`
}
