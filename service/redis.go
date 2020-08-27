package service

import (
	"context"
	"encoding/json"

	"github.com/go-redis/redis/v8"

	rc "go-kunpeng/config/redis"
	"go-kunpeng/model"
)

const (
	_userInfoRedisKey = "betahouse:heatae:user:userInfo"
	_roleInfoRedisKey = "betahouse:heatae:user:roleInfo"
	_jobInfoRedisKey  = "betahouse:heatae:user:jobInfo"
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

func SetUserInfoRedis(c *redis.Client, userId string, userInfo *model.UserInfo) error {
	var ctx = context.Background()

	info, _ := json.Marshal(userInfo)
	val := c.HSet(ctx, _userInfoRedisKey, userId, info)

	if val.Err() != nil {
		return val.Err()
	}

	return nil
}

func SetRoleInfoRedis(c *redis.Client, userId string, roleInfo *model.RoleInfo) error {
	var ctx = context.Background()

	info, _ := json.Marshal(roleInfo)
	val := c.HSet(ctx, _roleInfoRedisKey, userId, info)

	if val.Err() != nil {
		return val.Err()
	}

	return nil
}

func SetJobInfoRedis(c *redis.Client, userId string, jobInfo *model.JobInfo) error {
	var ctx = context.Background()

	info, _ := json.Marshal(jobInfo)
	val := c.HSet(ctx, _jobInfoRedisKey, userId, info)

	if val.Err() != nil {
		return val.Err()
	}

	return nil
}
