package service

import (
	"context"
	"fmt"
	"log"
	"testing"

	"go-kunpeng/model"
)

func TestCreateRedisClient(t *testing.T) {
	c, err := CreateRedisClient()
	if err != nil {
		log.Fatal(err.Error())
		return
	}

	err = SetAvatarUrlRedis(c, "201811302142192259540001201847", &model.AvatarUrl{Url: "123"})
	if err != nil {
		log.Fatal(err)
	}

	var ctx = context.Background()

	a := &model.ActivityRecord{
		ActivityRecordID: "201905292002334060129610022018",
		ActivityID:       "201905241117049855811810012019",
		UserID:           "201811302142192259540001201847",
		ScannerUserID:    "201811302142173663860001201808",
		Time:             0,
		Type:             "schoolActivity",
		Status:           "ENABLE",
		Term:             "2018B",
		Grades:           "",
		ExtInfo: map[string]interface{}{
			"scannerName": "庄子琛",
		},
		CreateTime:   1559131353000,
		ActivityTime: "0.0",
		ScannerName:  "庄子琛",
	}

	x := make([]model.ActivityRecord, 0)
	x = append(x, *a)

	err = AddActivityRecordRedis(c, &x, false)

	if err != nil {
		log.Fatal(err)
	}

	v, err := c.ZRange(ctx, fmt.Sprintf("%s:%s:%s", _activityRecordRedisKey, a.UserID, a.Type), 0, -1).Result()

	log.Println(v)
}
