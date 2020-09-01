package model

type ActivityRecord struct {
	ActivityRecordID string `json:"activityRecordId"`
	ActivityID       string `json:"activityId"`
	UserID           string `json:"userId"`
	ScannerUserID    string `json:"scannerUserId"`
	Time             int    `json:"time"`
	Type             string `json:"type"`
	Status           string `json:"status"`
	Term             string `json:"term"`
	Grades           string `json:"grades"`
	// ExtInfo 数据库中存储的是json字符串
	ExtInfo             map[string]interface{} `json:"extInfo"`
	CreateTime          int64                  `json:"createTime"`
	ActivityName        string                 `json:"activityName"`
	OrganizationMessage string                 `json:"organizationMessage"`
	Location            string                 `json:"location"`
	StartTime           int64                  `json:"startTime"`
	EndTime             int64                  `json:"endTime"`
	Score               int64                  `json:"score"`
	// ActivityTime 只是Time的格式化
	ActivityTime string `json:"activityTime"`
	ScannerName  string `json:"scannerName"`
}
