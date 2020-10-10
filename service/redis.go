package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/go-redis/redis/v8"

	rc "go-kunpeng/config/redis"
	"go-kunpeng/model"
)

const (
	_userRedisKey           = "betahouse:haetae:user:"
	_userInfoRedisKey       = "betahouse:haetae:user:userInfo"
	_roleInfoRedisKey       = "betahouse:haetae:user:roleInfo:"
	_jobInfoRedisKey        = "betahouse:haetae:user:jobInfo:"
	_avatarUrlRedisKey      = "betahouse:haetae:user:avatarUrl:"
	_activityRecordRedisKey = "betahouse:haetae:activity:activityRecord:"
	_activityRedisKey       = "betahouse:haetae:activity:activity:"
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

func CleanUserAllInfo(c *redis.Client, userId string) error {
	var ctx = context.Background()

	pipe := c.TxPipeline()

	defer pipe.Close()

	pipe.HDel(ctx, _userInfoRedisKey, userId)
	pipe.Del(ctx, _roleInfoRedisKey+userId)
	pipe.Del(ctx, _jobInfoRedisKey+userId)
	pipe.Del(ctx, _avatarUrlRedisKey+userId)

	_, err := pipe.Exec(ctx)

	if err != nil {
		_ = pipe.Discard()
		return err
	}

	return nil
}

func CleanAllUserAllInfo(c *redis.Client) (int32, error) {
	var ctx = context.Background()

	k, err := c.Keys(ctx, _userRedisKey+"*").Result()
	if err != nil {
		return 0, err
	}

	var successCount int32 = 0
	isFullSuccess := true

	for _, v := range k {
		_, err := c.Del(ctx, v).Result()
		if err != nil {
			isFullSuccess = false
		}
		successCount++
	}

	if isFullSuccess {
		return successCount, nil
	} else {
		return successCount, errors.New("some remove failed")
	}
}

func CleanUserAllActivityRecord(c *redis.Client, userId string) error {
	var ctx = context.Background()

	multiKey := fmt.Sprintf("%s%s*", _activityRecordRedisKey, userId)

	k, err := c.Keys(ctx, multiKey).Result()
	if err != nil {
		return err
	}

	pipe := c.TxPipeline()

	defer pipe.Close()

	for _, v := range k {
		pipe.Del(ctx, v)
	}

	_, err = pipe.Exec(ctx)

	if err != nil {
		_ = pipe.Discard()
		return err
	}

	return nil
}

func CleanAllUserAllActivityRecord(c *redis.Client) (int32, error) {
	var ctx = context.Background()

	k, err := c.Keys(ctx, _activityRecordRedisKey+"*").Result()
	if err != nil {
		return 0, err
	}

	var successCount int32 = 0
	isFullSuccess := true

	for _, v := range k {
		_, err := c.Del(ctx, v).Result()
		if err != nil {
			isFullSuccess = false
		}
		successCount++
	}

	if isFullSuccess {
		return successCount, nil
	} else {
		return successCount, errors.New("some remove failed")
	}
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

func CleanUserInfoRedis(c *redis.Client, userId string) error {
	var ctx = context.Background()

	val := c.HDel(ctx, _userInfoRedisKey, userId)

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

func CleanRoleInfoRedis(c *redis.Client, userId string) error {
	var ctx = context.Background()

	val := c.Del(ctx, _roleInfoRedisKey+userId)

	if val.Err() != nil {
		return val.Err()
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

func CleanJobInfoRedis(c *redis.Client, userId string) error {
	var ctx = context.Background()

	val := c.Del(ctx, _jobInfoRedisKey+userId)

	if val.Err() != nil {
		return val.Err()
	}

	return nil
}

func SetAvatarUrlRedis(c *redis.Client, userId string, avatarUrl *model.AvatarUrl) error {
	var ctx = context.Background()

	val := c.Set(ctx, _avatarUrlRedisKey+userId, avatarUrl.Url, 0)

	if val.Err() != nil {
		return val.Err()
	}

	return nil
}

func SetActivityRedis(c *redis.Client, activityId string, activity *model.Activity) error {
	var ctx = context.Background()

	if activityId == "" {
		return errors.New("activityId null")
	}

	record, _ := json.Marshal(activity)
	val := c.Set(ctx, _activityRedisKey+activityId, record, 0)

	if val.Err() != nil {
		return val.Err()
	}

	return nil
}

func CleanActivityRedis(c *redis.Client, activityId string) error {
	var ctx = context.Background()

	if activityId == "" {
		return errors.New("activityId null")
	}

	val := c.Del(ctx, _activityRedisKey+activityId)

	if val.Err() != nil {
		return val.Err()
	}

	return nil
}

func AddActivityRecordRedis(c *redis.Client, activityRecord *[]model.ActivityRecord) error {
	var ctx = context.Background()

	for _, v := range *activityRecord {
		record, _ := json.Marshal(v)
		key := fmt.Sprintf("%s%s:%s", _activityRecordRedisKey, v.UserID, v.Type)

		count, err := c.ZCard(ctx, key).Result()
		if err != nil {
			return err
		}

		//fmt.Println(count)

		val := c.ZAdd(ctx, key, &redis.Z{Score: float64(count), Member: record})

		if val.Err() != nil {
			return val.Err()
		}
	}

	return nil
}
