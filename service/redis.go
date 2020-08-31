package service

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis/v8"

	rc "go-kunpeng/config/redis"
	"go-kunpeng/model"
)

const (
	_userInfoRedisKey       = "betahouse:heatae:user:userInfo"
	_roleInfoRedisKey       = "betahouse:heatae:user:roleInfo:"
	_jobInfoRedisKey        = "betahouse:heatae:user:jobInfo:"
	_activityRecordRedisKey = "betahouse:heatae:activity:activityRecord"
)

func CreateRedisClient() (*redis.Client, error) {
	c := redis.NewClient(&redis.Options{
		Addr:     rc.Addr,
		Password: rc.Password,
		DB:       rc.DB,
	})

	var ctx = context.Background()
	_, err := c.Ping(ctx).Result()

	return c, err
}

func SetUserInfoRedis(c *redis.Client, userInfo *model.UserInfo) error {
	var ctx = context.Background()

	info, _ := json.Marshal(userInfo)
	val := c.HSet(ctx, _userInfoRedisKey, userInfo.UserID, info)

	if val.Err() != nil {
		return val.Err()
	}

	return nil
}

func AddRoleInfoRedis(c *redis.Client, userId string, roleInfo *model.RoleInfo) error {
	var ctx = context.Background()

	for _, role := range *roleInfo {
		val := c.SAdd(ctx, _roleInfoRedisKey+userId, role)

		if val.Err() != nil {
			return val.Err()
		}
	}

	return nil
}

func AddJobInfoRedis(c *redis.Client, userId string, jobInfo *model.JobInfo) error {
	var ctx = context.Background()

	for _, job := range *jobInfo {
		val := c.SAdd(ctx, _jobInfoRedisKey+userId, job)

		if val.Err() != nil {
			return val.Err()
		}
	}

	return nil
}

func AddActivityRecordRedis(c *redis.Client, activityRecord *model.ActivityRecord) error {
	var ctx = context.Background()

	record, _ := json.Marshal(activityRecord)
	key := fmt.Sprintf("%s:%s:%s", _activityRecordRedisKey, activityRecord.UserID, activityRecord.Type)

	count, err := c.ZCard(ctx, key).Result()
	if err != nil {
		return err
	}

	fmt.Println(count)

	val := c.ZAdd(ctx, key, &redis.Z{Score: float64(count), Member: record})

	if val.Err() != nil {
		return val.Err()
	}

	return nil
}
