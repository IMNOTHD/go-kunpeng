package service

import (
	"context"
	"encoding/json"
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
	var ctx = context.Background()

	userId := "189050214"
	userInfo := model.UserInfo{
		UserInfoID: "201811302142199201233000073047",
		UserID:     "201811302142192259540001201847",
		StuID:      "189050214",
		RealName:   "何锝升",
		Sex:        "男",
		Major:      "计算机科学与技术",
		ClassID:    "180923101",
		Grade:      "2018",
		EnrollDate: 1535731200000,
	}
	json.Unmarshal([]byte("{}"), &userInfo.ExtInfo)

	err = SetUserInfoRedis(c, userId, &userInfo)

	if err != nil {
		log.Fatal("Redis HMSet Error:", err)
	}

	roleInfo := model.RoleInfo{"p1", "p4"}
	err = AddRoleInfoRedis(c, userId, &roleInfo)

	if err != nil {
		log.Fatal(err.Error())
	}

	v, err := c.SMembers(ctx, _roleInfoRedisKey+userId).Result()

	fmt.Println(v)
}
